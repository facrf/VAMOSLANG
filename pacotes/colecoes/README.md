# 📦 Pacote `colecoes` (Biblioteca Padrão VAMOS-LANG)

O pacote **`colecoes`** faz parte da biblioteca padrão oficial em VAMOS-LANG puro, fornecendo estruturas de dados clássicas com suporte total a **tipos genéricos** (`[T qualquer]` e `[T comparavel]`).

---

## 📚 Estruturas Disponíveis

### 1. `Pilha[T qualquer]` (LIFO - Last In, First Out)
Implementação de pilha encadeada/vetorial genérica.

- **`NovaPilha[T]() *Pilha[T]`**: Cria uma pilha vazia.
- **`Empilhar(item T)`**: Insere um elemento no topo da pilha.
- **`Desempilhar() (T, booleano)`**: Remove e retorna o elemento do topo (ou falso se vazia).
- **`Topo() (T, booleano)`**: Retorna o elemento do topo sem removê-lo.
- **`Tamanho() inteiro`**: Retorna a quantidade de elementos.
- **`Vazia() booleano`**: Retorna `verdadeiro` se a pilha não possuir itens.

```vamos
importar "colecoes"

pilha := colecoes.NovaPilha[texto]()
pilha.Empilhar("Primeiro")
pilha.Empilhar("Segundo")

se topo, ok := pilha.Desempilhar(); ok {
    formatar.ImprimirLinha(topo) // Segundo
}
```

---

### 2. `Fila[T qualquer]` (FIFO - First In, First Out)
Implementação de fila genérica de atendimento sequencial.

- **`NovaFila[T]() *Fila[T]`**: Cria uma fila vazia.
- **`Enfileirar(item T)`**: Insere um elemento no fim da fila.
- **`Desenfileirar() (T, booleano)`**: Remove e retorna o primeiro elemento (ou falso se vazia).
- **`Frente() (T, booleano)`**: Retorna o primeiro elemento sem removê-lo.
- **`Tamanho() inteiro`**: Retorna a quantidade de elementos na fila.
- **`Vazia() booleano`**: Retorna `verdadeiro` se a fila estiver vazia.

```vamos
importar "colecoes"

fila := colecoes.NovaFila[inteiro]()
fila.Enfileirar(10)
fila.Enfileirar(20)

se item, ok := fila.Desenfileirar(); ok {
    formatar.ImprimirLinha(item) // 10
}
```

---

### 3. `Conjunto[T comparavel]` (Valores Únicos / Set)
Implementação de conjunto de valores matemáticos sem duplicatas.

- **`NovoConjunto[T]() *Conjunto[T]`**: Inicializa um conjunto vazio.
- **`Inserir(item T)` / `Adicionar(item T)`**: Adiciona um item ao conjunto (ignora duplicatas).
- **`Remover(item T)`**: Remove o elemento do conjunto.
- **`Contem(item T) booleano`**: Verifica se o item pertence ao conjunto.
- **`Tamanho() inteiro`**: Retorna o total de elementos únicos.
- **`ParaLista() []T`**: Exporta os itens do conjunto para uma fatia/lista.

```vamos
importar "colecoes"

conj := colecoes.NovoConjunto[texto]()
conj.Inserir("Go")
conj.Inserir("VAMOS")
conj.Inserir("Go") // Ignorado

formatar.ImprimirFormatado("Total de itens: %d\n", conj.Tamanho()) // 2
formatar.ImprimirFormatado("Contém 'VAMOS'? %v\n", conj.Contem("VAMOS")) // true
```
