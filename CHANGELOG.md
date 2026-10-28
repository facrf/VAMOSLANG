# 📝 Registro de Alterações (CHANGELOG) - VAMOS-LANG

Todas as alterações notáveis, novas funcionalidades, melhorias de arquitetura e correções de bugs deste projeto são documentadas neste arquivo.

---

## [1.1.0] - 28/10/2026

### 🎮 Novos Jogos e Melhorias em Arcade
- **Novo Jogo de Xadrez Clássico (`jogos/xadrez.vamos`)**:
  - Implementação completa das regras clássicas da FIDE para o terminal em 100% VAMOS-LANG.
  - Movimentos oficiais para todas as peças: Peão (avanço simples, avanço duplo inicial, captura diagonal e promoção automática à Dama), Cavalo (em "L"), Bispo (diagonais), Torre (linhas/colunas), Dama e Rei.
  - Detecção rigorosa de Xeque (`EstaEmXeque`) e filtragem de lances ilegais (`LanceDeixaReiEmXeque`).
  - Reconhecimento automático de Xeque-Mate e Empate por Afogamento (*Stalemate*).
  - Notação algébrica amigável no formato de coordenadas reais (ex: `e2 e4` ou `e2e4`).
  - Motor de Inteligência Artificial com avaliação de valor material, controle do centro e priorização de capturas.
  - Modo demonstração autônomo com duelo entre IAs (`vamos run jogos/xadrez.vamos --demo`).
- **Aprimoramento do Jogo da Cobrinha (`jogos/cobrinha.vamos`)**:
  - Nova mecânica de Fruta Dourada Especial (`✪`) que surge a cada 4 maçãs comuns com temporizador regressivo, concedendo bônus de **+50 pontos**.
  - IA tática refinada com cálculo de espaços livres adjacentes para evitar confinamentos.
  - Níveis de velocidade configuráveis (`--rapido`, `--insano`) e sistema de recorde por sessão.
- **Automação no Makefile (`make run-games`)**:
  - Execução contínua e sequencial dos 4 jogos clássicos: Cobrinha, Pacman, Damas e Xadrez.

---

### 🧪 Suporte de Primeira Classe a Testes Unitários Nativos (`*_teste.vamos`)
- **Pacote Canônico `teste`**:
  - Permite criar suítes de testes 100% em código VAMOS puro com arquivos terminando em `_teste.vamos`.
  - `importar "teste"` mapeado diretamente para o pacote padrão `testing` do Go.
  - Suporte completo a assinaturas e tipos: `t *teste.T` e `b *teste.B`.
  - Mapeamento transparente de todos os métodos de teste:
    - `t.Erro`, `t.ErroFormatado`, `t.Falhar`, `t.FalharAgora`, `t.Falhou`, `t.Fatal`, `t.FatalFormatado`.
    - `t.Registrar`, `t.RegistrarFormatado`, `t.Nome`, `t.Pular`, `t.PularFormatado`, `t.PularAgora`, `t.Pulou`.
    - `t.Executar` (subtestes), `t.Limpeza`, `t.Auxiliar` / `t.Ajudante`, `t.Paralelo`.
  - Conversão inteligente de nomenclatura de funções para atender à convenção estrita do Go:
    - `funcao TesteX(...)` ➔ `func TestX(...)`
    - `funcao ReferenciaX(...)` ➔ `func BenchmarkX(...)`
    - `funcao ExemploX(...)` ➔ `func ExampleX(...)`
- **Reescrita e Diagnóstico Amigável**:
  - Transpilação isolada para `_test.go` durante a preparação do workspace de testes.
  - `RewriteErrors` remapeia automaticamente linhas e colunas de falhas apontando para os arquivos `_teste.vamos` originais.
- **Novas Suítes Nativas na Biblioteca Padrão**:
  - [`pacotes/matematica_util/matematica_util_teste.vamos`](pacotes/matematica_util/matematica_util_teste.vamos): Testes das funções numéricas e matemáticas em VAMOS.
  - [`pacotes/texto_util/texto_util_teste.vamos`](pacotes/texto_util/texto_util_teste.vamos): Testes de manipulação de cadeias de caracteres em VAMOS.
  - [`pacotes/colecoes/colecoes_teste.vamos`](pacotes/colecoes/colecoes_teste.vamos): Testes de estruturas de dados genéricas (`Pilha[texto]`, `Fila[inteiro]`, `Conjunto[texto]`) em VAMOS.

---

### ⚙️ Compilador, Máquina Virtual (VVM) e Bytecode
- **Operador Módulo / Resto de Divisão (`%`)**:
  - Adicionado opcode `OpModulo = 30` em `compilador/bytecode.vamos`.
  - Emissão da AST e execução na VVM com proteção contra divisão por zero.
- **Aritmética e Comparações em Ponto Flutuante (`decimal64`) na VVM**:
  - Operações aritméticas (`OpSomar`, `OpSubtrair`, `OpMultiplicar`, `OpDividir`) e relacionais (`OpIgual`, `OpDiferente`, `OpMenor`, `OpMaior`, `OpMenorIgual`, `OpMaiorIgual`) com detecção automática de literais decimais via `conversao.AnalisarDecimal`.
- **Backend WebAssembly Nativo 1.0 (`compilador/wasm.vamos`)**:
  - Adicionada Seção 5 (Memória Linear com 1 página inicial) e exportação `"memoria"` na Seção 7.
  - Adicionado opcode de módulo assinado `WasmOpI32RemS = 0x6F`.
- **Transpilador Auto-Hospedado (`compilador/transpilador.vamos`)**:
  - Sincronizado com os métodos do pacote `teste`, métodos de instância de objetos e suporte a conversão de identificadores de testes.

---

### 🧠 Language Server Protocol (LSP)
- **Acesso Concorrente Thread-Safe**:
  - Sincronização do mapa de documentos abertos com `sync.RWMutex` (`getDoc`, `setDoc`, `deleteDoc`).
- **Navegação de Código (Go to Definition)**:
  - Implementado o método `textDocument/definition` rastreando declarações de funções, métodos e estruturas customizadas.
- **Visão de Estrutura de Documento (Document Symbols / Outline View)**:
  - Implementado o método `textDocument/documentSymbol` mapeando pacotes, funções e estruturas para o painel Outline do VS Code e Cursor.
- **Suíte de Testes LSP**:
  - Adicionados testes unitários automatizados em `lsp/server_test.go` cobrindo inicialização, hover, autocompletar, formatação, definição e símbolos.

---

### 🖥️ CLI e Ferramentas do Desenvolvedor
- **Isolamento de Workspace e Suporte Multi-Executável (`cmd/workspace.go`)**:
  - Adicionada contagem de pontos de entrada por diretório (`GetDirMainCount`) para evitar colisões de `package main` em pastas contendo múltiplos exemplos ou scripts independentes.
  - Preservação rigorosa de permissões executáveis (`si.Mode()`) na cópia de arquivos.
  - Preservação dos arquivos `.vamos` originais na pasta temporária do workspace.
- **Comando de Teste (`cmd/test.go`)**:
  - Suporte a flags arbitrárias (`vamos test -v`, `vamos test -race`), adicionando `./...` automaticamente quando nenhum pacote específico é informado.
- **Rastreamento Didático (`cmd/trace.go`)**:
  - Inclusão da cópia de `go.sum` para a pasta temporária de rastreamento, viabilizando o uso de dependências externas como `github.com/google/uuid`.
- **Gerador de Projetos (`cmd/init.go`)**:
  - Template padrão atualizado para gerar automaticamente arquivos de teste nativos (`pacotes/calculos/calculos_teste.vamos`) e instruções para `vamos test`.
- **Terminal Interativo REPL (`cmd/repl.go`)**:
  - Suporte a declaração de tipos estruturados (`tipo ... estrutura`) e funções (`funcao ...`) de nível superior na mesma sessão.
  - Importação dinâmica via `importar "pacote"` e auto-importação dos módulos canônicos.
  - Tradução e filtragem automática de mensagens de erro do compilador Go em tempo de execução.

---

### 📚 Documentação e Infraestrutura
- **Documentação Oficial**:
  - [`README.md`](README.md): Adicionada a seção completa de Testes Unitários Nativos, novos comandos da CLI e tabela atualizada da biblioteca padrão.
  - [`jogos/README.md`](jogos/README.md): Documentação detalhada dos 4 jogos clássicos, mecânicas e comandos de compilação para binários nativos.
  - [`docs.html`](docs.html): Regenerado via `vamos doc --html`.
- **Arquivo `.gitignore`**:
  - Configurado para ignorar binários executáveis (`vamosc*`, `*.vamosb`, `*.wasm`, `*.o`), pastas temporárias de workspace e banco de dados de notas.
