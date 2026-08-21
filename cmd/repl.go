package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"vamos/transpiler"
)

// RunREPL inicia o shell interativo do VAMOS-LANG
func RunREPL() int {
	fmt.Println("==========================================================")
	fmt.Printf(" VAMOS-LANG REPL v%s (%s)\n", Version, Codename)
	fmt.Println(" Digite expressões ou declarações. Digite :ajuda para comandos.")
	fmt.Println("==========================================================")

	scanner := bufio.NewScanner(os.Stdin)
	var history []string

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
			fmt.Println("  :sair       Encerra o REPL")
			continue
		}

		if line == ":limpar" || line == ":clear" {
			history = nil
			fmt.Println("✓ Sessão limpa com sucesso.")
			continue
		}

		if line == ":codigo" || line == ":source" {
			fmt.Println("// Código acumulado da sessão:")
			for _, h := range history {
				fmt.Println(h)
			}
			continue
		}

		// Monta o programa de teste para execução
		history = append(history, line)

		var progBuilder strings.Builder
		progBuilder.WriteString("pacote principal\n\nimportar (\n\t\"formatar\"\n\t\"tempo\"\n\t\"matematica\"\n)\n\nfuncao principal() {\n")
		for _, h := range history {
			if !strings.HasPrefix(h, "variavel ") && !strings.HasPrefix(h, "constante ") && !strings.HasPrefix(h, "tipo ") && !strings.HasPrefix(h, "funcao ") && !strings.Contains(h, ":=") && !strings.HasPrefix(h, "formatar.") && !strings.HasPrefix(h, "se ") && !strings.HasPrefix(h, "para ") {
				progBuilder.WriteString(fmt.Sprintf("\tformatar.ImprimirLinha(%s)\n", h))
			} else {
				progBuilder.WriteString(fmt.Sprintf("\t%s\n", h))
			}
		}
		progBuilder.WriteString("}\n")

		goCode, err := transpiler.TranspileSource(progBuilder.String())
		if err != nil {
			fmt.Printf("Erro de sintaxe: %v\n", err)
			history = history[:len(history)-1] // desfaz a última linha
			continue
		}

		tempFile, err := os.CreateTemp("", "vamos_repl_*.go")
		if err != nil {
			fmt.Printf("Erro de I/O: %v\n", err)
			continue
		}
		tempGoPath := tempFile.Name()
		_ = os.WriteFile(tempGoPath, []byte(goCode), 0644)
		tempFile.Close()

		cmd := exec.Command("go", "run", tempGoPath)
		out, err := cmd.CombinedOutput()
		_ = os.Remove(tempGoPath)

		if err != nil {
			errMsg := string(out)
			lines := strings.Split(errMsg, "\n")
			if len(lines) > 0 {
				fmt.Printf("Erro: %s\n", lines[0])
			}
			history = history[:len(history)-1]
		} else {
			res := strings.TrimSpace(string(out))
			if res != "" {
				fmt.Println(res)
			}
		}
	}

	return 0
}
