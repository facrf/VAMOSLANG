package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunTest executa o comando 'vamos test [flags...]' preparando o workspace e rodando go test
func RunTest(args []string) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter diretório atual: %v\n", err)
		return 1
	}

	tempDir, cleanup, err := PrepareWorkspace(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao preparar workspace para testes: %v\n", err)
		return 1
	}
	defer cleanup()

	testArgs := append([]string{"test"}, args...)
	if len(args) == 0 {
		testArgs = append(testArgs, "./...")
	}

	cmd := exec.Command("go", testArgs...)
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()

	if stdoutBuf.Len() > 0 {
		fmt.Print(RewriteErrors(stdoutBuf.String(), tempDir, cwd))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprint(os.Stderr, RewriteErrors(stderrBuf.String(), tempDir, cwd))
	}

	if execErr != nil {
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
