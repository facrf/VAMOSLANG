package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"vamos/transpiler"
)

// RunExecute executa o comando 'vamos run <arquivo.vamos> [argumentos...]'
func RunExecute(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos especificado para execução.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos run <arquivo.vamos> [argumentos do programa...]\n")
		return 1
	}

	srcPath := args[0]
	progArgs := args[1:]

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

	tempFile, err := os.CreateTemp("", "vamos_run_*.go")
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

	// Prepara execução com 'go run'
	cmdArgs := append([]string{"run", tempGoPath}, progArgs...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// Captura stderr para filtrar e reescrever nomes de arquivos em mensagens de erro
	var stderrPipe io.ReadCloser
	stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		cmd.Stderr = os.Stderr
	}

	// Repassa sinais como Ctrl+C (SIGINT) e SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar o processo: %v\n", err)
		return 1
	}

	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()

	if stderrPipe != nil {
		var errBuf bytes.Buffer
		_, _ = io.Copy(&errBuf, stderrPipe)
		filtered := rewriteCompilerErrors(errBuf.String(), tempGoPath, srcPath)
		if len(filtered) > 0 {
			fmt.Fprint(os.Stderr, filtered)
		}
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
