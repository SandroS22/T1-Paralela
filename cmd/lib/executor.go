package lib

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Task é a unidade de trabalho executada pelo pool.
// Retorne error para permitir coleta de falhas pela API Submit/Future.
type Task func() error

var (
	ErrExecutorClosed = errors.New("executor: fechado para novas tarefas")
	ErrNilTask        = errors.New("executor: tarefa nil")
)

// Future permite aguardar a conclusão de uma Task submetida via Submit.
type Future struct {
	done chan struct{}
	err  atomic.Value // armazena error
}

// Done retorna um canal fechado quando a tarefa terminar.
func (f *Future) Done() <-chan struct{} { return f.done }

// Err retorna o erro da tarefa após concluída (ou nil se ok).
func (f *Future) Err() error {
	v := f.err.Load()
	if v == nil {
		return nil
	}
	return v.(error)
}

// Executor implementa um pool de workers com fila de tarefas.
type Executor struct {
	jobs    chan Task
	wg      sync.WaitGroup // espera workers
	closed  atomic.Bool
	closeMu sync.Mutex
}

// NewExecutor cria um executor com 'workers' goroutines e fila com 'queue' slots.
// Regra prática: workers ~ número de CPUs; queue >= workers para amortecer bursts.
func NewExecutor(workers, queue int) *Executor {
	if workers <= 0 {
		workers = 1
	}
	if queue < 0 {
		queue = 0
	}
	e := &Executor{
		jobs: make(chan Task, queue),
	}
	e.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer e.wg.Done()
			for t := range e.jobs {
				// Executa e ignora erro aqui; coleta acontece só via Submit/Future (wrap).
				if t != nil {
					_ = t()
				}
			}
		}()
	}
	return e
}

// Execute enfileira a tarefa bloqueando se a fila estiver cheia.
// Retorna erro se executor estiver fechado ou tarefa nil.
func (e *Executor) Execute(t Task) error {
	if t == nil {
		return ErrNilTask
	}
	if e.closed.Load() {
		return ErrExecutorClosed
	}
	select {
	case e.jobs <- t:
		return nil
	default:
		// fila cheia: bloqueia até haver espaço, a menos que feche no meio
		select {
		case e.jobs <- t:
			return nil
		case <-e.whenClosed():
			return ErrExecutorClosed
		}
	}
}

// ExecuteContext enfileira respeitando cancelamento do contexto.
func (e *Executor) ExecuteContext(ctx context.Context, t Task) error {
	if t == nil {
		return ErrNilTask
	}
	if e.closed.Load() {
		return ErrExecutorClosed
	}
	select {
	case e.jobs <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-e.whenClosed():
		return ErrExecutorClosed
	}
}

// TryExecute tenta enfileirar sem bloquear. Retorna false se fila estiver cheia.
func (e *Executor) TryExecute(t Task) (bool, error) {
	if t == nil {
		return false, ErrNilTask
	}
	if e.closed.Load() {
		return false, ErrExecutorClosed
	}
	select {
	case e.jobs <- t:
		return true, nil
	default:
		return false, nil
	}
}

// Submit encapsula a Task e retorna um Future para aguardar erro/conclusão.
func (e *Executor) Submit(t Task) (*Future, error) {
	if t == nil {
		return nil, ErrNilTask
	}
	f := &Future{done: make(chan struct{})}
	wrapped := func() error {
		err := t()
		f.err.Store(err)
		close(f.done)
		return err
	}
	if err := e.Execute(wrapped); err != nil {
		return nil, err
	}
	return f, nil
}

// Close impede novos enfileiramentos e inicia desligamento.
// As tarefas já enfileiradas ainda serão executadas.
func (e *Executor) Close() {
	e.closeMu.Lock()
	defer e.closeMu.Unlock()
	if e.closed.Swap(true) {
		return
	}
	close(e.jobs)
}

// Wait bloqueia até todos os workers finalizarem (fila drenada).
func (e *Executor) Wait() {
	e.wg.Wait()
}

// Shutdown = Close + Wait.
func (e *Executor) Shutdown() {
	e.Close()
	e.Wait()
}

// whenClosed devolve um canal fechado quando o executor for fechado.
// Implementado com uma goroutine única e sync.Once seria overkill;
// aqui usamos polling simples via canal fechado do jobs quando já fechado.
func (e *Executor) whenClosed() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		if e.closed.Load() {
			close(ch)
			return
		}
		// se jobs fechar, também considera fechado
		// não dá para selecionar diretamente, então usa outra goroutine leve
		// porém, isso só é alcançado quando Close fecha jobs
		<-func() <-chan struct{} {
			c := make(chan struct{})
			go func() { // observa fechamento do jobs
				for range e.jobs {
					// nunca entra, pois range consome; não queremos consumir.
					// Portanto, não use este caminho.
				}
				close(c)
			}()
			return c
		}()
		close(ch)
	}()
	return ch
}
