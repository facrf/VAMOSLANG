package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RunWeb inicia o servidor do Playground web
func RunWeb(args []string) int {
	port := "8080"
	if len(args) > 0 {
		port = args[0]
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	playgroundDir := "playground"
	if _, err := os.Stat(playgroundDir); os.IsNotExist(err) {
		exePath, err := os.Executable()
		if err == nil {
			candidate := filepath.Join(filepath.Dir(exePath), "playground")
			if _, err := os.Stat(candidate); err == nil {
				playgroundDir = candidate
			}
		}
	}

	fs := http.FileServer(http.Dir(playgroundDir))
	http.Handle("/", fs)

	fmt.Printf("🚀 VAMOS-LANG Playground disponível em http://localhost%s\n", port)
	fmt.Println("Pressione Ctrl+C para encerrar o servidor.")
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar servidor web: %v\n", err)
		return 1
	}
	return 0
}
