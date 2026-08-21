package main

import (
	"fmt"
	"os"
	"strings"

	"vamos/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintHelp()
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch command {
	case "run", "rodar", "executar":
		os.Exit(cmd.RunExecute(args))
	case "build", "compilar", "construir":
		os.Exit(cmd.RunBuild(args))
	case "test", "testar", "teste":
		os.Exit(cmd.RunTest(args))
	case "transpile", "converter", "transpilar":
		os.Exit(cmd.RunTranspile(args))
	case "go2vamos", "descompilar":
		os.Exit(cmd.RunGoToVamos(args))
	case "fmt", "formatar":
		os.Exit(cmd.RunFmt(args))
	case "lint", "checar", "verificar":
		os.Exit(cmd.RunLint(args))
	case "init", "iniciar":
		os.Exit(cmd.RunInit(args))
	case "repl", "interativo":
		os.Exit(cmd.RunREPL())
	case "playground", "web":
		os.Exit(cmd.RunWeb(args))
	case "version", "versao", "-v", "--version":
		cmd.PrintVersion()
		os.Exit(0)
	case "help", "ajuda", "-h", "--help":
		cmd.PrintHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Erro: comando desconhecido '%s'\n\n", command)
		cmd.PrintHelp()
		os.Exit(1)
	}
}
