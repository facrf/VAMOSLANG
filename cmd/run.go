package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"vamos/transpiler"
)

// RunExecute executa o comando 'vamos run <arquivo.vamos> [argumentos...]'
func RunExecute(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos informado para execução.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos run <arquivo.vamos> [argumentos...]\n")
		return 1
	}

	srcFile := args[0]
	progArgs := args[1:]

	absSrc, err := filepath.Abs(srcFile)
	if err != nil {
		absSrc = srcFile
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = filepath.Dir(absSrc)
	}

	content, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcFile, err)
		return 1
	}

	goCode, err := transpiler.TranspileSource(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcFile, err)
		return 1
	}

	projectRoot := FindProjectRoot(srcFile)
	hasOtherVamosFiles := false

	// Verifica se o projeto possui outros arquivos .vamos (ex: pacotes locais)
	if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
		_ = filepath.WalkDir(projectRoot, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if strings.HasSuffix(p, ".vamos") && p != absSrc {
				hasOtherVamosFiles = true
				return filepath.SkipAll
			}
			return nil
		})
	}

	// Caso A: Projeto com módulos locais ou múltiplos arquivos .vamos
	if hasOtherVamosFiles {
		tempDir, cleanup, err := PrepareWorkspace(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao preparar workspace: %v\n", err)
			return 1
		}
		defer cleanup()

		relPath, err := filepath.Rel(projectRoot, absSrc)
		if err != nil {
			relPath = filepath.Base(absSrc)
		}
		goEntryRelPath := strings.TrimSuffix(relPath, ".vamos") + ".go"

		runArgs := append([]string{"run", goEntryRelPath}, progArgs...)
		cmd := exec.Command("go", runArgs...)
		cmd.Dir = tempDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		var stderrPipe io.ReadCloser
		stderrPipe, err = cmd.StderrPipe()
		if err != nil {
			cmd.Stderr = os.Stderr
		}

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
			filtered := RewriteErrors(errBuf.String(), tempDir, projectRoot)
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

	// Caso B: Execução direta no diretório para arquivos simples / scripts únicos
	srcDir := filepath.Dir(absSrc)
	tempFile, err := os.CreateTemp(srcDir, "vamos_run_*.go")
	if err != nil {
		tempFile, err = os.CreateTemp("", "vamos_run_*.go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo temporário: %v\n", err)
			return 1
		}
	}
	tempGoPath := tempFile.Name()
	defer os.Remove(tempGoPath)

	if _, err := tempFile.WriteString(goCode); err != nil {
		tempFile.Close()
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo temporário: %v\n", err)
		return 1
	}
	tempFile.Close()

	runArgs := append([]string{"run", tempGoPath}, progArgs...)
	cmd := exec.Command("go", runArgs...)
	cmd.Dir = cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderrPipe io.ReadCloser
	stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		cmd.Stderr = os.Stderr
	}

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
		filtered := RewriteFileErrors(errBuf.String(), tempGoPath, srcFile)
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
