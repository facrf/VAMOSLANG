package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTranslateCompilerMessages(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "undefined: minhaVariavel",
			expected: "identificador não definido: minhaVariavel",
		},
		{
			input:    "cannot use x as value in return statement",
			expected: "não é possível usar x como valor no retorno",
		},
		{
			input:    "too many arguments in call to funcao",
			expected: "muitos argumentos na chamada para funcao",
		},
		{
			input:    "not enough arguments in call to funcao",
			expected: "argumentos insuficientes na chamada para funcao",
		},
		{
			input:    "declared and not used: total",
			expected: "declarado e não utilizado: total",
		},
		{
			input:    "panic: runtime error: index out of range [5] with length 3",
			expected: "pânico (erro em tempo de execução): índice fora dos limites [5] with length 3",
		},
		{
			input:    "nil pointer dereference",
			expected: "desreferência de ponteiro nulo",
		},
	}

	for _, tt := range tests {
		got := TranslateCompilerMessages(tt.input)
		if got != tt.expected {
			t.Errorf("TranslateCompilerMessages(%q) = %q, esperado %q", tt.input, got, tt.expected)
		}
	}
}

func TestRewriteFileErrors(t *testing.T) {
	tempPath := "/tmp/vamos_run_123456.go"
	origPath := "exemplos/ola_mundo.vamos"
	rawError := "/tmp/vamos_run_123456.go:12:5: undefined: x\n./vamos_run_123456.go:15:1: syntax error: unexpected token"

	rewritten := RewriteFileErrors(rawError, tempPath, origPath)

	if strings.Contains(rewritten, "vamos_run_123456") {
		t.Errorf("Erro reescrito ainda contém caminho temporário: %s", rewritten)
	}

	if !strings.Contains(rewritten, "exemplos/ola_mundo.vamos:12:5:") {
		t.Errorf("Erro reescrito não contém o caminho original correto: %s", rewritten)
	}

	if !strings.Contains(rewritten, "identificador não definido: x") {
		t.Errorf("Mensagem de erro não foi traduzida para Português: %s", rewritten)
	}

	if !strings.Contains(rewritten, "erro de sintaxe:") {
		t.Errorf("Erro de sintaxe não foi traduzido para Português: %s", rewritten)
	}
}

func TestRewriteErrors(t *testing.T) {
	tempDir := "/tmp/vamos_workspace_98765"
	rawOutput := "/tmp/vamos_workspace_98765/main.go:20:10: declared and not used: y"

	rewritten := RewriteErrors(rawOutput, tempDir, ".")

	if strings.Contains(rewritten, tempDir) {
		t.Errorf("Saída ainda contém tempDir: %s", rewritten)
	}

	if !strings.Contains(rewritten, "main.vamos:20:10:") {
		t.Errorf("Extensão .go: não foi substituída por .vamos: %s", rewritten)
	}

	if !strings.Contains(rewritten, "declarado e não utilizado: y") {
		t.Errorf("Mensagem não foi traduzida: %s", rewritten)
	}
}

func TestFindProjectRoot(t *testing.T) {
	root := FindProjectRoot("exemplos/ola_mundo.vamos")
	if _, err := os.Stat(filepath.Join(root, "go.mod")); os.IsNotExist(err) {
		t.Errorf("FindProjectRoot('exemplos/ola_mundo.vamos') retornou %q que não possui go.mod", root)
	}
}

func TestPrepareWorkspace(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Erro ao obter diretório atual: %v", err)
	}

	tempDir, cleanup, err := PrepareWorkspace(cwd)
	if err != nil {
		t.Fatalf("PrepareWorkspace falhou: %v", err)
	}
	defer cleanup()

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("Diretório temporário do workspace não foi criado: %s", tempDir)
	}

	// Verifica se go.mod foi copiado para o diretório temporário
	if _, err := os.Stat(filepath.Join(tempDir, "go.mod")); os.IsNotExist(err) {
		t.Errorf("go.mod não encontrado no workspace temporário")
	}
}

func TestRunVMAndBytecode(t *testing.T) {
	root := FindProjectRoot(".")
	exemploPath := filepath.Join(root, "exemplos", "maquina_virtual.vamos")
	codeVM := RunVM([]string{exemploPath})
	if codeVM != 0 {
		t.Errorf("RunVM retornou código de erro: %d", codeVM)
	}

	codeBytecode := RunBytecode([]string{exemploPath})
	if codeBytecode != 0 {
		t.Errorf("RunBytecode retornou código de erro: %d", codeBytecode)
	}
}

func TestRunWasm(t *testing.T) {
	root := FindProjectRoot(".")
	exemploPath := filepath.Join(root, "exemplos", "wasm_modulo.vamos")
	saidaWasm := filepath.Join(os.TempDir(), "teste_wasm.wasm")
	defer os.Remove(saidaWasm)

	codeWasm := RunWasm([]string{exemploPath, "-o", saidaWasm})
	if codeWasm != 0 {
		t.Errorf("RunWasm retornou código de erro: %d", codeWasm)
	}

	dados, err := os.ReadFile(saidaWasm)
	if err != nil {
		t.Fatalf("Arquivo .wasm não foi gerado: %v", err)
	}
	if len(dados) < 8 || string(dados[:4]) != "\x00asm" {
		t.Errorf("Módulo WebAssembly inválido ou cabeçalho incorreto")
	}
}
