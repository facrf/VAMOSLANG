package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunBuild executa o comando 'vamos build <arquivo.vamos> [-o <binario>]'
func RunBuild(args []string) int {
	var srcPath string
	var outFile string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Erro: flag '%s' requer o nome do binário de destino.\n", arg)
				return 1
			}
		} else if strings.HasPrefix(arg, "-o=") {
			outFile = strings.TrimPrefix(arg, "-o=")
		} else if strings.HasPrefix(arg, "--output=") {
			outFile = strings.TrimPrefix(arg, "--output=")
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
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos especificado para compilação.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos build <arquivo.vamos> [-o <binario>]\n")
		return 1
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

	targetBinary := outFile
	if targetBinary == "" {
		base := filepath.Base(srcPath)
		ext := filepath.Ext(base)
		targetBinary = strings.TrimSuffix(base, ext)
	}

	absTargetBinary, err := filepath.Abs(targetBinary)
	if err != nil {
		absTargetBinary = targetBinary
	}

	// Cria arquivo temporário Go no diretório temporário do sistema
	tempFile, err := os.CreateTemp("", "vamos_build_*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo temporário: %v\n", err)
		return 1
	}
	tempGoPath := tempFile.Name()
	defer os.Remove(tempGoPath)

	if _, err := tempFile.WriteString(goCode); err != nil {
		tempFile.Close()
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo temporário: %v\n", err)
		return 1
	}
	tempFile.Close()

	// Executa go build -o <binario> <tempGoPath>
	cmd := exec.Command("go", "build", "-o", absTargetBinary, tempGoPath)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()

	if stdoutBuf.Len() > 0 {
		fmt.Print(rewriteCompilerErrors(stdoutBuf.String(), tempGoPath, srcPath))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprint(os.Stderr, rewriteCompilerErrors(stderrBuf.String(), tempGoPath, srcPath))
	}

	if execErr != nil {
		fmt.Fprintf(os.Stderr, "Falha na compilação do binário.\n")
		return 1
	}

	// Estatísticas do binário gerado
	fi, err := os.Stat(absTargetBinary)
	sizeStr := ""
	if err == nil {
		sizeKB := float64(fi.Size()) / 1024.0
		if sizeKB > 1024 {
			sizeStr = fmt.Sprintf(" (%.2f MB)", sizeKB/1024.0)
		} else {
			sizeStr = fmt.Sprintf(" (%.1f KB)", sizeKB)
		}
	}

	fmt.Printf("✓ Binário gerado com sucesso: %s%s\n", targetBinary, sizeStr)
	return 0
}

func rewriteCompilerErrors(output, tempPath, originalPath string) string {
	tempBase := filepath.Base(tempPath)
	origBase := filepath.Base(originalPath)

	res := strings.ReplaceAll(output, tempPath, originalPath)
	res = strings.ReplaceAll(res, tempBase, origBase)
	return res
}
