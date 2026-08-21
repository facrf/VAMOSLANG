package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunFmt executa o comando 'vamos fmt [caminho...]'
func RunFmt(args []string) int {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	count := 0
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			formatted, err := transpiler.FormatSource(string(content))
			if err != nil {
				return err
			}
			if string(content) != formatted {
				if err := os.WriteFile(path, []byte(formatted), 0644); err != nil {
					return err
				}
				fmt.Printf("✓ Formatado: %s\n", path)
				count++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao formatar arquivos: %v\n", err)
		return 1
	}

	if count == 0 {
		fmt.Println("✓ Todos os arquivos .vamos já estão devidamente formatados.")
	} else {
		fmt.Printf("✓ Total de %d arquivo(s) formatado(s) com sucesso.\n", count)
	}
	return 0
}
