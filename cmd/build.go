package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vamos/transpiler"
)

// RunBuild executa o comando 'vamos build [flags...] [arquivo.vamos]'
func RunBuild(args []string) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter diretório atual: %v\n", err)
		return 1
	}

	var forwardArgs []string
	var targetFiles []string
	var outFile string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Erro: flag '%s' requer o nome do binário de destino.\n", arg)
				return 1
			}
		} else if strings.HasPrefix(arg, "-o=") {
			outFile = strings.TrimPrefix(arg, "-o=")
		} else if strings.HasPrefix(arg, "--output=") {
			outFile = strings.TrimPrefix(arg, "--output=")
		} else if strings.HasSuffix(arg, ".vamos") {
			targetFiles = append(targetFiles, arg)
		} else {
			forwardArgs = append(forwardArgs, arg)
		}
	}

	// Caso A: Arquivo específico .vamos
	if len(targetFiles) == 1 {
		srcPath := targetFiles[0]
		absSrcPath, err := filepath.Abs(srcPath)
		if err != nil {
			absSrcPath = srcPath
		}

		if outFile == "" {
			base := filepath.Base(srcPath)
			ext := filepath.Ext(base)
			outFile = strings.TrimSuffix(base, ext)
		}

		absOutFile, err := filepath.Abs(outFile)
		if err != nil {
			absOutFile = outFile
		}

		projectRoot := FindProjectRoot(srcPath)
		hasOtherVamosFiles := false

		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			_ = filepath.WalkDir(projectRoot, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				if strings.HasSuffix(p, ".vamos") && p != absSrcPath {
					hasOtherVamosFiles = true
					return filepath.SkipAll
				}
				return nil
			})
		}

		// Se tiver módulos locais, compila no workspace
		if hasOtherVamosFiles {
			tempDir, cleanup, err := PrepareWorkspace(projectRoot)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao preparar workspace: %v\n", err)
				return 1
			}
			defer cleanup()

			relPath, err := filepath.Rel(projectRoot, absSrcPath)
			if err != nil {
				relPath = filepath.Base(absSrcPath)
			}
			goEntryRelPath := strings.TrimSuffix(relPath, ".vamos") + ".go"

			buildArgs := []string{"build", "-o", absOutFile}
			buildArgs = append(buildArgs, forwardArgs...)
			buildArgs = append(buildArgs, goEntryRelPath)

			cmd := exec.Command("go", buildArgs...)
			cmd.Dir = tempDir

			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			execErr := cmd.Run()

			if stdoutBuf.Len() > 0 {
				fmt.Print(RewriteErrors(stdoutBuf.String(), tempDir, projectRoot))
			}
			if stderrBuf.Len() > 0 {
				fmt.Fprint(os.Stderr, RewriteErrors(stderrBuf.String(), tempDir, projectRoot))
			}

			if execErr != nil {
				fmt.Fprintf(os.Stderr, "Falha na compilação do binário.\n")
				return 1
			}

			printBinarySuccess(outFile, absOutFile)
			return 0
		}

		// Script único simples
		content, err := os.ReadFile(srcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
			return 1
		}

		goCode, err := transpiler.TranspileSource(string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcPath, err)
			return 1
		}

		srcDir := filepath.Dir(absSrcPath)
		tempFile, err := os.CreateTemp(srcDir, "vamos_build_*.go")
		if err != nil {
			tempFile, err = os.CreateTemp("", "vamos_build_*.go")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao criar arquivo temporário: %v\n", err)
				return 1
			}
		}
		tempGoPath := tempFile.Name()
		defer os.Remove(tempGoPath)

		if _, err := tempFile.WriteString(goCode); err != nil {
			tempFile.Close()
			fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo temporário: %v\n", err)
			return 1
		}
		tempFile.Close()

		buildArgs := []string{"build", "-o", absOutFile}
		buildArgs = append(buildArgs, forwardArgs...)
		buildArgs = append(buildArgs, tempGoPath)

		cmd := exec.Command("go", buildArgs...)
		cmd.Dir = srcDir

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		execErr := cmd.Run()

		if stdoutBuf.Len() > 0 {
			fmt.Print(RewriteFileErrors(stdoutBuf.String(), tempGoPath, srcPath))
		}
		if stderrBuf.Len() > 0 {
			fmt.Fprint(os.Stderr, RewriteFileErrors(stderrBuf.String(), tempGoPath, srcPath))
		}

		if execErr != nil {
			fmt.Fprintf(os.Stderr, "Falha na compilação do binário.\n")
			return 1
		}

		printBinarySuccess(outFile, absOutFile)
		return 0
	}

	// Caso B: Build de todo o projeto / workspace
	if outFile == "" {
		outFile = filepath.Base(cwd)
	}

	absOutFile, err := filepath.Abs(outFile)
	if err != nil {
		absOutFile = outFile
	}

	tempDir, cleanup, err := PrepareWorkspace(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao preparar workspace: %v\n", err)
		return 1
	}
	defer cleanup()

	// Identifica o ponto de entrada principal do pacote main
	entryTarget := "."
	if !hasMainInDir(tempDir) {
		if hasMainInDir(filepath.Join(tempDir, "cmd")) {
			entryTarget = "./cmd"
		} else {
			// Procura em subdiretórios de cmd/
			_ = filepath.WalkDir(filepath.Join(tempDir, "cmd"), func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() && hasMainInDir(p) {
					rel, _ := filepath.Rel(tempDir, p)
					entryTarget = "./" + rel
					return filepath.SkipAll
				}
				return nil
			})
		}
	}

	buildArgs := []string{"build", "-o", absOutFile}
	buildArgs = append(buildArgs, forwardArgs...)
	buildArgs = append(buildArgs, entryTarget)

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()

	if stdoutBuf.Len() > 0 {
		fmt.Print(RewriteErrors(stdoutBuf.String(), tempDir, cwd))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprint(os.Stderr, RewriteErrors(stderrBuf.String(), tempDir, cwd))
	}

	if execErr != nil {
		fmt.Fprintf(os.Stderr, "Falha na compilação.\n")
		return 1
	}

	printBinarySuccess(outFile, absOutFile)
	return 0
}

func hasMainInDir(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			content, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err == nil && strings.Contains(string(content), "package main") {
				return true
			}
		}
	}
	return false
}

func printBinarySuccess(outFile, absOutFile string) {
	fi, err := os.Stat(absOutFile)
	sizeStr := ""
	if err == nil {
		sizeKB := float64(fi.Size()) / 1024.0
		if sizeKB > 1024 {
			sizeStr = fmt.Sprintf(" (%.2f MB)", sizeKB/1024.0)
		} else {
			sizeStr = fmt.Sprintf(" (%.1f KB)", sizeKB)
		}
	}
	fmt.Printf("✓ Binário gerado com sucesso: %s%s\n", outFile, sizeStr)
}
