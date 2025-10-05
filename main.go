package main

import (
	"fmt"
	"paralela/cmd/paralel"
	"paralela/cmd/seq"
	"paralela/cmd/util"
	"runtime"
	"time"
)


func main(){

	runtime.GOMAXPROCS(runtime.NumCPU())
	// Definições para o teste fixo
	const ARRAY_SIZE = 1000000 
	const FIXED_SEED = 42
	
	data := util.GenerateDeterministicArray(ARRAY_SIZE, FIXED_SEED)
	
	// Criar uma cópia do array para a ordenação, mantendo o original intacto se necessário
	arrayToSort := make([]int, ARRAY_SIZE)
	copy(arrayToSort, data)

	// Medição de Tempo
	start := time.Now()
	// Executa a ordenação sequencial
	sortedArray := seq.MergeSort(arrayToSort)
	elapsedTime := time.Since(start)
	fmt.Printf("Tempo de execução (Tempo Sequencial): %s\n", elapsedTime)

	if len(sortedArray) > 0 {
		fmt.Printf("Primeiro elemento: %d, Último elemento: %d\n", sortedArray[0], sortedArray[len(sortedArray)-1])
	}

	start = time.Now()
	sortedArray = paralel.ParalelMergeSort(arrayToSort)
	elapsedTime = time.Since(start)

	
	fmt.Printf("Tempo de execução (Tempo Paralelo): %s\n", elapsedTime)
	
	if len(sortedArray) > 0 {
		fmt.Printf("Primeiro elemento: %d, Último elemento: %d\n", sortedArray[0], sortedArray[len(sortedArray)-1])
	}
}
