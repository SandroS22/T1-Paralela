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

	workers := runtime.GOMAXPROCS(runtime.NumCPU())

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
	var seqResult [][]int

	for i := 0; i < NUM_TASKS; i++ {
		seqResult = append(seqResult, seq.MergeSort(seqInputs[i]))
	}

	elapsedTime := time.Since(start)
	log.Printf("Tempo de execução (Tarefas sequênciais): %s\n", elapsedTime)
	if(!allInnerSlicesSorted(seqResult)) {
		log.Println("erro: Ha vetores nao organizados! (Tarefas sequênciais)")
	}

	// Método de processamento paralelo (Multithread de tarefas e merge sort sequencial)
	parInputs := cloneDatasets(datasets)
	var parResult [][]int
	start = time.Now()

	runTasksInParallel(parInputs, func(arr []int) {
		parResult = append(parResult, seq.MergeSort(arr))
	})

	elapsedTime = time.Since(start)

	log.Printf("Tempo total (Tarefas em Paralelo, mergesort seq): %s\n", elapsedTime)
	if(!allInnerSlicesSorted(parResult)) {
		log.Println("erro: Ha vetores nao organizados! (Tarefas em Paralelo, mergesort seq)")
	}

	// Método de processamento paralelo (Multithread de tarefas e merge sort sequencial)
	parAlgoInputs := cloneDatasets(datasets)
	var parAlgoInputsResult [][]int
	start = time.Now()

	runTasksInParallel(parAlgoInputs, func(arr []int) {
		parResult = append(parAlgoInputsResult, parallel.ParallelMergeSort(arr))
	})

	elapsedTime = time.Since(start)
	log.Printf("Tempo total (Tarefas em Paralelo, mergesort paralelo): %s\n", elapsedTime)
	if(!allInnerSlicesSorted(parAlgoInputsResult)) {
		log.Println("erro: Ha vetores nao organizados! (Tarefas em Paralelo, mergesort seq)")
	}

	// crie o executor antes do loop:
	start = time.Now()
	exec := lib.NewExecutor(workers, workers*2) // fila 2x workers

	// submeta cada vetor como uma tarefa
	for i := range datasets {
		arr := datasets[i] // captura
		_ = exec.Execute(func() error {
			_ = seq.MergeSort(arr)
			return nil
		})
	}
	// finalize
	exec.Shutdown()
	elapsedTime = time.Since(start)
	log.Printf("Tempo total (Biblioteca): %s\n", elapsedTime)
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

func isSorted(slice []int) bool {
	for i := 1; i < len(slice); i++ {
		if slice[i] < slice[i-1] {
			return false
		}
	}
	return true
}

// Verifica se todos os slices internos estão ordenados
func allInnerSlicesSorted(matrix [][]int) bool {
	for _, inner := range matrix {
		if !isSorted(inner) {
			return false
		}
	}
	return true
}
