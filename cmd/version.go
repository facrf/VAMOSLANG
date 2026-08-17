package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const (
	Version   = "1.0.0"
	Codename  = "Pindorama"
	Author    = "Equipe VAMOS-LANG"
	License   = "MIT"
)

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

func getGoVersion() string {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return "Go não encontrado no PATH"
	}
	return strings.TrimSpace(string(out))
}
