package parallel

import (
	"paralela/cmd/seq"
	"runtime"
	"sync"
)

const minChunk = 1 << 14 // = 2^14 = 16384. Ajuste com base no seu hardware.

// depthOrDefault define quantos níveis paralelos abrir baseado
// no número de CPUs lógicas que o Go está usando.
// Regra prática: ~2 * log2(NCPU).
func depthOrDefault(d int) int {
	if d > 0 {
		return d
	}
	n := runtime.GOMAXPROCS(0) //Não altera o número de threads máximas que podem executar goroutines simultaneamente.

	//log2 aproximado por contagem de bits.
	depth := 0
	for n > 1 {
		n >>= 1 //bit shift a direita (n = n/2, mas atuando com bits)
		depth++
	}

	return 2*depth + 1
}

func Merge(left, right []int) []int {
	// Cria o slice de resultado com capacidade total das duas fatias de entrada.
	result := make([]int, 0, len(left)+len(right))

	i, j := 0, 0 // i para 'left', j para 'right'

	// Itera enquanto houver elementos em ambas as fatias
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Adiciona os elementos restantes de 'left' (se houver)
	for ; i < len(left); i++ {
		result = append(result, left[i])
	}

	// Adiciona os elementos restantes de 'right' (se houver)
	for ; j < len(right); j++ {
		result = append(result, right[j])
	}

	return result
}

func ParallelMergeSort(arr []int) []int {
	return ParallelMergeSortDepth(arr, depthOrDefault(0))
}

func ParallelMergeSortDepth(arr []int, depth int) []int {
	// Caso base da recursão: uma fatia com 0 ou 1 elemento está sempre ordenada.
	if len(arr) <= 1 {
		return arr
	}

	// Se o subproblema for pequeno ou profundidade esgotou, cai para sequencial.
	if len(arr) < minChunk || depth <= 0 {
		return seq.MergeSort(arr)
	}

	mid := len(arr) / 2
	var left, right []int

	// Semáforo
	// Sincronizar (wait) goroutines
	var wg sync.WaitGroup

	// Adiciona 1 goroutine ao semafaro
	wg.Add(1)
	go func() {
		// Garante que ira reduzir uma goroutine ao final da funcao
		defer wg.Done()
		left = ParallelMergeSortDepth(arr[:mid], depth-1)
	}()

	// Paralela ramo direito se ainda tiver depth
	right = ParallelMergeSortDepth(arr[mid:], depth-1)

	// Espera semafaro chegar a 0. Garantindo sincronizacao
	wg.Wait()

	return Merge(left, right)
}
