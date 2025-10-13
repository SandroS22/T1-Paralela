package main

import (
	"log"
	"paralela/cmd/lib"
	"paralela/cmd/parallel"
	"paralela/cmd/seq"
	"paralela/cmd/util"
	"runtime"
	"sync"
	"time"
)

func main() {
	// Busca o número de núcleos físicos disponíveis para goroutines
	workers := runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(workers)
	log.Printf("Tempo total (Número máximo de goroutines (workers): %d\n", workers)

	// Definições para o teste fixo
	const (
		ARRAY_SIZE = 1000000
		NUM_TASKS  = 16 // Número de vetores = tarefas
		FIXED_SEED = 42
	)

	// Gera o conjunto de tarefas: vetores independentes. Cada vetor = 1 tarefa
	datasets := make([][]int, NUM_TASKS)
	for i := 0; i < NUM_TASKS; i++ {
		// Seeds diferentes por tarefa, mas determinísticas.
		datasets[i] = util.GenerateDeterministicArray(ARRAY_SIZE, int64(FIXED_SEED+int64(i)))
	}

	// 1. Método de processamento sequencial
	seqInputs := cloneDatasets(datasets)
	start := time.Now()

	for i := 0; i < NUM_TASKS; i++ {
		_ = seq.MergeSort(seqInputs[i])
	}

	elapsedSeq := time.Since(start)
	log.Printf("Tempo de execução (Tarefas sequênciais): %s\n", elapsedSeq)

	// 2A. Método de processamento paralelo entre tarefas
	// Cada tarefa usa MergeSort sequencial
	parInputs := cloneDatasets(datasets)
	start = time.Now()

	runTasksInParallel(workers, parInputs, func(arr []int) {
		_ = seq.MergeSort(arr)
	})

	elapsedTimeParSeq := time.Since(start)
	log.Printf("Tempo total (Tarefas em Paralelo, mergesort seq): %s\n", elapsedTimeParSeq)

	// 2B. Método de processamento paralelo entre tarefas
	// Cada tarefa usa MergeSort paralelo
	parAlgoInputs := cloneDatasets(datasets)
	start = time.Now()

	runTasksInParallel(workers, parAlgoInputs, func(arr []int) {
		_ = parallel.ParallelMergeSort(arr)
	})

	elapsedTimeParPar := time.Since(start)
	log.Printf("Tempo total (Tarefas em Paralelo, mergesort paralelo): %s\n", elapsedTimeParPar)

	// 3A. Método usando biblioteca + MergeSort sequencial
	{
		libInputs := cloneDatasets(datasets)
		start := time.Now()
		exec := lib.NewExecutor(workers, workers*2)
		for i := range libInputs {
			arr := libInputs[i]
			_ = exec.Execute(func() error {
				_ = seq.MergeSort(arr)
				return nil
			})
		}
		exec.Shutdown()
		elapsedLibA := time.Since(start)
		log.Printf("Tempo total (Biblioteca 3A, mergesort seq): %s\n", elapsedLibA)
	}

	// 3B. Método usando biblioteca + MergeSort paralelo (paralelismo interno)
	{
		libInputs := cloneDatasets(datasets)
		start := time.Now()
		exec := lib.NewExecutor(workers, workers*2)
		for i := range libInputs {
			arr := libInputs[i]
			_ = exec.Execute(func() error {
				_ = parallel.ParallelMergeSort(arr)
				return nil
			})
		}
		exec.Shutdown()
		elapsedLibB := time.Since(start)
		log.Printf("Tempo total (Biblioteca 3B, mergesort paralelo): %s\n", elapsedLibB)
	}
}

// Criação de uma deep copy de uma matriz de inteiros.
// Ou seja, uma matriz idêntica, porém com endereços diferentes.
func cloneDatasets(src [][]int) [][]int {
	out := make([][]int, len(src))
	for i := range src {
		out[i] = make([]int, len(src[i]))
		copy(out[i], src[i])
	}
	return out
}

// runTasksInParallel: pool simples controlado por 'workers'
func runTasksInParallel(workers int, datasets [][]int, work func([]int)) {
	jobs := make(chan []int, len(datasets))
	var wg sync.WaitGroup

	// Workers
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for arr := range jobs {
				work(arr)
			}
		}()
	}

	// Enfileira jobs (tarefas)
	for i := range datasets {
		jobs <- datasets[i]
	}
	close(jobs)

	wg.Wait()
}
