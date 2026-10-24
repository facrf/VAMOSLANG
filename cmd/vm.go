package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunVM executa código .vamos ou .vamosb diretamente na Máquina Virtual VVM
func RunVM(args []string) int {
	if len(args) == 0 {
		fmt.Println("Erro: informe o arquivo .vamos ou .vamosb para executar na VM.")
		fmt.Println("Uso: vamos vm <arquivo.vamos|.vamosb>")
		return 1
	}

	root := FindProjectRoot(args[0])
	vamoscPath := filepath.Join(root, "compilador", "vamosc")

	// Se o compilador vamosc ainda não foi construído, constrói automaticamente
	if _, err := os.Stat(vamoscPath); err != nil {
		buildArgs := []string{"build", "-o", vamoscPath, filepath.Join(root, "compilador", "main.vamos")}
		if code := RunBuild(buildArgs); code != 0 {
			return code
		}
	}

	execArgs := append([]string{"vm"}, args...)
	cmd := exec.Command(vamoscPath, execArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

// RunBytecode gera ou desassembla bytecode VVM (.vamosb)
func RunBytecode(args []string) int {
	if len(args) == 0 {
		fmt.Println("Erro: informe o arquivo .vamos ou .vamosb para processar bytecode.")
		fmt.Println("Uso: vamos bytecode <arquivo.vamos> [-o saida.vamosb]")
		return 1
	}

	root := FindProjectRoot(args[0])
	vamoscPath := filepath.Join(root, "compilador", "vamosc")

	if _, err := os.Stat(vamoscPath); err != nil {
		buildArgs := []string{"build", "-o", vamoscPath, filepath.Join(root, "compilador", "main.vamos")}
		if code := RunBuild(buildArgs); code != 0 {
			return code
		}
	}

	execArgs := append([]string{"bytecode"}, args...)
	cmd := exec.Command(vamoscPath, execArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
