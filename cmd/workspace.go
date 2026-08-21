package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// FindProjectRoot localiza a raiz do projeto procurando por go.mod a partir do diretório do arquivo.
func FindProjectRoot(srcPath string) string {
	abs, err := filepath.Abs(srcPath)
	if err != nil {
		cwd, _ := os.Getwd()
		return cwd
	}

	dir := abs
	if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
		dir = filepath.Dir(dir)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Se não encontrar go.mod, usa o diretório do arquivo fonte
	if fi, err := os.Stat(abs); err == nil && !fi.IsDir() {
		return filepath.Dir(abs)
	}
	return dir
}

// PrepareWorkspace transcreve e prepara todo o workspace para execução ou compilação com suporte completo a Go modules.
func PrepareWorkspace(rootDir string) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "vamos_workspace_*")
	if err != nil {
		return "", nil, fmt.Errorf("falha ao criar diretório temporário: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}

	hasGoMod := false

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			base := d.Name()
			if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" || base == "bin" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tempDir, relPath)

		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("erro ao ler '%s': %w", path, err)
			}

			goCode, err := transpiler.TranspileSource(string(content))
			if err != nil {
				return fmt.Errorf("erro ao transpilar '%s': %w", path, err)
			}

			goTargetPath := strings.TrimSuffix(targetPath, ".vamos") + ".go"
			if err := os.MkdirAll(filepath.Dir(goTargetPath), 0755); err != nil {
				return err
			}

			return os.WriteFile(goTargetPath, []byte(goCode), 0644)
		}

		if d.Name() == "go.mod" {
			hasGoMod = true
		}

		// Copia arquivos do projeto (configs, JSON, dados, go.mod, go.sum, etc.)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		return CopyFile(path, targetPath)
	})

	if err != nil {
		cleanup()
		return "", nil, err
	}

	if !hasGoMod {
		modContent := "module app\n\ngo 1.20\n"
		_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(modContent), 0644)
	}

	return tempDir, cleanup, nil
}

// CopyFile copia um arquivo de origem para o destino.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// RewriteErrors filtra mensagens do compilador Go para apontar para os arquivos .vamos originais.
func RewriteErrors(output, tempDir, rootDir string) string {
	res := strings.ReplaceAll(output, tempDir+string(filepath.Separator), "")
	res = strings.ReplaceAll(res, tempDir, "")
	lines := strings.Split(res, "\n")
	for i, line := range lines {
		if strings.Contains(line, ".go:") {
			lines[i] = strings.Replace(line, ".go:", ".vamos:", 1)
		}
	}
	return strings.Join(lines, "\n")
}

// RewriteFileErrors substitui o caminho do arquivo temporário pelo caminho do arquivo .vamos original.
func RewriteFileErrors(output, tempPath, originalPath string) string {
	return strings.ReplaceAll(output, tempPath, originalPath)
}
