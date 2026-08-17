package transpiler

import (
	"strings"
	"testing"
)

func TestTranspileOlaMundo(t *testing.T) {
	input := `pacote principal

importar "formatar"

funcao principal() {
	formatar.ImprimirLinha("Olá, Mundo!")
}
`
	expected := `package main

import "fmt"

func main() {
	fmt.Println("Olá, Mundo!")
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro inesperado ao transpilar: %v", err)
	}

	if got != expected {
		t.Errorf("Código transpilado incorreto.\nEsperado:\n%s\nObtido:\n%s", expected, got)
	}
}

func TestTranspileKeywords(t *testing.T) {
	input := `pacote principal

importar "formatar"

tipo Usuario estrutura {
	Nome texto
	Idade inteiro
	Ativo booleano
}

tipo Processador interface {
	Executar() erro
}

constante Maximo inteiro = 100
variavel Contador inteiro = 0

funcao checar(x inteiro) texto {
	se x > 10 {
		retornar "maior"
	} senao {
		retornar "menor ou igual"
	}
}

funcao iterar() {
	para i := 0; i < 10; i++ {
		se i == 5 {
			continuar
		}
		se i == 8 {
			interromper
		}
	}
}

funcao testarSwitch(valor inteiro) {
	escolha valor {
	caso 1:
		prosseguir
	caso 2:
		retornar
	padrao:
		retornar
	}
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar palavras-chave: %v", err)
	}

	keywordsToCheck := []string{
		"package main",
		"type Usuario struct",
		"type Processador interface",
		"const Maximo int = 100",
		"var Contador int = 0",
		"func checar(x int) string",
		"if x > 10",
		"else {",
		"return \"maior\"",
		"for i := 0; i < 10; i++",
		"continue",
		"break",
		"switch valor {",
		"case 1:",
		"fallthrough",
		"case 2:",
		"default:",
	}

	for _, kw := range keywordsToCheck {
		if !strings.Contains(got, kw) {
			t.Errorf("Esperava encontrar %q no código gerado, mas não encontrou:\n%s", kw, got)
		}
	}
}

func TestTranspileStringsAndCommentsPreservation(t *testing.T) {
	input := `pacote principal

importar "formatar"

// Este comentário tem funcao, variavel, pacote e retornar!
/*
   Comentário em bloco com
   palavras: se senao para criar canal
*/
funcao principal() {
	msg := "A string contem: pacote principal funcao variavel se senao retornar"
	raw := ` + "`" + `Raw string com:
	tipo estrutura interface
	para intervalo
	` + "`" + `
	formatar.ImprimirLinha(msg, raw)
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar: %v", err)
	}

	// Comentários e strings literais não devem ter palavras-chave traduzidas
	if !strings.Contains(got, "// Este comentário tem funcao, variavel, pacote e retornar!") {
		t.Errorf("Comentário de linha foi alterado incorretamente")
	}
	if !strings.Contains(got, "se senao para criar canal") {
		t.Errorf("Comentário de bloco foi alterado incorretamente")
	}
	if !strings.Contains(got, `"A string contem: pacote principal funcao variavel se senao retornar"`) {
		t.Errorf("String literal interpretada foi alterada incorretamente")
	}
	if !strings.Contains(got, "tipo estrutura interface") {
		t.Errorf("String literal raw foi alterada incorretamente")
	}
}

func TestTranspileConcurrency(t *testing.T) {
	input := `pacote principal

importar (
	"formatar"
	"tempo"
	"sincronizar"
)

funcao trabalhador(id inteiro, c canal texto, wg *sincronizar.GrupoEspera) {
	adiar wg.Concluido()
	tempo.Dormir(tempo.Milissegundo * 50)
	c <- formatar.Formatar("Trabalhador %d pronto", id)
}

funcao principal() {
	canalMsg := criar(canal texto, 2)
	variavel wg sincronizar.GrupoEspera

	wg.Adicionar(2)
	disparar trabalhador(1, canalMsg, &wg)
	disparar trabalhador(2, canalMsg, &wg)

	wg.Esperar()
	fechar(canalMsg)

	para msg := intervalo canalMsg {
		formatar.ImprimirLinha(msg)
	}
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar concorrência: %v", err)
	}

	expectedPhrases := []string{
		`import (`,
		`"fmt"`,
		`"time"`,
		`"sync"`,
		`func trabalhador(id int, c chan string, wg *sync.WaitGroup)`,
		`defer wg.Done()`,
		`time.Sleep(time.Millisecond * 50)`,
		`c <- fmt.Sprintf("Trabalhador %d pronto", id)`,
		`canalMsg := make(chan string, 2)`,
		`var wg sync.WaitGroup`,
		`wg.Add(2)`,
		`go trabalhador(1, canalMsg, &wg)`,
		`wg.Wait()`,
		`close(canalMsg)`,
		`for msg := range canalMsg {`,
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(got, phrase) {
			t.Errorf("Esperava encontrar %q no código gerado, mas não encontrou:\n%s", phrase, got)
		}
	}
}

func TestLinePreservation(t *testing.T) {
	input := `pacote principal

importar (
	"formatar"
	"tempo"
)

// Linha de comentário 1
// Linha de comentário 2

funcao somar(a inteiro, b inteiro) inteiro {
	retornar a + b
}

funcao principal() {
	x := 10
	y := 20
	res := somar(x, y)
	formatar.ImprimirLinha("Resultado:", res)
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar: %v", err)
	}

	inputLines := strings.Count(input, "\n")
	gotLines := strings.Count(got, "\n")

	if inputLines != gotLines {
		t.Errorf("Número de linhas não preservado: esperado %d linhas, obteve %d linhas", inputLines, gotLines)
	}
}

func TestTranspileMapsSlicesAndBuiltins(t *testing.T) {
	input := `pacote principal

importar "formatar"

funcao principal() {
	// Slices e Maps
	tabela := criar(mapa[texto]inteiro)
	tabela["um"] = 1
	tabela["dois"] = 2
	excluir(tabela, "um")

	lista := criar([]inteiro, 0, 10)
	lista = anexar(lista, 42)
	tam := tamanho(lista)
	capa := capacidade(lista)

	formatar.ImprimirLinha("Tamanho:", tam, "Capacidade:", capa)

	// Recuperação de Pânico
	adiar funcao() {
		se r := recuperar(); r != nulo {
			formatar.ImprimirLinha("Recuperado de:", r)
		}
	}()
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar maps e slices: %v", err)
	}

	expectedPhrases := []string{
		`tabela := make(map[string]int)`,
		`delete(tabela, "um")`,
		`lista := make([]int, 0, 10)`,
		`lista = append(lista, 42)`,
		`tam := len(lista)`,
		`capa := cap(lista)`,
		`defer func() {`,
		`if r := recover(); r != nil {`,
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(got, phrase) {
			t.Errorf("Esperava encontrar %q no código gerado, mas não encontrou:\n%s", phrase, got)
		}
	}
}

func TestTranspileSelectAndGoto(t *testing.T) {
	input := `pacote principal

funcao testarSelect(c1 canal inteiro, c2 canal texto) {
	selecionar {
	caso v := <-c1:
		imprimirln("Recebido c1:", v)
	caso msg := <-c2:
		imprimirln("Recebido c2:", msg)
	padrao:
		ir_para fim
	}

fim:
	retornar
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar select e goto: %v", err)
	}

	expectedPhrases := []string{
		`func testarSelect(c1 chan int, c2 chan string)`,
		`select {`,
		`case v := <-c1:`,
		`println("Recebido c1:", v)`,
		`case msg := <-c2:`,
		`println("Recebido c2:", msg)`,
		`default:`,
		`goto fim`,
		`fim:`,
		`return`,
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(got, phrase) {
			t.Errorf("Esperava encontrar %q no código gerado, mas não encontrou:\n%s", phrase, got)
		}
	}
}

