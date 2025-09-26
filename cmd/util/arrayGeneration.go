package util

import "math/rand"

func GenerateDeterministicArray(size int, seed int64) []int {
	// Cria uma nova fonte de aleatoriedade com a seed fixa.
	source := rand.NewSource(seed)
	r := rand.New(source)
	
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		// Gera números grandes para aumentar o custo de comparação e embaralhar bem
		arr[i] = r.Intn(1000000) 
	}
	return arr
}
