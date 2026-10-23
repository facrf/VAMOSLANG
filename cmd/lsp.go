package cmd

import (
	"os"

	"vamos/lsp"
)

// RunLSP inicia o Language Server Protocol (LSP) para integração com editores (VS Code, Cursor, Neovim).
func RunLSP() int {
	server := lsp.NewServer(os.Stdin, os.Stdout)
	if err := server.Run(); err != nil {
		return 1
	}
	return 0
}
