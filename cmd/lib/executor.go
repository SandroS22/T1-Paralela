package lib

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Task é a unidade de trabalho executada pelos workers.
// Retorne error para permitir coleta de falhas pela API Submit/Future.
type Task func() error

var (
	ErrExecutorClosed = errors.New("executor: fechado para novas tarefas")
	ErrNilTask        = errors.New("executor: tarefa nil")
)

// Executor implementa um pool de workers com fila de tarefas.
type Executor struct {
	jobs   chan Task
	wg     sync.WaitGroup // espera workers
	closed atomic.Bool
	mu     sync.Mutex // protege fechamento + enfileiramento
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
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed.Load() {
		return ErrExecutorClosed
	}

	e.jobs <- t
	return nil
}

// Close impede novos enfileiramentos e inicia desligamento.
// As tarefas já enfileiradas ainda serão executadas.
func (e *Executor) Close() {
	e.mu.Lock()
	if e.closed.Swap(true) {
		e.mu.Unlock()
		return
	}
	close(e.jobs)
	e.mu.Unlock()
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
