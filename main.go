package main

import (
	"fmt"
	"paralela/cmd"
	"paralela/cmd/util"
	"time"
)


func main(){
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
	sortedArray := cmd.MergeSort(arrayToSort)
	
	elapsed := time.Since(start)
	
	fmt.Printf("Ordenação Sequencial Concluída.\n")
	fmt.Printf("Tempo de Execução (Tempo Sequencial): %s\n", elapsed)
	
	if len(sortedArray) > 0 {
		fmt.Printf("Primeiro elemento: %d, Último elemento: %d\n", sortedArray[0], sortedArray[len(sortedArray)-1])
	}
}
