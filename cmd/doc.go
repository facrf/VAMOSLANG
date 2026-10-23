package cmd

import (
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DocSymbol representa um elemento documentado do código VAMOS-LANG.
type DocSymbol struct {
	Name        string
	Kind        string // "pacote", "funcao", "metodo", "estrutura", "interface", "constante"
	Signature   string
	Description string
	File        string
	Line        int
}

// RunDoc executa o gerador e visualizador de documentação VAMOS-LANG.
func RunDoc(args []string) int {
	targetPath := "."
	htmlOutput := ""
	querySymbol := ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--html" {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				htmlOutput = args[i+1]
				i++
			} else {
				htmlOutput = "docs.html"
			}
		} else if strings.HasPrefix(arg, "--html=") {
			htmlOutput = strings.TrimPrefix(arg, "--html=")
		} else if !strings.HasPrefix(arg, "-") {
			if _, err := os.Stat(arg); err == nil {
				targetPath = arg
			} else {
				querySymbol = arg
			}
		}
	}

	symbols, err := extractDocs(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao extrair documentação: %v\n", err)
		return 1
	}

	if len(symbols) == 0 {
		fmt.Println("Nenhum símbolo ou comentário de documentação encontrado.")
		return 0
	}

	// 1. Busca por símbolo específico
	if querySymbol != "" {
		found := false
		lowerQuery := strings.ToLower(querySymbol)
		for _, s := range symbols {
			if strings.ToLower(s.Name) == lowerQuery || strings.Contains(strings.ToLower(s.Name), lowerQuery) {
				found = true
				fmt.Printf("\n=== [%s] %s ===\n", strings.ToUpper(s.Kind), s.Name)
				fmt.Printf("Arquivo: %s:%d\n", s.File, s.Line)
				if s.Signature != "" {
					fmt.Printf("Assinatura: %s\n", s.Signature)
				}
				if s.Description != "" {
					fmt.Printf("\nDescrição:\n%s\n", s.Description)
				}
				fmt.Println(strings.Repeat("-", 50))
			}
		}
		if !found {
			fmt.Printf("Símbolo '%s' não encontrado na documentação.\n", querySymbol)
		}
		return 0
	}

	// 2. Exportação para HTML
	if htmlOutput != "" {
		err := generateHTML(symbols, htmlOutput, targetPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao gerar arquivo HTML: %v\n", err)
			return 1
		}
		fmt.Printf("✓ Documentação HTML gerada com sucesso em: %s\n", htmlOutput)
		return 0
	}

	// 3. Exibição Geral no Terminal
	fmt.Println("====================================================================")
	fmt.Println("                📖 DOCUMENTAÇÃO VAMOS-LANG                          ")
	fmt.Println("====================================================================")

	// Agrupa por tipo
	packages := filterKind(symbols, "pacote")
	types := filterKind(symbols, "estrutura", "interface")
	functions := filterKind(symbols, "funcao")
	methods := filterKind(symbols, "metodo")

	if len(packages) > 0 {
		fmt.Println("\n📦 Pacotes:")
		for _, p := range packages {
			fmt.Printf("  • %-20s %s (%s)\n", p.Name, p.Description, p.File)
		}
	}

	if len(types) > 0 {
		fmt.Println("\n🏛️  Tipos e Estruturas:")
		for _, t := range types {
			desc := t.Description
			if desc == "" {
				desc = "Sem descrição."
			}
			fmt.Printf("  • %-20s [%s] - %s\n", t.Name, t.Kind, desc)
			if t.Signature != "" {
				fmt.Printf("    %s\n", t.Signature)
			}
		}
	}

	if len(functions) > 0 {
		fmt.Println("\n⚡ Funções:")
		for _, f := range functions {
			desc := f.Description
			if desc == "" {
				desc = "Sem descrição."
			}
			fmt.Printf("  • %-25s %s\n", f.Name, desc)
			if f.Signature != "" {
				fmt.Printf("    %s\n", f.Signature)
			}
		}
	}

	if len(methods) > 0 {
		fmt.Println("\n🔧 Métodos:")
		for _, m := range methods {
			fmt.Printf("  • %-25s %s\n", m.Name, m.Description)
			if m.Signature != "" {
				fmt.Printf("    %s\n", m.Signature)
			}
		}
	}

	fmt.Printf("\nTotal de %d símbolo(s) documentado(s).\n", len(symbols))
	fmt.Println("Dica: Use 'vamos doc <nome-simbolo>' para detalhes ou 'vamos doc --html docs.html' para exportar.")
	return 0
}

func filterKind(symbols []DocSymbol, kinds ...string) []DocSymbol {
	var res []DocSymbol
	for _, s := range symbols {
		for _, k := range kinds {
			if s.Kind == k {
				res = append(res, s)
				break
			}
		}
	}
	return res
}

func extractDocs(root string) ([]DocSymbol, error) {
	var symbols []DocSymbol

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && d.Name() != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".vamos") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		fileSymbols := parseFileDocs(string(data), path)
		symbols = append(symbols, fileSymbols...)
		return nil
	})

	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i].Name < symbols[j].Name
	})

	return symbols, err
}

func parseFileDocs(content, filePath string) []DocSymbol {
	var symbols []DocSymbol
	lines := strings.Split(content, "\n")
	var pendingComments []string

	for idx, line := range lines {
		lineNum := idx + 1
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			commentText := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
			if !strings.HasPrefix(commentText, "==") && !strings.HasPrefix(commentText, "--") {
				pendingComments = append(pendingComments, commentText)
			}
			continue
		}

		if trimmed == "" {
			if len(pendingComments) > 0 && len(pendingComments) > 4 {
				pendingComments = nil
			}
			continue
		}

		docText := strings.Join(pendingComments, " ")
		pendingComments = nil

		// 1. Pacote
		if strings.HasPrefix(trimmed, "pacote ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				symbols = append(symbols, DocSymbol{
					Name:        parts[1],
					Kind:        "pacote",
					Signature:   trimmed,
					Description: docText,
					File:        filePath,
					Line:        lineNum,
				})
			}
			continue
		}

		// 2. Tipos (estrutura ou interface)
		if strings.HasPrefix(trimmed, "tipo ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 3 {
				kind := parts[2]
				if kind == "estrutura" || kind == "interface" {
					symbols = append(symbols, DocSymbol{
						Name:        parts[1],
						Kind:        kind,
						Signature:   trimmed,
						Description: docText,
						File:        filePath,
						Line:        lineNum,
					})
				}
			}
			continue
		}

		// 3. Funções e Métodos
		if strings.HasPrefix(trimmed, "funcao ") {
			rest := strings.TrimPrefix(trimmed, "funcao ")
			if strings.HasPrefix(rest, "(") {
				// Método com receptor: funcao (r Receptor) Nome() ...
				closeParen := strings.Index(rest, ")")
				if closeParen != -1 {
					methodRest := strings.TrimSpace(rest[closeParen+1:])
					methodParts := strings.Fields(methodRest)
					if len(methodParts) > 0 {
						nameWithArgs := methodParts[0]
						methodName := strings.Split(nameWithArgs, "(")[0]
						symbols = append(symbols, DocSymbol{
							Name:        methodName,
							Kind:        "metodo",
							Signature:   strings.TrimSuffix(trimmed, "{"),
							Description: docText,
							File:        filePath,
							Line:        lineNum,
						})
					}
				}
			} else {
				// Função regular: funcao Nome(args) ...
				parts := strings.Fields(rest)
				if len(parts) > 0 {
					fnName := strings.Split(parts[0], "(")[0]
					symbols = append(symbols, DocSymbol{
						Name:        fnName,
						Kind:        "funcao",
						Signature:   strings.TrimSuffix(trimmed, "{"),
						Description: docText,
						File:        filePath,
						Line:        lineNum,
					})
				}
			}
			continue
		}
	}

	return symbols
}

func generateHTML(symbols []DocSymbol, outputPath, rootPath string) error {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8">
  <title>Documentação da API | VAMOS-LANG</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #0d1117;
      --card-bg: #161b22;
      --border: #30363d;
      --text: #c9d1d9;
      --accent: #2ea043;
      --accent-hover: #3fb950;
      --keyword: #ff7b72;
      --type: #79c0ff;
    }
    body {
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      background: var(--bg);
      color: var(--text);
      margin: 0;
      padding: 0;
      display: flex;
      min-height: 100vh;
    }
    aside {
      width: 280px;
      background: var(--card-bg);
      border-right: 1px solid var(--border);
      padding: 24px;
      position: sticky;
      top: 0;
      height: 100vh;
      box-sizing: border-box;
      overflow-y: auto;
    }
    main {
      flex: 1;
      padding: 40px 60px;
      max-width: 1000px;
    }
    h1 { color: #fff; margin-top: 0; font-size: 28px; }
    h2 { color: var(--accent-hover); border-bottom: 1px solid var(--border); padding-bottom: 8px; margin-top: 40px; }
    .badge {
      display: inline-block;
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      padding: 2px 8px;
      border-radius: 12px;
      background: #238636;
      color: #fff;
      margin-left: 8px;
    }
    .badge-tipo { background: #1f6feb; }
    .badge-metodo { background: #8957e5; }
    .card {
      background: var(--card-bg);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 20px;
      margin-bottom: 24px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    }
    .symbol-title {
      font-size: 20px;
      font-weight: 600;
      color: #fff;
      display: flex;
      align-items: center;
      margin-bottom: 12px;
    }
    pre {
      background: #090d13;
      padding: 12px 16px;
      border-radius: 6px;
      overflow-x: auto;
      font-family: 'JetBrains Mono', monospace;
      font-size: 13px;
      border: 1px solid var(--border);
      color: var(--type);
    }
    .file-ref { font-size: 12px; color: #8b949e; margin-top: 8px; }
    a { color: #58a6ff; text-decoration: none; }
    a:hover { text-decoration: underline; }
    ul.nav-list { list-style: none; padding: 0; margin: 0; }
    ul.nav-list li { margin: 8px 0; }
    .search-box {
      width: 100%;
      padding: 8px 12px;
      background: #090d13;
      border: 1px solid var(--border);
      border-radius: 6px;
      color: #fff;
      margin-bottom: 16px;
      box-sizing: border-box;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <aside>
    <div style="font-weight:700; font-size:18px; color:#fff; margin-bottom:16px;">🇧🇷 VAMOS-LANG Doc</div>
    <input type="text" id="filterInput" class="search-box" placeholder="Filtrar símbolos..." onkeyup="filterDocs()">
    <ul class="nav-list" id="navList">
`)

	for _, s := range symbols {
		sb.WriteString(fmt.Sprintf(`      <li><a href="#%s">%s</a> <span style="font-size:10px; color:#8b949e;">(%s)</span></li>`+"\n",
			html.EscapeString(s.Name), html.EscapeString(s.Name), s.Kind))
	}

	sb.WriteString(fmt.Sprintf(`    </ul>
  </aside>
  <main>
    <h1>Documentação Oficial do Projeto</h1>
    <p style="color:#8b949e;">Gerado automaticamente a partir dos fontes VAMOS-LANG em <code>%s</code>.</p>
`, html.EscapeString(rootPath)))

	for _, s := range symbols {
		badgeClass := "badge"
		if s.Kind == "estrutura" || s.Kind == "interface" {
			badgeClass = "badge badge-tipo"
		} else if s.Kind == "metodo" {
			badgeClass = "badge badge-metodo"
		}

		desc := s.Description
		if desc == "" {
			desc = "<em>Nenhuma descrição fornecida.</em>"
		} else {
			desc = html.EscapeString(desc)
		}

		sb.WriteString(fmt.Sprintf(`    <div class="card doc-item" id="%s" data-name="%s">
      <div class="symbol-title">
        <span>%s</span>
        <span class="%s">%s</span>
      </div>
`, html.EscapeString(s.Name), strings.ToLower(s.Name), html.EscapeString(s.Name), badgeClass, s.Kind))

		if s.Signature != "" {
			sb.WriteString(fmt.Sprintf(`      <pre><code>%s</code></pre>`+"\n", html.EscapeString(s.Signature)))
		}

		sb.WriteString(fmt.Sprintf(`      <p>%s</p>
      <div class="file-ref">Definido em: <code>%s:%d</code></div>
    </div>
`, desc, html.EscapeString(s.File), s.Line))
	}

	sb.WriteString(`  </main>
  <script>
    function filterDocs() {
      const q = document.getElementById('filterInput').value.toLowerCase();
      const items = document.querySelectorAll('.doc-item');
      items.forEach(el => {
        const name = el.getAttribute('data-name');
        el.style.display = name.includes(q) ? 'block' : 'none';
      });
    }
  </script>
</body>
</html>`)

	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
