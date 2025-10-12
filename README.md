# T1-Paralela

### Objetivo: Criar um algoritmo sequencial e paralelo. Além disso criar uma biblioteca própria de paralelismo e testa-la junto com as anteriores

## Contribuidores
* Felipe Delduqui
* Sandro Santana Ribeiro

## Soluções

### Solução Sequencial

Cada vetor é ordenado individualmente, um de cada vez, utilizando o algoritmo MergeSort.

Serve como base de comparação para as versões paralelas.

### Solução Paralela 1 (Entre tarefas)

O programa executa várias tarefas independentes em paralelo, onde cada tarefa é a ordenação de um vetor.

As tarefas são distribuídas a um pool de threads (workers) equivalente ao número de núcleos da CPU.

Cada tarefa usa o MergeSort sequencial internamente, garantindo isonomia no teste de paralelismo entre tarefas.

### Solução Paralela 2 (Dentro da tarefa)

Dentro de cada tarefa, usou-se um MergeSort paralelo

Essa versão cria goroutines recursivamente até um limite controlado de profundidade e tamanho mínimo de partição, evitando overhead excessivo.

Implementada em cmd/parallel/parallelMergeSort.go.

### Solução Paralela 3 (Biblioteca própria)

Foi desenvolvida uma biblioteca genérica de execução paralela, localizada em cmd/lib/executor.go.

A biblioteca implementa um pool de threads e o padrão Produtor/Consumidor.

## Padrões de Projeto Aplicados

* Worker Pool / Produtor–Consumidor – gerenciamento das tarefas paralelas.

* Fork–Join – usado no MergeSort paralelo.

* Divide and Conquer – estrutura de recursão e fusão do MergeSort.

* Strategy – escolha entre diferentes modos de execução (sequencial, paralelo 1, paralelo 2).

## Compilação
Para compilação, utilize os seguintes códigos no terminal:

```bash
$ go build main.go
$ ./main
```

## Estrutura de pastas
As pastas foram estruturadas da seguinte maneira:
```text
paralela/
│
├── cmd/
│   ├── lib/
│   │   └── executor.go            # Biblioteca de execução paralela
│   ├── parallel/
│   │   └── parallelMergeSort.go   # MergeSort paralelo com limitação de profundidade
│   ├── seq/
│   │   └── sequentialMergeSort.go # MergeSort sequencial
│   └── util/
│       └── arrayGeneration.go     # Geração determinística de vetores
│
├── main.go                        # Execução e comparação das versões
└── go.mod                         # Configuração do módulo Go

```
