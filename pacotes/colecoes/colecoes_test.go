package colecoes

import (
	"testing"
)

func TestPilhaGenerica(t *testing.T) {
	p := NovaPilha[string]()
	if !p.Vazia() || p.Tamanho() != 0 {
		t.Errorf("Pilha recém-criada deveria estar vazia")
	}

	p.Empilhar("VAMOS")
	p.Empilhar("LANG")

	if p.Tamanho() != 2 {
		t.Errorf("Tamanho esperado 2, obteve %d", p.Tamanho())
	}

	topo, ok := p.Topo()
	if !ok || topo != "LANG" {
		t.Errorf("Topo esperado 'LANG', obteve '%s'", topo)
	}

	item1, ok1 := p.Desempilhar()
	if !ok1 || item1 != "LANG" {
		t.Errorf("Desempilhar esperado 'LANG', obteve '%s'", item1)
	}

	item2, ok2 := p.Desempilhar()
	if !ok2 || item2 != "VAMOS" {
		t.Errorf("Desempilhar esperado 'VAMOS', obteve '%s'", item2)
	}

	if !p.Vazia() {
		t.Errorf("Pilha deveria estar vazia após desempilhar tudo")
	}
}

func TestFilaGenerica(t *testing.T) {
	f := NovaFila[int]()
	if !f.Vazia() {
		t.Errorf("Fila recém-criada deveria estar vazia")
	}

	f.Enfileirar(10)
	f.Enfileirar(20)
	f.Enfileirar(30)

	if f.Tamanho() != 3 {
		t.Errorf("Tamanho esperado 3, obteve %d", f.Tamanho())
	}

	frente, ok := f.Frente()
	if !ok || frente != 10 {
		t.Errorf("Frente esperada 10, obteve %d", frente)
	}

	it1, _ := f.Desenfileirar()
	it2, _ := f.Desenfileirar()
	it3, _ := f.Desenfileirar()

	if it1 != 10 || it2 != 20 || it3 != 30 {
		t.Errorf("Ordem incorreta na fila: %d, %d, %d", it1, it2, it3)
	}

	if !f.Vazia() {
		t.Errorf("Fila deveria estar vazia")
	}
}

func TestConjuntoGenerico(t *testing.T) {
	c := NovoConjunto[string]()
	c.Inserir("Go")
	c.Add("VAMOS")
	c.Inserir("Go") // Duplicado

	if c.Tamanho() != 2 {
		t.Errorf("Conjunto com valor duplicado deveria ter tamanho 2, obteve %d", c.Tamanho())
	}

	if !c.Contem("Go") || !c.Contem("VAMOS") {
		t.Errorf("Conjunto deveria conter 'Go' e 'VAMOS'")
	}
	if c.Contem("Java") {
		t.Errorf("Conjunto não deveria conter 'Java'")
	}

	c.Remover("Go")
	if c.Contem("Go") || c.Tamanho() != 1 {
		t.Errorf("Elemento 'Go' deveria ter sido removido")
	}
}
