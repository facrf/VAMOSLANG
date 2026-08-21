package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunTranspile executa o comando 'vamos transpile <arquivo.vamos> [-o <destino.go>] [--stdout]'
func RunTranspile(args []string) int {
	var srcPath string
	var outFile string
	toStdout := false

	// Análise flexível de argumentos e flags
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Erro: flag '%s' requer um caminho de destino.\n", arg)
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
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos especificado para transpilação.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos transpile <arquivo.vamos> [-o <destino.go>] [--stdout]\n")
		return 1
	}

	if !strings.HasSuffix(srcPath, ".vamos") {
		fmt.Fprintf(os.Stderr, "Aviso: o arquivo '%s' não possui a extensão .vamos\n", srcPath)
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
		return 1
	}

	goCode, err := transpiler.TranspileSource(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcPath, err)
		return 1
	}

	if toStdout {
		fmt.Print(goCode)
		return 0
	}

	target := outFile
	if target == "" {
		ext := filepath.Ext(srcPath)
		base := strings.TrimSuffix(srcPath, ext)
		target = base + ".go"
	}

	if err := os.WriteFile(target, []byte(goCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar arquivo gerado '%s': %v\n", target, err)
		return 1
	}

	fmt.Printf("✓ Arquivo transpilado com sucesso: %s -> %s\n", srcPath, target)
	return 0
}
