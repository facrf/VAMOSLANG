package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"vamos/transpiler"
)

// RunTrace instrumenta o código VAMOS-LANG para exibir o rastreamento didático de cada instrução executada.
func RunTrace(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos especificado para rastreamento.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos trace <arquivo.vamos> [argumentos...]\n")
		fmt.Fprintf(os.Stderr, "Exemplo: vamos trace exemplos/calculadora.vamos\n")
		return 1
	}

	srcFile := args[0]
	progArgs := args[1:]

	content, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcFile, err)
		return 1
	}

	instrumentedSource := InstrumentSourceForTrace(string(content))

	// Transpila o código instrumentado
	goCode, err := transpiler.TranspileSource(instrumentedSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro de transpilação no rastreio: %v\n", err)
		return 1
	}

	tempDir, err := os.MkdirTemp("", "vamos_trace_*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar diretório temporário: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tempDir)

	goFilePath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(goFilePath, []byte(goCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar arquivo temporário: %v\n", err)
		return 1
	}

	// Copia go.mod e go.sum se existirem no projeto original
	projectRoot := FindProjectRoot(srcFile)
	if modData, err := os.ReadFile(filepath.Join(projectRoot, "go.mod")); err == nil {
		_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), modData, 0644)
	} else {
		_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module trace_app\n\ngo 1.22\n"), 0644)
	}
	if sumData, err := os.ReadFile(filepath.Join(projectRoot, "go.sum")); err == nil {
		_ = os.WriteFile(filepath.Join(tempDir, "go.sum"), sumData, 0644)
	}

	fmt.Println("==================================================================")
	fmt.Printf(" 🔍 RASTREAMENTO DIDÁTICO DE EXECUÇÃO: %s\n", filepath.Base(srcFile))
	fmt.Println(" Visualização passo a passo de instruções e fluxo de controle")
	fmt.Println("==================================================================")

	cmdArgs := append([]string{"run", "main.go"}, progArgs...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = tempDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	execErr := cmd.Run()
	fmt.Println("\n==================================================================")
	if execErr != nil {
		fmt.Println(" ❌ Execução finalizada com erro.")
		return 1
	}
	fmt.Println(" ✓ Execução e rastreamento finalizados com sucesso!")
	fmt.Println("==================================================================")
	return 0
}

// InstrumentSourceForTrace injeta chamadas de log visual antes de cada declaração dentro de blocos de função.
func InstrumentSourceForTrace(source string) string {
	lines := strings.Split(source, "\n")
	var sb strings.Builder

	inFunction := false
	inBlock := 0
	inComposite := 0
	hasFormatImport := strings.Contains(source, `"formatar"`)

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Garante import de "formatar" se o arquivo não tiver
		if strings.HasPrefix(trimmed, "pacote ") {
			sb.WriteString(line + "\n")
			if !hasFormatImport {
				sb.WriteString("importar \"formatar\"\n")
			}
			continue
		}

		if strings.HasPrefix(trimmed, "funcao ") && strings.Contains(trimmed, "{") {
			inFunction = true
			inBlock++
			sb.WriteString(line + "\n")
			continue
		}

		// Detecta início de literal composto em multilinha
		if inFunction && strings.HasSuffix(trimmed, "{") &&
			!strings.HasPrefix(trimmed, "se ") &&
			!strings.HasPrefix(trimmed, "senao") &&
			!strings.HasPrefix(trimmed, "para ") &&
			!strings.HasPrefix(trimmed, "escolha ") &&
			!strings.HasPrefix(trimmed, "funcao ") {
			inComposite++
		}

		// Injeta trace apenas em declarações válidas fora de literais compostos
		if inFunction && inComposite == 0 && shouldTraceLine(trimmed) {
			escapedLine := escapeForTrace(trimmed)
			indent := getIndent(line)
			traceCall := fmt.Sprintf(`%sformatar.ImprimirLinha("\033[90m[Passo L.%03d]\033[0m \033[34m>>\033[0m \033[36m%s\033[0m")`, indent, lineNum, escapedLine)
			sb.WriteString(traceCall + "\n")
		}

		// Detecta fechamento de literal composto
		if inComposite > 0 {
			if strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, "},") {
				inComposite--
			}
		}

		// Conta abertura e fechamento de chaves
		openBraces := strings.Count(line, "{")
		closeBraces := strings.Count(line, "}")
		inBlock += openBraces - closeBraces
		if inBlock <= 0 {
			inFunction = false
			inBlock = 0
			inComposite = 0
		}

		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func shouldTraceLine(trimmed string) bool {
	if trimmed == "" || trimmed == "}" || trimmed == "{" || trimmed == "}," || trimmed == ");" {
		return false
	}
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
		return false
	}
	if strings.HasPrefix(trimmed, "caso ") || strings.HasPrefix(trimmed, "padrao:") || strings.HasPrefix(trimmed, "escolha ") {
		return false
	}
	// Não traça declarações de tipos ou pacotes
	if strings.HasPrefix(trimmed, "tipo ") || strings.HasPrefix(trimmed, "pacote ") || strings.HasPrefix(trimmed, "importar ") {
		return false
	}
	// Não traça itens de listas ou tuplas em literais
	if strings.HasSuffix(trimmed, ",") || strings.HasPrefix(trimmed, "{") {
		return false
	}
	return true
}

func escapeForTrace(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `'`)
	if len(s) > 60 {
		s = s[:57] + "..."
	}
	return s
}

func getIndent(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		if unicode.IsSpace(r) {
			buf.WriteRune(r)
		} else {
			break
		}
	}
	return buf.String()
}
