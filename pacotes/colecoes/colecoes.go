// Pacote colecoes - Estruturas de dados genéricas em VAMOS-LANG
// Fornece implementações de Pilha, Fila e Conjunto com tipos genéricos.

package colecoes

// Pilha representa uma estrutura de dados LIFO (Last-In, First-Out)
type Pilha[T any] struct {
	itens []T
}

// NovaPilha inicializa uma pilha vazia
func NovaPilha[T any]() *Pilha[T] {
	return &Pilha[T]{
		itens: make([]T, 0),
	}
}

// Empilhar insere um elemento no topo da pilha
func (p *Pilha[T]) Empilhar(item T) {
	p.itens = append(p.itens, item)
}

// Desempilhar remove e retorna o elemento do topo da pilha
func (p *Pilha[T]) Desempilhar() (T, bool) {
	if p.Vazia() {
		var zero T
		return zero, false
	}
	ultimoIndice := len(p.itens) - 1
	item := p.itens[ultimoIndice]
	p.itens = p.itens[:ultimoIndice]
	return item, true
}

// Topo retorna o elemento do topo da pilha sem removê-lo
func (p *Pilha[T]) Topo() (T, bool) {
	if p.Vazia() {
		var zero T
		return zero, false
	}
	return p.itens[len(p.itens)-1], true
}

// Tamanho retorna a quantidade de itens na pilha
func (p *Pilha[T]) Tamanho() int {
	return len(p.itens)
}

// Vazia verifica se a pilha não possui elementos
func (p *Pilha[T]) Vazia() bool {
	return len(p.itens) == 0
}

// Fila representa uma estrutura de dados FIFO (First-In, First-Out)
type Fila[T any] struct {
	itens []T
}

// NovaFila inicializa uma fila vazia
func NovaFila[T any]() *Fila[T] {
	return &Fila[T]{
		itens: make([]T, 0),
	}
}

// Enfileirar insere um elemento no fim da fila
func (f *Fila[T]) Enfileirar(item T) {
	f.itens = append(f.itens, item)
}

// Desenfileirar remove e retorna o primeiro elemento da fila
func (f *Fila[T]) Desenfileirar() (T, bool) {
	if f.Vazia() {
		var zero T
		return zero, false
	}
	item := f.itens[0]
	f.itens = f.itens[1:]
	return item, true
}

// Frente retorna o primeiro elemento da fila sem removê-lo
func (f *Fila[T]) Frente() (T, bool) {
	if f.Vazia() {
		var zero T
		return zero, false
	}
	return f.itens[0], true
}

// Tamanho retorna a quantidade de elementos na fila
func (f *Fila[T]) Tamanho() int {
	return len(f.itens)
}

// Vazia verifica se a fila está vazia
func (f *Fila[T]) Vazia() bool {
	return len(f.itens) == 0
}

// Conjunto representa uma coleção de valores únicos
type Conjunto[T comparable] struct {
	elementos map[T]bool
}

// NovoConjunto inicializa um conjunto vazio
func NovoConjunto[T comparable]() *Conjunto[T] {
	return &Conjunto[T]{
		elementos: make(map[T]bool),
	}
}

// Inserir adiciona um valor no conjunto
func (c *Conjunto[T]) Inserir(item T) {
	c.elementos[item] = true
}

// Add adiciona um valor no conjunto (compatibilidade direta com Adicionar transpilado)
func (c *Conjunto[T]) Add(item T) {
	c.elementos[item] = true
}

// Remover elimina um valor do conjunto
func (c *Conjunto[T]) Remover(item T) {
	delete(c.elementos, item)
}

// Contem verifica se o elemento pertence ao conjunto
func (c *Conjunto[T]) Contem(item T) bool {
	_, existe := c.elementos[item]
	return existe
}

// Tamanho retorna o total de elementos únicos
func (c *Conjunto[T]) Tamanho() int {
	return len(c.elementos)
}

// ParaLista exporta os elementos do conjunto para uma lista
func (c *Conjunto[T]) ParaLista() []T {
	lista := make([]T, 0)
	for k := range c.elementos {
		lista = append(lista, k)
	}
	return lista
}
