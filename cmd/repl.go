package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunREPL inicia o shell interativo do VAMOS-LANG
func RunREPL() int {
	fmt.Println("==========================================================")
	fmt.Printf(" VAMOS-LANG REPL v%s (%s)\n", Version, Codename)
	fmt.Println(" Digite expressões, variáveis, funções ou importe pacotes.")
	fmt.Println(" Digite :ajuda para comandos especiais ou :sair para sair.")
	fmt.Println("==========================================================")

	scanner := bufio.NewScanner(os.Stdin)
	var history []string
	customImports := []string{"formatar", "tempo", "matematica", "cordas"}
	var topLevelDecls []string
	var bodyStatements []string

	projectRoot := FindProjectRoot(".")

	for {
		fmt.Print("vamos> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if line == ":sair" || line == ":exit" || line == ":q" {
			fmt.Println("Até logo!")
			break
		}

		if line == ":ajuda" || line == ":help" {
			fmt.Println("Comandos do REPL:")
			fmt.Println("  :ajuda      Exibe esta ajuda")
			fmt.Println("  :codigo     Exibe o histórico de código acumulado na sessão")
			fmt.Println("  :limpar     Limpa a sessão atual do REPL")
			fmt.Println("  :pacotes    Exibe os pacotes importados na sessão")
			fmt.Println("  :sair       Encerra o REPL")
			fmt.Println("\nDicas de sintaxe:")
			fmt.Println("  - Expressões: 10 + 20 * 2 ou cordas.ParaMaiusculas(\"teste\")")
			fmt.Println("  - Variáveis:  x := 42 ou variavel nome texto = \"Ana\"")
			fmt.Println("  - Importar:   importar \"so\" ou importar \"texto_util\"")
			fmt.Println("  - Funções:    funcao dobro(n inteiro) inteiro { retornar n * 2 }")
			continue
		}

		if line == ":limpar" || line == ":clear" {
			history = nil
			topLevelDecls = nil
			bodyStatements = nil
			customImports = []string{"formatar", "tempo", "matematica", "cordas"}
			fmt.Println("✓ Sessão limpa com sucesso.")
			continue
		}

		if line == ":pacotes" || line == ":imports" {
			fmt.Println("// Pacotes importados nesta sessão:")
			for _, imp := range customImports {
				fmt.Printf("  importar %q\n", imp)
			}
			continue
		}

		if line == ":codigo" || line == ":source" {
			fmt.Println("// Código acumulado da sessão:")
			for _, h := range history {
				fmt.Println(h)
			}
			continue
		}

		// Identifica importação dinâmica
		if strings.HasPrefix(line, "importar ") || strings.HasPrefix(line, "import ") {
			rawPath := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "importar"), "import"))
			rawPath = strings.Trim(rawPath, `"'(); `)
			if rawPath != "" {
				jaImportado := false
				for _, imp := range customImports {
					if imp == rawPath {
						jaImportado = true
						break
					}
				}
				if !jaImportado {
					customImports = append(customImports, rawPath)
					history = append(history, fmt.Sprintf("importar %q", rawPath))
					fmt.Printf("✓ Pacote %q importado com sucesso.\n", rawPath)
				} else {
					fmt.Printf("Pacote %q já está importado.\n", rawPath)
				}
				continue
			}
		}

		// Identifica declarações top-level (tipos e funções auxiliares)
		isTopLevel := false
		if strings.HasPrefix(line, "tipo ") || (strings.HasPrefix(line, "funcao ") && !strings.HasPrefix(line, "funcao principal")) {
			isTopLevel = true
		}

		var formattedStmt string
		if isTopLevel {
			formattedStmt = line
		} else {
			if !strings.HasPrefix(line, "variavel ") &&
				!strings.HasPrefix(line, "constante ") &&
				!strings.Contains(line, ":=") &&
				!strings.HasPrefix(line, "formatar.") &&
				!strings.HasPrefix(line, "se ") &&
				!strings.HasPrefix(line, "para ") &&
				!strings.HasPrefix(line, "escolha ") &&
				!strings.HasPrefix(line, "disparar ") &&
				!strings.HasPrefix(line, "adiar ") {
				formattedStmt = fmt.Sprintf("formatar.ImprimirLinha(%s)", line)
			} else {
				formattedStmt = line
			}
		}

		// Monta o programa de teste para execução
		var progBuilder strings.Builder
		progBuilder.WriteString("pacote principal\n\nimportar (\n")
		for _, imp := range customImports {
			progBuilder.WriteString(fmt.Sprintf("\t%q\n", imp))
		}
		progBuilder.WriteString(")\n\n")

		for _, tl := range topLevelDecls {
			progBuilder.WriteString(tl + "\n\n")
		}
		if isTopLevel {
			progBuilder.WriteString(formattedStmt + "\n\n")
		}

		progBuilder.WriteString("funcao principal() {\n")
		for _, stmt := range bodyStatements {
			progBuilder.WriteString("\t" + stmt + "\n")
		}
		if !isTopLevel {
			progBuilder.WriteString("\t" + formattedStmt + "\n")
		}
		progBuilder.WriteString("}\n")

		goCode, err := transpiler.TranspileSource(progBuilder.String())
		if err != nil {
			fmt.Printf("Erro de sintaxe: %v\n", err)
			continue
		}

		tempDir, err := os.MkdirTemp("", "vamos_repl_*")
		if err != nil {
			fmt.Printf("Erro de I/O: %v\n", err)
			continue
		}

		// Configura módulo temporário apontando para o workspace se houver
		if modData, err := os.ReadFile(filepath.Join(projectRoot, "go.mod")); err == nil {
			_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), modData, 0644)
			if sumData, err := os.ReadFile(filepath.Join(projectRoot, "go.sum")); err == nil {
				_ = os.WriteFile(filepath.Join(tempDir, "go.sum"), sumData, 0644)
			}
		} else {
			_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module repl_app\n\ngo 1.22\n"), 0644)
		}

		tempGoPath := filepath.Join(tempDir, "main.go")
		_ = os.WriteFile(tempGoPath, []byte(goCode), 0644)

		cmd := exec.Command("go", "run", "main.go")
		cmd.Dir = tempDir
		out, err := cmd.CombinedOutput()
		_ = os.RemoveAll(tempDir)

		if err != nil {
			translated := TranslateCompilerMessages(string(out))
			lines := strings.Split(translated, "\n")
			for _, l := range lines {
				if strings.TrimSpace(l) != "" {
					fmt.Printf("Erro: %s\n", l)
					break
				}
			}
		} else {
			if isTopLevel {
				topLevelDecls = append(topLevelDecls, formattedStmt)
				history = append(history, line)
				fmt.Println("✓ Declaração registrada.")
			} else {
				bodyStatements = append(bodyStatements, formattedStmt)
				history = append(history, line)
				res := strings.TrimSpace(string(out))
				if res != "" {
					fmt.Println(res)
				}
			}
		}
	}

	return 0
}
