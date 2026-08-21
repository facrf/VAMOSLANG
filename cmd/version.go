package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const (
	Version  = "1.1.0"
	Codename = "Pindorama"
	Author   = "Equipe VAMOS-LANG"
	License  = "MIT"
)

const Banner = `
 __      __     __  __  ____   _____     _           _   _  _____ 
 \ \    / /\   |  \/  |/ __ \ / ____|   | |         / \ | \ |/ ____|
  \ \  / /  \  | \  / | |  | | (___     | |        / _ \|  \| | |  __ 
   \ \/ / /\ \ | |\/| | |  | |\___ \    | |       / /_\ \ . ' | | |_ |
    \  / ____ \| |  | | |__| |____) |   | |____  / _____ \ |\  | |__| |
     \/_/    \_\_|  |_|\____/|_____/    |______|/_/     \_\_| \_|\_____|
`

// PrintVersion exibe as informações detalhadas da versão do VAMOS-LANG.
func PrintVersion() {
	goVer := getGoVersion()

	fmt.Printf("┌──────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│               VAMOS-LANG - Versão %-22s │\n", Version)
	fmt.Printf("│   Linguagem Go 100%% em Português do Brasil (PT-BR)       │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│ Codinome:     %-42s │\n", Codename)
	fmt.Printf("│ Plataforma:   %-42s │\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	fmt.Printf("│ Compilador:   %-42s │\n", goVer)
	fmt.Printf("│ Licença:      %-42s │\n", License)
	fmt.Printf("│ Autor:        %-42s │\n", Author)
	fmt.Printf("└──────────────────────────────────────────────────────────┘\n")
}

// PrintHelp exibe o menu de ajuda completo da CLI.
func PrintHelp() {
	fmt.Print(Banner)
	fmt.Println("VAMOS-LANG: A linguagem de programação Go em Português do Brasil (PT-BR)")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  vamos <comando> [opções] [argumentos]")
	fmt.Println()
	fmt.Println("Comandos Principais (suporte a aliases em PT e EN):")
	fmt.Println("  run | rodar | executar <arquivo.vamos> [args...]   Transpila e executa diretamente")
	fmt.Println("  build | compilar | construir [flags...]            Transpila o projeto e compila binário")
	fmt.Println("  test | testar | teste [flags...]                   Transpila e executa os testes (go test)")
	fmt.Println("  transpile | converter | transpilar <origem> [-o]   Converte .vamos para Go (.go)")
	fmt.Println("  go2vamos | descompilar <arquivo.go> [-o]           Converte Go (.go) de volta para VAMOS")
	fmt.Println("  fmt | formatar [caminho...]                        Formata arquivos .vamos automaticamente")
	fmt.Println("  lint | checar | verificar [caminho...]             Verifica problemas de estilo e boas práticas")
	fmt.Println("  init | iniciar <nome-modulo>                       Cria um novo projeto com estrutura padrão")
	fmt.Println("  repl | interativo                                  Inicia o terminal interativo (REPL)")
	fmt.Println("  playground | web [porta]                           Inicia o Playground web interativo")
	fmt.Println("  version | versao                                   Exibe informações de versão")
	fmt.Println("  help | ajuda                                       Exibe esta tela de ajuda")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  vamos run teste.vamos")
	fmt.Println("  vamos build -o meu_app")
	fmt.Println("  vamos fmt .")
	fmt.Println("  vamos repl")
	fmt.Println("  vamos init meu_novo_app")
	fmt.Println("  vamos go2vamos main.go -o main.vamos")
	fmt.Println("  vamos web 8080")
	fmt.Println()
}

func getGoVersion() string {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return "Go não encontrado no PATH"
	}
	return strings.TrimSpace(string(out))
}
