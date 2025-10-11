package main

import (
	"fmt"
	"paralela/cmd/parallel"
	"paralela/cmd/seq"
	"paralela/cmd/util"
	"runtime"
	"sync"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Definições para o teste fixo
	const (
		ARRAY_SIZE = 1000000
		NUM_TASKS  = 16 // Número de vetores = tarefas
		FIXED_SEED = 42
	)

	// Gera o conjunto de tarefas: vetores independentes.
	datasets := make([][]int, NUM_TASKS)
	for i := 0; i < NUM_TASKS; i++ {
		// Seeds diferentes por tarefa, mas determinísticas.
		datasets[i] = util.GenerateDeterministicArray(ARRAY_SIZE, int64(FIXED_SEED+int64(i)))
	}

	// Método de processamento sequencial
	seqInputs := cloneDatasets(datasets)
	start := time.Now()

	for i := 0; i < NUM_TASKS; i++ {
		_ = seq.MergeSort(seqInputs[i])
	}

	elapsedTime := time.Since(start)
	fmt.Printf("Tempo de execução (Tarefas sequênciais, merge sort seq): %s\n", elapsedTime)

	// Método de processamento paralelo (Multithread de tarefas e merge sort sequencial)
	parInputs := cloneDatasets(datasets)
	start = time.Now()

	runTasksInParallel(parInputs, func(arr []int) {
		_ = seq.MergeSort(arr)
	})

	elapsedTime = time.Since(start)

	fmt.Printf("Tempo total (Tarefas em Paralelo, mergesort seq): %s\n", elapsedTime)

	// Método de processamento paralelo (Multithread de tarefas e merge sort sequencial)

	parAlgoInputs := cloneDatasets(datasets)
	start = time.Now()

	runTasksInParallel(parAlgoInputs, func(arr []int) {
		_ = parallel.ParallelMergeSort(arr)
	})

	elapsedTime = time.Since(start)
	fmt.Printf("Tempo total (Tarefas em Paralelo, mergesort paralelo): %s\n", elapsedTime)
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

func runTasksInParallel(datasets [][]int, work func([]int)) {
	workers := runtime.GOMAXPROCS(0) //usa NCPU (CPU lógicas) workers (goroutines)
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
