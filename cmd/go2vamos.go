package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunGoToVamos executa o comando 'vamos go2vamos <arquivo.go> [-o <destino.vamos>] [--stdout]'
func RunGoToVamos(args []string) int {
	var srcPath string
	var outFile string
	toStdout := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Erro: flag '%s' requer o nome do arquivo de destino.\n", arg)
				return 1
			}
		} else if strings.HasPrefix(arg, "-o=") {
			outFile = strings.TrimPrefix(arg, "-o=")
		} else if strings.HasPrefix(arg, "--output=") {
			outFile = strings.TrimPrefix(arg, "--output=")
		} else if arg == "--stdout" || arg == "-stdout" {
			toStdout = true
		} else if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Erro: flag desconhecida '%s'\n", arg)
			return 1
		} else if srcPath == "" {
			srcPath = arg
		} else {
			fmt.Fprintf(os.Stderr, "Erro: argumento inesperado '%s'\n", arg)
			return 1
		}
	}

	if srcPath == "" {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .go especificado para conversão reversa.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos go2vamos <arquivo.go> [-o <destino.vamos>] [--stdout]\n")
		return 1
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
		return 1
	}

	vamosCode, err := transpiler.TranspileGoToVamos(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter Go para VAMOS: %v\n", err)
		return 1
	}

	if toStdout {
		fmt.Print(vamosCode)
		return 0
	}

	target := outFile
	if target == "" {
		ext := filepath.Ext(srcPath)
		base := strings.TrimSuffix(srcPath, ext)
		target = base + ".vamos"
	}

	if err := os.WriteFile(target, []byte(vamosCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar arquivo .vamos '%s': %v\n", target, err)
		return 1
	}

	fmt.Printf("✓ Código Go convertido para VAMOS-LANG com sucesso: %s -> %s\n", srcPath, target)
	return 0
}
