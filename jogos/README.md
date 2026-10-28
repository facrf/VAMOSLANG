# 🎮 Jogos Clássicos em VAMOS-LANG

Uma suíte de jogos clássicos de arcade e tabuleiro implementados **100% em Português do Brasil (PT-BR)** com a linguagem [VAMOS-LANG](file:///storage/www/projetos/VAMOSLANG/README.md).

Estes jogos servem como demonstração prática e didática do poder, expressividade e compatibilidade do ecossistema VAMOS-LANG, abrangendo:
- Estruturas de dados (`tipo ... estrutura`) e métodos com receptores (`funcao (j *Jogo) ...`).
- Concorrência, temporização (`tempo.Dormir`) e aleatoriedade (`matematica/aleatorio`).
- Manipulação de fatias (slices), matrizes bidimensionais e strings (`cordas`).
- Entrada interativa via terminal (`leitor.NovoScanner(so.EntradaPadrao)`).
- Algoritmos de inteligência artificial (IA) para modos de demonstração e partidas contra o computador.

---

## 🕹️ Índice dos Jogos

1. [🐍 Jogo da Cobrinha (`cobrinha.vamos`)](#-1-jogo-da-cobrinha-cobrinhavamos)
2. [👾 Pacman / Come-Come (`pacman.vamos`)](#-2-pacman--come-come-pacmanvamos)
3. [🏁 Jogo de Damas (`damas.vamos`)](#-3-jogo-de-damas-damasvamos)
4. [♔ Jogo de Xadrez (`xadrez.vamos`)](#-4-jogo-de-xadrez-xadrezvamos)

---

## 🐍 1. Jogo da Cobrinha (`cobrinha.vamos`)

Arquivo: [`jogos/cobrinha.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/cobrinha.vamos)

O clássico jogo onde o jogador controla uma serpente que cresce a cada alimento consumido, desviando das bordas do cenário e do próprio corpo.

### 🌟 Destaques de Implementação:
- **Estrutura 2D**: Representação cartesiana através do tipo `Ponto { X inteiro; Y inteiro }`.
- **Corpo Dinâmico**: Fatias dinâmicas em VAMOS (`[]Ponto`) simulando avanço da cabeça (`adicionar`) e recolhimento da cauda (`j.Corpo[1:]`).
- **IA Autônoma para Demonstração**: Algoritmo de busca direcionada à maçã com prevenção antecipada de colisões.

### 🎮 Como Jogar:
```bash
# Modo Interativo (use W, A, S, D):
./vamos run jogos/cobrinha.vamos

# Modo Demonstração (IA joga automaticamente):
./vamos run jogos/cobrinha.vamos --demo
```

**Legenda do Tabuleiro:**
- `◉` Cabeça da cobrinha
- `○` Segmentos do corpo
- `★` Maçã comum (+10 pontos)
- `✪` Fruta dourada especial bônus (+50 pontos temporários)

---

## 👾 2. Pacman / Come-Come (`pacman.vamos`)

Arquivo: [`jogos/pacman.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/pacman.vamos)

Adaptação do clássico arcade dos anos 80. O Pacman percorre um labirinto consumindo pastilhas e desviando dos fantasmas (*Blinky*, *Pinky* e *Inky*).

### 🌟 Destaques de Implementação:
- **Matriz de Labirinto**: Tabuleiro `[][]runa` com paredes intransponíveis (`██`), corredores e pastilhas.
- **Pílulas de Poder (`●`)**: Ao coletar uma pílula especial, o Pacman ganha um período temporário onde os fantasmas entram no estado assustado (`?`) e podem ser devorados por +200 pontos.
- **Comportamento dos Fantasmas**: IA de perseguição vetorial em modo normal e de fuga calculada durante o estado de poder do Pacman.

### 🎮 Como Jogar:
```bash
# Modo Interativo (W: Cima, S: Baixo, A: Esquerda, D: Direita):
./vamos run jogos/pacman.vamos

# Modo Demonstração (Simulação automática de labirinto):
./vamos run jogos/pacman.vamos --demo
```

**Legenda do Labirinto:**
- `C` Pacman (jogador)
- `B`, `P`, `I` Fantasmas (*Blinky*, *Pinky*, *Inky*)
- `?` Fantasma vulnerável (assustado)
- `·` Pastilha simples (+10 pontos)
- `●` Pílula de poder (+50 pontos)
- `██` Paredes do labirinto

---

## 🏁 3. Jogo de Damas (`damas.vamos`)

Arquivo: [`jogos/damas.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/damas.vamos)

O tradicional jogo de estratégia em tabuleiro oficial 8x8. O jogador comanda as peças Brancas contra a inteligência artificial do Computador que gerencia as Pretas.

### 🌟 Destaques de Implementação:
- **Regras Oficiais**:
  - Movimentos diagonais simples de avanço.
  - Capturas obrigatórias por salto duplo diagonal com eliminação da peça adversária.
  - Promoção automática à Dama (`⛁` ou `⛃`) ao atingir o lado oposto do tabuleiro.
  - Damas possuem mobilidade multidirecional (frente e trás).
- **Notação Algébrica**: Entrada amigável no formato de coordenadas reais de xadrez e damas (ex: `c3 d4` ou lances de salto como `e3 g5`).
- **IA Heurística**: Avaliador tático que prioriza capturas, avanços estratégicos e defesas de coroação.

### 🎮 Como Jogar:
```bash
# Modo 1 Jogador vs Computador:
./vamos run jogos/damas.vamos

# Modo Demonstração (Partida IA vs IA):
./vamos run jogos/damas.vamos --demo
```

**Legenda das Peças:**
- `⛀` Peça Branca normal
- `⛁` Dama Branca (coroada)
- `⛂` Peça Preta normal
- `⛃` Dama Preta (coroada)
- `·` Casa escura jogável

---

## ♔ 4. Jogo de Xadrez (`xadrez.vamos`)

Arquivo: [`jogos/xadrez.vamos`](file:///storage/www/projetos/VAMOSLANG/jogos/xadrez.vamos)

O mais clássico jogo de estratégia do mundo, implementado com as regras oficiais da FIDE e notação algébrica.

### 🌟 Destaques de Implementação:
- **Peças e Movimentação Oficial**:
  - `♙`/`♟` Peões (passo simples, avanço duplo inicial de 2 casas, captura diagonal e promoção automática à Dama ao atingir a 8ª fileira).
  - `♘`/`♞` Cavalos (salto em "L" ultrapassando peças intervenientes).
  - `♗`/`♝` Bispos (diagonais livres).
  - `♖`/`♜` Torres (linhas e colunas livres).
  - `♕`/`♛` Damas (mobilidade combinada ortogonal e diagonal).
  - `♔`/`♚` Reis (passo unitário multidirecional).
- **Segurança do Rei e Xeque**:
  - Detecção rigorosa de Xeque (`EstaEmXeque`).
  - Filtragem de lances ilegais (`LanceDeixaReiEmXeque`) que exporiam o próprio Rei a ataque.
  - Reconhecimento automático de Xeque-Mate e Empate por Afogamento (*Stalemate*).
- **Notação Algébrica**: Entrada amigável aceitando formatos como `e2 e4` ou `e2e4`.
- **Inteligência Artificial Heurística**: Motor com avaliação de valor material, controle estratégico do centro e priorização de capturas.

### 🎮 Como Jogar:
```bash
# Modo Interativo (Humano Brancas vs Computador Pretas):
./vamos run jogos/xadrez.vamos

# Modo Demonstração (Duelo entre duas IAs táticas):
./vamos run jogos/xadrez.vamos --demo
```

---

## 🔨 Compilando os Jogos para Binários Nativos

Você pode compilar qualquer um dos jogos diretamente para um binário executável independente:

```bash
# Compilar a Cobrinha para binário executável nativo
./vamos build -o cobrinha jogos/cobrinha.vamos
./cobrinha --demo

# Compilar o Pacman
./vamos build -o pacman jogos/pacman.vamos
./pacman --demo

# Compilar Damas
./vamos build -o damas jogos/damas.vamos
./damas --demo

# Compilar Xadrez
./vamos build -o xadrez jogos/xadrez.vamos
./xadrez --demo
```

---

## 🔍 Verificação de Qualidade e Formatação

Todos os jogos seguem as convenções estritas da linguagem:

```bash
# Formatação automática de código
./vamos fmt jogos/

# Análise estática e validação rigorosa de sintaxe
./vamos lint --strict jogos/
```
