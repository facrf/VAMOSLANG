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

func TestTranspileMapsSlicesAndBuiltins(t *testing.T) {
	input := `pacote principal

importar "formatar"

funcao principal() {
	tabela := criar(mapa[texto]inteiro)
	tabela["um"] = 1
	excluir(tabela, "um")

	lista := criar([]inteiro, 0, 10)
	lista = adicionar(lista, 42)
	lista = anexar(lista, 99)
	tam := tamanho(lista)
	capa := capacidade(lista)

	formatar.ImprimirLinha("Tamanho:", tam, "Capacidade:", capa)

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
		`lista = append(lista, 99)`,
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

func TestTranspileExpandedStdlib(t *testing.T) {
	input := `pacote principal

importar (
	"leitor"
	"caminho"
	"ordenacao"
	"contexto"
	"so"
)

funcao testarStdlib() {
	scanner := leitor.NovoScanner(so.EntradaPadrao)
	_ = scanner
	caminhoCompleto := caminho.Juntar("dir", "arquivo.txt")
	_ = caminhoCompleto
	nomes := []texto{"Carlos", "Ana", "Beatriz"}
	ordenacao.Strings(nomes)
	ctx := contexto.PlanoDeFundo()
	_ = ctx
}
`
	got, err := TranspileSource(input)
	if err != nil {
		t.Fatalf("Erro ao transpilar biblioteca expandida: %v", err)
	}

	expectedPhrases := []string{
		`"bufio"`,
		`"path/filepath"`,
		`"sort"`,
		`"context"`,
		`"os"`,
		`scanner := bufio.NewScanner(os.Stdin)`,
		`caminhoCompleto := filepath.Join("dir", "arquivo.txt")`,
		`nomes := []string{"Carlos", "Ana", "Beatriz"}`,
		`sort.Strings(nomes)`,
		`ctx := context.Background()`,
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(got, phrase) {
			t.Errorf("Esperava encontrar %q no código gerado, mas não encontrou:\n%s", phrase, got)
		}
	}
}

func TestTranspileGoToVamos(t *testing.T) {
	input := `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Olá Go -> VAMOS")
	time.Sleep(time.Millisecond * 10)
}
`
	got, err := TranspileGoToVamos(input)
	if err != nil {
		t.Fatalf("Erro ao descompilar Go para VAMOS: %v", err)
	}

	expectedPhrases := []string{
		`pacote principal`,
		`importar (`,
		`"formatar"`,
		`"tempo"`,
		`funcao principal() {`,
		`formatar.ImprimirLinha("Olá Go -> VAMOS")`,
		`tempo.Dormir(tempo.Milissegundo * 10)`,
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(got, phrase) {
			t.Errorf("Esperava encontrar %q no código VAMOS gerado, mas não encontrou:\n%s", phrase, got)
		}
	}
}

func TestFormatSourceWithRawString(t *testing.T) {
	input := `pacote principal

funcao principal() {
banner := ` + "`" + `
{
  "chave": "valor"
}
` + "`" + `
}`
	formatted, err := FormatSource(input)
	if err != nil {
		t.Fatalf("Erro ao formatar: %v", err)
	}

	// Verifica se as linhas dentro do raw string preservaram suas posições e chaves
	if !strings.Contains(formatted, `  "chave": "valor"`) {
		t.Errorf("Formatador corrompeu o conteúdo de raw string multilinhas:\n%s", formatted)
	}
}

func TestLintSource(t *testing.T) {
	input := `pacote Principal

funcao teste() {
	_ = err
}
`
	issues := LintSource(input)
	if len(issues) < 2 {
		t.Errorf("Esperava pelo menos 2 avisos no linter, obteve %d", len(issues))
	}
}
