package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGet instala dependências e bibliotecas externas para o projeto VAMOS-LANG.
func RunGet(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: nenhum pacote informado para instalação.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos get <caminho-do-pacote-go>\n")
		fmt.Fprintf(os.Stderr, "Exemplo: vamos get github.com/google/uuid\n")
		return 1
	}

	pkg := args[0]
	fmt.Printf("==> Baixando e instalando pacote externo: %s\n", pkg)

	cmdArgs := append([]string{"get"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter pacote '%s': %v\n", pkg, err)
		return 1
	}

	// Executa go mod tidy para organizar go.mod e go.sum
	_ = exec.Command("go", "mod", "tidy").Run()

	fmt.Printf("✓ Pacote '%s' instalado e registrado com sucesso!\n", pkg)
	return 0
}
