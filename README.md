# 🚀 VAMOS-LANG

> **A linguagem de programação Go (Golang) com sintaxe 100% traduzida para Português do Brasil (PT-BR).**

[![Linguagem Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Licença MIT](https://img.shields.io/badge/Licença-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Produção-brightgreen.svg)]()

---

## 📖 Visão Geral

**VAMOS-LANG** é uma linguagem de programação de alto desempenho e tipagem estática com sintaxe idêntica ao Go, porém totalmente expressa no idioma Português (PT-BR). 

O ecossistema VAMOS-LANG inclui:
- **Transpilador Ultra-Rápido**: Em Go puro, sem dependências externas pesadas.
- **CLI Completa (`vamos`)**: Com comandos e aliases em Português e Inglês.
- **Extensão Oficial para VS Code**: Realce de sintaxe oficial e snippets inteligentes.
- **Formatador Oficial de Código (`vamos fmt`)**: Padronização automática de código `.vamos`.
- **Linter Estático & Validador de Sintaxe (`vamos lint --strict`)**: Validação de convenções e diagnóstico sintático baseado no Go AST.
- **Language Server Protocol Nativo (`vamos lsp`)**: Suporte a autocompletar inteligente, hover com documentação, formatação e diagnósticos em tempo real no VS Code/Cursor.
- **Gerenciador de Dependências Externas (`vamos get`)**: Instalação e resolução de pacotes de terceiros.
- **Conversor Bidirecional (`vamos go2vamos`)**: Converte código Go nativo para VAMOS-LANG.
- **REPL Interativo (`vamos repl`)**: Shell de linha de comando para testar expressões instantaneamente.
- **Playground Web com WebAssembly (Wasm)**: Ambiente interativo online rodando o motor Go nativo diretamente no navegador (`vamos.wasm`).
- **Gerador Automático de Documentação (`vamos doc`)**: Extração de documentação no terminal e exportação em HTML navegável.
- **Rastreamento Didático (`vamos trace`)**: Execução passo a passo detalhando linhas, fluxo e variáveis para fins educacionais.
- **Micro-Framework Web Integrada (`pacotes/servidor`)**: Roteamento rápido, middlewares de log/CORS e respostas JSON/Texto em Português.

---

## 💻 Comandos da CLI (`vamos`)

A CLI oferece suporte completo com comandos e aliases em Português e Inglês:

| Comando / Aliases | Descrição |
| :--- | :--- |
| `vamos run` \| `rodar` \| `executar <arquivo.vamos> [args...]` | Transpila e executa diretamente com `go run` |
| `vamos build` \| `compilar` \| `construir [flags...]` | Transpila o projeto e compila para binário nativo executável |
| `vamos test` \| `testar` \| `teste [flags...]` | Transpila o projeto e executa os testes (`go test`) |
| `vamos transpile` \| `converter` \| `transpilar <origem> [-o]` | Converte arquivo `.vamos` para Go puro (`.go`) |
| `vamos go2vamos` \| `descompilar <arquivo.go> [-o]` | Converte código Go nativo de volta para `.vamos` |
| `vamos fmt` \| `formatar [caminho...]` | Formata automaticamente os arquivos `.vamos` |
| `vamos lint` \| `checar` \| `verificar [--strict] [caminho...]` | Validador de estilo e verificador estrito de sintaxe |
| `vamos init` \| `iniciar <nome> [--template=<api\|cli\|jogo>]` | Gera projeto a partir de modelos prontos |
| `vamos doc` \| `documentacao [simbolo] [--html]` | Visualizador e gerador de documentação oficial em HTML |
| `vamos trace` \| `rastrear <arquivo.vamos>` | Executa em modo didático com rastreamento visual passo a passo |
| `vamos get` \| `obter` \| `instalar <pacote-go>` | Instala dependências e módulos externos |
| `vamos lsp` \| `servidor-linguagem` | Inicia o servidor oficial LSP (Language Server Protocol) |
| `vamos repl` \| `interativo` | Inicia o terminal interativo (REPL) |
| `vamos playground` \| `web [porta]` | Inicia o servidor local do Playground Web |
| `vamos version` \| `versao` | Exibe detalhes da versão e compilador |
| `vamos help` \| `ajuda` | Exibe a tela de ajuda com todos os comandos |

---

## 🎨 Extensão para VS Code & Cursor

O projeto inclui uma extensão oficial na pasta [`vscode-extension/`](file:///storage/www/projetos/VAMOSLANG/vscode-extension):
- **Realce de Sintaxe Oficial**: TextMate Grammar para palavras-chave, tipos, literais e métodos.
- **Snippets Inteligentes**: Digite `pkgm`, `funcao`, `se`, `sesenao`, `para`, `intervalo`, `estrutura`, `interface`, `disparar`, `adiar`, `imp` ou `impf`.
- **Configuração de Linguagem**: Auto-fechamento de chaves, parênteses e indentação automática.

### Como Instalar:
```bash
cp -r /storage/www/projetos/VAMOSLANG/vscode-extension ~/.vscode/extensions/vamos-lang
# ou no Cursor:
cp -r /storage/www/projetos/VAMOSLANG/vscode-extension ~/.cursor/extensions/vamos-lang
```

---

## 🌐 Playground Web Interativo

O **Playground Web** na pasta [`playground/`](file:///storage/www/projetos/VAMOSLANG/playground) roda diretamente no navegador:
- **Transpilação em Tempo Real**: Visualize o código Go correspondente à medida que digita.
- **Conversão Reversa Go ➔ VAMOS**: Cole código Go para convertê-lo instantaneamente para VAMOS.
- **Compartilhamento por Link**: Gere URLs compartilháveis (`#code=...`) com seu código embutido.
- **Simulador de Execução**: Veja a saída formatada do programa no console integrado.
- **Formatador Online**: Botão para formatar o código no navegador.
- **Inspetor Léxico de Tokens**: Tabela detalhada de tokens identificados pelo scanner.

Para iniciar:
```bash
./vamos web 8080
```
Acesse em: `http://localhost:8080`

---

## 📊 Tabela de Referência de Sintaxe (Go vs VAMOS)

### 1. Palavras Reservadas (Keywords)

| VAMOS-LANG | Golang | Descrição / Uso |
| :--- | :--- | :--- |
| `pacote` | `package` | Declaração do pacote |
| `principal` | `main` | Ponto de entrada padrão (`pacote principal`, `funcao principal()`) |
| `importar` | `import` | Importação de pacotes |
| `funcao` | `func` | Declaração de função / método |
| `variavel` | `var` | Declaração de variável |
| `constante` | `const` | Declaração de constante |
| `tipo` | `type` | Declaração de tipo |
| `estrutura` | `struct` | Declaração de struct |
| `interface` | `interface` | Declaração de interface |
| `mapa` | `map` | Tipo mapa associativo |
| `canal` | `chan` | Tipo canal de comunicação |
| `se` | `if` | Condicional |
| `senao` | `else` | Ramo alternativo da condicional |
| `para` | `for` | Laço de repetição |
| `intervalo` | `range` | Iteração sobre coleções/canais |
| `escolha` | `switch` | Estrutura de seleção múltipla |
| `caso` | `case` | Ramo do switch / select |
| `padrao` | `default` | Ramo padrão do switch / select |
| `prosseguir` | `fallthrough` | Continuação no próximo caso |
| `interromper` | `break` | Interrupção de laço / escolha |
| `continuar` | `continue` | Próxima iteração do laço |
| `ir_para` | `goto` | Salto incondicional |
| `retornar` | `return` | Retorno de função |
| `disparar` | `go` | Execução concorrente (goroutine) |
| `adiar` | `defer` | Execução adiada ao final da função |
| `selecionar` | `select` | Multiplexação de canais |

---

### 2. Tipos Primitivos, Valores e Funções Embutidas

| VAMOS-LANG | Golang | Descrição |
| :--- | :--- | :--- |
| `texto` | `string` | Cadeia de caracteres UTF-8 |
| `inteiro` | `int` | Inteiro padrão da arquitetura |
| `inteiro64` | `int64` | Inteiro de 64 bits |
| `decimal` / `decimal64` | `float64` | Ponto flutuante de 64 bits |
| `booleano` | `bool` | Valor booleano |
| `verdadeiro` | `true` | Literal booleano verdadeiro |
| `falso` | `false` | Literal booleano falso |
| `nulo` | `nil` | Ponteiro / valor nulo |
| `erro` | `error` | Tipo de interface de erro |
| `criar` | `make` | Aloca slices, mapas e canais |
| `adicionar` / `anexar` | `append` | Anexa itens a uma fatia (slice) |
| `tamanho` | `len` | Retorna o tamanho da coleção |
| `capacidade` | `cap` | Retorna a capacidade da coleção |
| `fechar` | `close` | Fecha um canal |
| `excluir` / `deletar` | `delete` | Remove chave de um mapa |
| `panico` | `panic` | Interrompe o fluxo normal por erro crítico |
| `recuperar` | `recover` | Captura pânicos dentro de funções adiadas |

---

### 3. Biblioteca Padrão (Pacotes e Métodos Comuns)

| VAMOS-LANG | Golang | Exemplo em VAMOS |
| :--- | :--- | :--- |
| `formatar.ImprimirLinha(...)` | `fmt.Println(...)` | `formatar.ImprimirLinha("Olá!")` |
| `formatar.ImprimirFormatado(...)` | `fmt.Printf(...)` | `formatar.ImprimirFormatado("%d\n", 42)` |
| `formatar.CriarErro(...)` | `fmt.Errorf(...)` | `err := formatar.CriarErro("falha: %v", e)` |
| `formatar.Formatar(...)` | `fmt.Sprintf(...)` | `txt := formatar.Formatar("%s", v)` |
| `tempo.Dormir(...)` | `time.Sleep(...)` | `tempo.Dormir(tempo.Segundo)` |
| `tempo.Agora()` | `time.Now()` | `agora := tempo.Agora()` |
| `sincronizar.GrupoEspera` | `sync.WaitGroup` | `variavel wg sincronizar.GrupoEspera` |
| `erros.Novo(...)` | `errors.New(...)` | `retornar erros.Novo("falha")` |
| `so.Sair(...)` | `os.Exit(...)` | `so.Sair(0)` |
| `conversao.ParaInteiro(...)` | `strconv.Atoi(...)` | `num, _ := conversao.ParaInteiro("42")` |
| `cripto/sha256.Soma256(...)` | `crypto/sha256.Sum256(...)` | `hash := sha256.Soma256([]byte("msg"))` |
| `matematica/aleatorio.InteiroAte(...)` | `math/rand.Intn(...)` | `sorte := aleatorio.InteiroAte(100)` |
| `banco_dados.Abrir(...)` | `database/sql.Open(...)` | `db, _ := banco_dados.Abrir("sqlite", "db.sqlite")` |

---

## 🚀 Exemplo Completo (`teste.vamos`)

```vamos
pacote principal

importar (
	"formatar"
	"sincronizar"
	"tempo"
)

tipo MensagemProcessamento estrutura {
	Identificador inteiro
	Descricao     texto
	Pontuacao     decimal64
	Prioritaria   booleano
}

funcao (m MensagemProcessamento) ExibirDetalhes() {
	tagPrioridade := "NORMAL"
	se m.Prioritaria {
		tagPrioridade = "ALTA"
	}
	formatar.ImprimirFormatado(
		"  [Item %d] %-25s | Pontuação: %.2f | Prioridade: %s\n",
		m.Identificador,
		m.Descricao,
		m.Pontuacao,
		tagPrioridade,
	)
}

funcao executarTrabalhador(id inteiro, canalSaida canal MensagemProcessamento, wg *sincronizar.GrupoEspera) {
	adiar wg.Concluido()
	tempo.Dormir(tempo.Milissegundo * tempo.Duracao(100*id))

	msg := MensagemProcessamento{
		Identificador: id,
		Descricao:     formatar.Formatar("Payload da tarefa #%d", id),
		Pontuacao:     decimal64(id) * 15.5,
		Prioritaria:   (id%2 != 0),
	}

	canalSaida <- msg
}

funcao principal() {
	formatar.ImprimirLinha("==================================================")
	formatar.ImprimirLinha("         VAMOS-LANG: Teste Integrado Completo     ")
	formatar.ImprimirLinha("==================================================")

	constante totalTarefas inteiro = 4
	canalDados := criar(canal MensagemProcessamento, totalTarefas)
	variavel wg sincronizar.GrupoEspera

	formatar.ImprimirLinha("1. Disparando tarefas concorrentes:")
	para i := 1; i <= totalTarefas; i++ {
		wg.Adicionar(1)
		disparar executarTrabalhador(i, canalDados, &wg)
	}

	disparar funcao() {
		wg.Esperar()
		fechar(canalDados)
	}()

	formatar.ImprimirLinha("\n2. Recebendo e exibindo dados processados:")
	para item := intervalo canalDados {
		item.ExibirDetalhes()
	}
	formatar.ImprimirLinha("==================================================")
}
```

**Executando:**
```bash
./vamos rodar teste.vamos
```

---

## 🎮 Jogos Clássicos Implementados em VAMOS-LANG

O repositório inclui implementações completas de jogos clássicos na pasta [`jogos/`](file:///storage/www/projetos/VAMOSLANG/jogos):
- 🐍 **Jogo da Cobrinha** ([`jogos/cobrinha.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/cobrinha.vamos)): Grade 2D, crescimento dinâmico e IA com busca autônoma de alimentos.
- 👾 **Pacman / Come-Come** ([`jogos/pacman.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/pacman.vamos)): Labirinto clássico, pastilhas, pílulas de poder temporárias e 3 fantasmas com perseguição vetorial e fuga.
- 🏁 **Jogo de Damas** ([`jogos/damas.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/damas.vamos)): Tabuleiro 8x8 oficial, capturas obrigatórias por salto, promoção a Damas coroadas e IA heurística (1P vs CPU ou Demo).

Consulte [`jogos/README.md`](file:///storage/www/projetos/VAMOSLANG/jogos/README.md) para regras completas e modos de execução.

```bash
# Executar demonstrações com IA autônoma:
./vamos run jogos/cobrinha.vamos --demo
./vamos run jogos/pacman.vamos --demo
./vamos run jogos/damas.vamos --demo

# Ou rodar via make:
make run-games
```

---

## 🏗️ Estrutura do Projeto

```
/storage/www/projetos/VAMOSLANG/
├── main.go                   # Ponto de entrada da CLI oficial
├── teste.vamos               # Exemplo de teste integrado oficial
├── lsp/                      # Language Server Protocol (LSP JSON-RPC 2.0 nativo)
├── jogos/                    # Jogos clássicos desenvolvidos em VAMOS-LANG
│   ├── cobrinha.vamos        # Jogo da Cobrinha (Snake) com IA e interativo
│   ├── pacman.vamos          # Pacman com labirinto e fantasmas
│   ├── damas.vamos           # Jogo de Damas 8x8 com IA tática
│   └── README.md             # Guia completo dos jogos
├── vscode-extension/         # Extensão oficial para VS Code / Cursor
│   ├── package.json
│   ├── extension.js          # Cliente integrado ao vamos lsp
│   ├── language-configuration.json
│   ├── syntaxes/
│   │   └── vamos.tmLanguage.json
│   ├── snippets/
│   │   └── vamos.json
│   └── README.md
├── playground/               # Playground Web Interativo
│   ├── index.html
│   ├── style.css
│   ├── transpiler.js
│   ├── app.js
│   ├── vamos.wasm            # Motor Wasm Go nativo compilado
│   ├── wasm_exec.js          # Ponte de execução WebAssembly
│   └── README.md
├── blocodenotasBR/           # Aplicação real completa (Bloco de Notas CLI + JSON)
├── pacotes/                  # Pacotes e bibliotecas nativas em PT-BR
│   └── servidor/             # Micro-framework web com roteamento e middlewares
├── exemplos/                 # Coleção de exemplos práticos
│   ├── ola_mundo.vamos
│   ├── concorrencia.vamos
│   ├── estruturas_interfaces.vamos
│   ├── calculadora.vamos
│   ├── servidor_web.vamos
│   ├── cripto_e_conversao.vamos
│   ├── nano.vamos            # Editor de texto estilo GNU nano em terminal
│   └── api_rest.vamos        # API REST com a micro-framework servidor
├── wasm/                     # Ponto de entrada para compilação WebAssembly
├── transpiler/               # Motor do transpilador, validador Go AST e linter
├── cmd/                      # Comandos modulares da CLI (run, build, lsp, get, etc.)
├── go.mod                    # Definição do módulo Go
├── Makefile                  # Automação de compilação, testes e demonstrações
└── README.md                 # Documentação completa
```

---

## 📄 Licença

Este projeto está licenciado sob a licença [MIT](LICENSE).
