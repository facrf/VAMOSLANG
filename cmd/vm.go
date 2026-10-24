package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// findVamosc localiza o executável do compilador nativo vamosc
func findVamosc() (string, error) {
	// 1. Procura no PATH do sistema
	if p, err := exec.LookPath("vamosc"); err == nil {
		return p, nil
	}

	// 2. Procura a partir do diretório atual de trabalho
	cwd, _ := os.Getwd()
	rootCwd := FindProjectRoot(cwd)
	candidate := filepath.Join(rootCwd, "compilador", "vamosc")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	// 3. Procura a partir da localização do binário 'vamos'
	if exe, err := os.Executable(); err == nil {
		rootExe := FindProjectRoot(filepath.Dir(exe))
		candidateExe := filepath.Join(rootExe, "compilador", "vamosc")
		if _, err := os.Stat(candidateExe); err == nil {
			return candidateExe, nil
		}
	}

	// 4. Se existir código fonte compilador/main.vamos, compila sob demanda
	mainVamos := filepath.Join(rootCwd, "compilador", "main.vamos")
	if _, err := os.Stat(mainVamos); err == nil {
		buildArgs := []string{"build", "-o", candidate, mainVamos}
		if code := RunBuild(buildArgs); code == 0 {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("compilador auto-hospedado 'vamosc' não encontrado")
}

// RunVM executa código .vamos ou .vamosb diretamente na Máquina Virtual VVM
func RunVM(args []string) int {
	if len(args) == 0 {
		fmt.Println("Erro: informe o arquivo .vamos ou .vamosb para executar na VM.")
		fmt.Println("Uso: vamos vm <arquivo.vamos|.vamosb>")
		return 1
	}

	vamoscPath, err := findVamosc()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		return 1
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

	vamoscPath, err := findVamosc()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		return 1
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

// RunWasm compila código VAMOS-LANG diretamente em módulo WebAssembly 1.0 (.wasm)
func RunWasm(args []string) int {
	if len(args) == 0 {
		fmt.Println("Erro: informe o arquivo .vamos para compilar em WebAssembly.")
		fmt.Println("Uso: vamos wasm <arquivo.vamos> [-o saida.wasm]")
		return 1
	}

	vamoscPath, err := findVamosc()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		return 1
	}

	execArgs := append([]string{"wasm"}, args...)
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
