# T1-Paralela

## Criar um algoritmo sequencial e paralelo. Além disso criar uma biblioteca própria de paralelismo e testa-la junto com as anteriores

### Contribuidores
* Sandro Santana Ribeiro
* Felipe Delduqui

### Solução Sequencial

Cada vetor é ordenado individualmente, um de cada vez, utilizando o algoritmo MergeSort.

Serve como base de comparação para as versões paralelas.

### Solução Paralela 1

O programa executa várias tarefas independentes em paralelo, onde cada tarefa é a ordenação de um vetor.

As tarefas são distribuídas a um pool de threads (workers) equivalente ao número de núcleos da CPU.

Cada tarefa usa o MergeSort sequencial internamente, garantindo isonomia no teste de paralelismo entre tarefas.

Código principal em main.go.

### Solução Paralela 2

Dentro de cada tarefa, usou-se um MergeSort paralelo

Essa versão cria goroutines recursivamente até um limite controlado de profundidade e tamanho mínimo de partição, evitando overhead excessivo.

### Compilação
Para compilação, utilize os seguintes códigos no terminal:
```bash
$ go build main.go
$ ./main
```

### Estrutura de pastas
```bash
paralela/
│
├── cmd/
│   ├── parallel/
│   │   └── parallelMergeSort.go   # MergeSort paralelo com limitação de profundidade
│   ├── seq/
│   │   └── sequentialMergeSort.go # MergeSort sequencial
│   └── util/
│       └── arrayGeneration.go     # Geração determinística de vetores
│
├── main.go                        # Execução das comparações de desempenho
└── go.mod                         # Configuração do módulo Go
```
