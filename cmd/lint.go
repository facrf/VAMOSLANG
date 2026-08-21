package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunLint executa o comando 'vamos lint [caminho...]'
func RunLint(args []string) int {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	totalIssues := 0
	_ = filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && (strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err == nil {
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
		fmt.Println("✓ Análise concluída: nenhum aviso ou problema de estilo encontrado.")
		return 0
	}
	fmt.Printf("\nTotal de %d aviso(s) encontrado(s).\n", totalIssues)
	return 0
}
