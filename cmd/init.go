package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunInit executa o comando 'vamos init <nome-ou-caminho-modulo>'
func RunInit(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: informe o nome do módulo para o novo projeto.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos init <nome-do-modulo>\n")
		return 1
	}

	targetInput := args[0]
	projectDir := targetInput
	modName := filepath.Base(targetInput)

	if targetInput == "." {
		cwd, err := os.Getwd()
		if err == nil {
			modName = filepath.Base(cwd)
		} else {
			modName = "app"
		}
		projectDir = "."
	}

	// Sanitiza o nome do módulo para ser válido no Go
	modName = strings.TrimPrefix(modName, ".")
	modName = strings.TrimPrefix(modName, "/")
	if modName == "" {
		modName = "app"
	}

	fmt.Printf("==> Criando novo projeto VAMOS-LANG: %s\n", modName)

	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "calculos"), 0755)

	// go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.20\n", modName)
	_ = os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644)

	// cmd/principal.vamos
	mainContent := fmt.Sprintf(`// Ponto de entrada da aplicação
pacote principal

importar (
	"formatar"
	"tempo"
	"%s/pacotes/calculos"
)

funcao principal() {
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirLinha("   Novo Projeto VAMOS-LANG Iniciado!    ")
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirFormatado("Executado em: %%s\n", tempo.Agora().Formato("02/01/2006 15:04:05"))

	resultado := calculos.Somar(15, 27)
	formatar.ImprimirFormatado("Teste de Módulo Local (15 + 27): %%d\n", resultado)
}
`, modName)
	_ = os.WriteFile(filepath.Join(projectDir, "cmd", "principal.vamos"), []byte(mainContent), 0644)

	// pacotes/calculos/calculos.vamos
	calcContent := `pacote calculos

funcao Somar(a inteiro, b inteiro) inteiro {
	retornar a + b
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "calculos", "calculos.vamos"), []byte(calcContent), 0644)

	// .gitignore
	gitIgnore := `# Binários
bin/
*.exe
*.out
app
`
	_ = os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitIgnore), 0644)

	fmt.Println("✓ Estrutura de pastas gerada com sucesso!")
	fmt.Printf("\nPara começar a desenvolver:\n")
	if projectDir != "." {
		fmt.Printf("  cd %s\n", projectDir)
	}
	fmt.Printf("  vamos run cmd/principal.vamos\n\n")
	return 0
}
