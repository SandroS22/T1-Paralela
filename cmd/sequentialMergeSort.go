package cmd


// Merge combina dois slices já ordenadas em uma única fatia ordenada.
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

// MergeSort é a função recursiva que implementa a lógica de "divisão e conquista".
func MergeSort(arr []int) []int {
	// Caso base da recursão: uma fatia com 0 ou 1 elemento está sempre ordenada.
	if len(arr) <= 1 {
		return arr
	}
	
	mid := len(arr) / 2
	
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])
	
	return Merge(left, right)
}
