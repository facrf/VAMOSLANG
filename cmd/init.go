package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunInit executa o comando 'vamos init <nome-do-modulo> [--template=<api|cli|jogo|padrao>]'
func RunInit(args []string) int {
	var targetInput string
	template := "padrao"

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--template=") {
			template = strings.ToLower(strings.TrimPrefix(arg, "--template="))
		} else if strings.HasPrefix(arg, "--modelo=") {
			template = strings.ToLower(strings.TrimPrefix(arg, "--modelo="))
		} else if arg == "-t" || arg == "--template" || arg == "-m" || arg == "--modelo" {
			if i+1 < len(args) {
				template = strings.ToLower(args[i+1])
				i++
			}
		} else if !strings.HasPrefix(arg, "-") && targetInput == "" {
			targetInput = arg
		}
	}

	if targetInput == "" {
		fmt.Fprintf(os.Stderr, "Erro: informe o nome do módulo para o novo projeto.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos init <nome-do-modulo> [--template=api|cli|jogo|padrao]\n")
		return 1
	}

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

	modName = strings.TrimPrefix(modName, ".")
	modName = strings.TrimPrefix(modName, "/")
	if modName == "" {
		modName = "app"
	}

	fmt.Printf("==> Criando novo projeto VAMOS-LANG: %s [Modelo: %s]\n", modName, template)

	// go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.20\n", modName)
	_ = os.MkdirAll(projectDir, 0755)
	_ = os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644)

	// .gitignore
	gitIgnore := "# Binários e temporários\nbin/\n*.exe\n*.out\napp\nnotas_db.json\n"
	_ = os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitIgnore), 0644)

	switch template {
	case "api", "web", "rest":
		criarTemplateAPI(projectDir, modName)
	case "cli", "terminal":
		criarTemplateCLI(projectDir, modName)
	case "jogo", "game":
		criarTemplateJogo(projectDir, modName)
	default:
		criarTemplatePadrao(projectDir, modName)
	}

	fmt.Println("✓ Estrutura de pastas gerada com sucesso!")
	fmt.Printf("\nPara começar a desenvolver:\n")
	if projectDir != "." {
		fmt.Printf("  cd %s\n", projectDir)
	}
	fmt.Printf("  vamos run cmd/principal.vamos\n\n")
	return 0
}

func criarTemplatePadrao(projectDir, modName string) {
	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "calculos"), 0755)

	mainContent := fmt.Sprintf(`// Ponto de entrada padrão do projeto
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

	calcContent := `pacote calculos

funcao Somar(a inteiro, b inteiro) inteiro {
	retornar a + b
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "calculos", "calculos.vamos"), []byte(calcContent), 0644)
}

func criarTemplateAPI(projectDir, modName string) {
	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "controladores"), 0755)

	mainContent := fmt.Sprintf(`// API REST em VAMOS-LANG
pacote principal

importar (
	"formatar"
	"servidor_http"
	"%s/pacotes/controladores"
)

funcao principal() {
	porta := ":8080"
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirLinha("       🚀 API REST VAMOS-LANG           ")
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirFormatado("Servidor rodando em http://localhost%%s\n", porta)
	formatar.ImprimirLinha("Rotas disponíveis:")
	formatar.ImprimirLinha("  GET  /              -> Mensagem de boas-vindas")
	formatar.ImprimirLinha("  GET  /api/status    -> Status de saúde da API")
	formatar.ImprimirLinha("  GET  /api/tarefas   -> Lista de tarefas em JSON")

	servidor_http.Manipulador("/", controladores.InicioHandler)
	servidor_http.Manipulador("/api/status", controladores.StatusHandler)
	servidor_http.Manipulador("/api/tarefas", controladores.TarefasHandler)

	err := servidor_http.OuvirEServir(porta, nulo)
	se err != nulo {
		formatar.ImprimirFormatado("Erro no servidor: %%v\n", err)
	}
}
`, modName)
	_ = os.WriteFile(filepath.Join(projectDir, "cmd", "principal.vamos"), []byte(mainContent), 0644)

	ctrlContent := `pacote controladores

importar (
	"formatar"
	"json"
	"servidor_http"
	"tempo"
)

tipo Tarefa estrutura {
	Id        inteiro ` + "`json:\"id\"`" + `
	Titulo    texto   ` + "`json:\"titulo\"`" + `
	Concluida booleano ` + "`json:\"concluida\"`" + `
}

funcao InicioHandler(w servidor_http.EscritorResposta, r *servidor_http.Requisicao) {
	formatar.FimprimirLinha(w, "Bem-vindo à API desenvolvida 100% em VAMOS-LANG!")
}

funcao StatusHandler(w servidor_http.EscritorResposta, r *servidor_http.Requisicao) {
	w.Header().Set("Content-Type", "application/json")
	resposta := mapa[texto]texto{
		"status":    "operacional",
		"linguagem": "VAMOS-LANG (Go em PT-BR)",
		"horario":   tempo.Agora().Formato("02/01/2006 15:04:05"),
	}
	dados, _ := json.Serializar(resposta)
	w.Write(dados)
}

funcao TarefasHandler(w servidor_http.EscritorResposta, r *servidor_http.Requisicao) {
	w.Header().Set("Content-Type", "application/json")
	tarefas := []Tarefa{
		{Id: 1, Titulo: "Escrever código em VAMOS-LANG", Concluida: verdadeiro},
		{Id: 2, Titulo: "Publicar API REST", Concluida: falso},
	}
	dados, _ := json.SerializarIdentado(tarefas, "", "  ")
	w.Write(dados)
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "controladores", "controladores.vamos"), []byte(ctrlContent), 0644)
}

func criarTemplateCLI(projectDir, modName string) {
	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "comandos"), 0755)

	mainContent := fmt.Sprintf(`// Ferramenta de Linha de Comando (CLI)
pacote principal

importar (
	"formatar"
	"so"
	"%s/pacotes/comandos"
)

funcao principal() {
	args := so.Argumentos[1:]
	se tamanho(args) == 0 {
		comandos.ExibirAjuda()
		retornar
	}

	cmd := args[0]
	escolha cmd {
	caso "saudacao", "ola":
		nome := "Mundo"
		se tamanho(args) > 1 {
			nome = args[1]
		}
		comandos.Saudar(nome)
	caso "info":
		comandos.ExibirInfo()
	padrao:
		formatar.ImprimirFormatado("Comando desconhecido: %%s\n", cmd)
		comandos.ExibirAjuda()
	}
}
`, modName)
	_ = os.WriteFile(filepath.Join(projectDir, "cmd", "principal.vamos"), []byte(mainContent), 0644)

	cmdContent := `pacote comandos

importar "formatar"

funcao ExibirAjuda() {
	formatar.ImprimirLinha("Uso: app <comando> [argumentos]")
	formatar.ImprimirLinha("Comandos:")
	formatar.ImprimirLinha("  saudacao [nome]   Exibe uma mensagem amigável")
	formatar.ImprimirLinha("  info              Exibe detalhes da aplicação")
}

funcao Saudar(nome texto) {
	formatar.ImprimirFormatado("Olá, %s! Bem-vindo ao VAMOS-LANG.\n", nome)
}

funcao ExibirInfo() {
	formatar.ImprimirLinha("Aplicação CLI construída com VAMOS-LANG (Go 100% PT-BR)")
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "comandos", "comandos.vamos"), []byte(cmdContent), 0644)
}

func criarTemplateJogo(projectDir, modName string) {
	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "grade"), 0755)

	mainContent := fmt.Sprintf(`// Jogo 2D em Terminal
pacote principal

importar (
	"formatar"
	"%s/pacotes/grade"
)

funcao principal() {
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirLinha("     🎮 JOGO EM TERMINAL VAMOS-LANG     ")
	formatar.ImprimirLinha("========================================")

	g := grade.NovaGrade(12, 6)
	g.Posicionar(5, 3, "★")
	g.Renderizar()
	formatar.ImprimirLinha("Jogo inicializado com sucesso!")
}
`, modName)
	_ = os.WriteFile(filepath.Join(projectDir, "cmd", "principal.vamos"), []byte(mainContent), 0644)

	gradeContent := `pacote grade

importar (
	"cordas"
	"formatar"
)

tipo Grade estrutura {
	Largura inteiro
	Altura  inteiro
	Matriz  [][]texto
}

funcao NovaGrade(l, a inteiro) *Grade {
	m := criar([][]texto, a)
	para y := 0; y < a; y++ {
		m[y] = criar([]texto, l)
		para x := 0; x < l; x++ {
			m[y][x] = "·"
		}
	}
	retornar &Grade{Largura: l, Altura: a, Matriz: m}
}

funcao (g *Grade) Posicionar(x, y inteiro, simbolo texto) {
	se x >= 0 && x < g.Largura && y >= 0 && y < g.Altura {
		g.Matriz[y][x] = simbolo
	}
}

funcao (g *Grade) Renderizar() {
	formatar.ImprimirLinha("┌" + cordas.Repetir("──", g.Largura) + "┐")
	para y := 0; y < g.Altura; y++ {
		linha := "│"
		para x := 0; x < g.Largura; x++ {
			linha += g.Matriz[y][x] + " "
		}
		linha += "│"
		formatar.ImprimirLinha(linha)
	}
	formatar.ImprimirLinha("└" + cordas.Repetir("──", g.Largura) + "┘")
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "grade", "grade.vamos"), []byte(gradeContent), 0644)
}
