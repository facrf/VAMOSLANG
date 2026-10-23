package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunLint executa o comando 'vamos lint [--strict] [caminho...]'
func RunLint(args []string) int {
	targetPath := "."
	strict := false

	for _, arg := range args {
		if arg == "--strict" || arg == "-s" || arg == "--estrito" {
			strict = true
		} else if !strings.HasPrefix(arg, "-") {
			targetPath = arg
		}
	}

	totalIssues := 0
	hasErrors := false

	_ = filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && (strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if strict {
				diagnostics, _ := transpiler.ValidateSource(string(content))
				for _, diag := range diagnostics {
					prefix := "Aviso"
					if diag.Severity == transpiler.SeverityError {
						prefix = "Erro"
						hasErrors = true
					}
					fmt.Printf("%s:%d:%d: [%s] %s\n", path, diag.Line, diag.Column, prefix, diag.Message)
					totalIssues++
				}
			} else {
				issues := transpiler.LintSource(string(content))
				for _, iss := range issues {
					fmt.Printf("%s:%d: %s\n", path, iss.Line, iss.Message)
					totalIssues++
				}
			}
		}
		return nil
	})

	if totalIssues == 0 {
		fmt.Println("✓ Análise concluída: nenhum aviso ou problema encontrado.")
		return 0
	}
	fmt.Printf("\nTotal de %d ocorrência(s) encontrada(s).\n", totalIssues)
	if hasErrors {
		return 1
	}
	return 0
}
