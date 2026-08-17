package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unicode"
)

// ============================================================================
// 1. CONSTANTES, VERSÃO E BANNER
// ============================================================================

const (
	Version  = "1.1.0"
	Codename = "Pindorama"
	Author   = "Equipe VAMOS-LANG"
)

const banner = `
 __      __     __  __  ____   _____     _           _   _  _____ 
 \ \    / /\   |  \/  |/ __ \ / ____|   | |         / \ | \ |/ ____|
  \ \  / /  \  | \  / | |  | | (___     | |        / _ \|  \| | |  __ 
   \ \/ / /\ \ | |\/| | |  | |\___ \    | |       / /_\ \ . ' | | |_ |
    \  / ____ \| |  | | |__| |____) |   | |____  / _____ \ |\  | |__| |
     \/_/    \_\_|  |_|\____/|_____/    |______|/_/     \_\_| \_|\_____|
`

// ============================================================================
// 2. TABELAS DE MAPEAMENTO ESTRITAS (PT-BR -> GO E REVERSO GO -> PT-BR)
// ============================================================================

// KeywordsMap mapeia palavras reservadas do VAMOS-LANG para Go.
var KeywordsMap = map[string]string{
	"pacote":      "package",
	"principal":   "main",
	"importar":    "import",
	"funcao":      "func",
	"variavel":    "var",
	"constante":   "const",
	"tipo":        "type",
	"estrutura":   "struct",
	"interface":   "interface",
	"mapa":        "map",
	"canal":       "chan",
	"se":          "if",
	"senao":       "else",
	"para":        "for",
	"intervalo":   "range",
	"escolha":     "switch",
	"caso":        "case",
	"padrao":      "default",
	"prosseguir":  "fallthrough",
	"interromper": "break",
	"continuar":   "continue",
	"ir_para":     "goto",
	"retornar":    "return",
	"disparar":    "go",
	"adiar":       "defer",
	"selecionar":  "select",
}

// ReverseKeywordsMap mapeia palavras reservadas do Go para VAMOS-LANG.
var ReverseKeywordsMap = map[string]string{
	"package":     "pacote",
	"main":        "principal",
	"import":      "importar",
	"func":        "funcao",
	"var":         "variavel",
	"const":       "constante",
	"type":        "tipo",
	"struct":      "estrutura",
	"interface":   "interface",
	"map":         "mapa",
	"chan":        "canal",
	"if":          "se",
	"else":        "senao",
	"for":         "para",
	"range":       "intervalo",
	"switch":      "escolha",
	"case":        "caso",
	"default":     "padrao",
	"fallthrough": "prosseguir",
	"break":       "interromper",
	"continue":    "continuar",
	"goto":        "ir_para",
	"return":      "retornar",
	"go":          "disparar",
	"defer":       "adiar",
	"select":      "selecionar",
}

// TypesAndBuiltinsMap mapeia tipos primitivos, literais e funções embutidas.
var TypesAndBuiltinsMap = map[string]string{
	// Tipos
	"texto":             "string",
	"inteiro":           "int",
	"inteiro8":          "int8",
	"inteiro16":         "int16",
	"inteiro32":         "int32",
	"inteiro64":         "int64",
	"uinteiro":          "uint",
	"uinteiro8":         "uint8",
	"uinteiro16":        "uint16",
	"uinteiro32":        "uint32",
	"uinteiro64":        "uint64",
	"uinteiro_ponteiro": "uintptr",
	"decimal":           "float64",
	"decimal32":         "float32",
	"decimal64":         "float64",
	"complexo64":        "complex64",
	"complexo128":       "complex128",
	"booleano":          "bool",
	"byte":              "byte",
	"runa":              "rune",
	"caractere":         "rune",
	"qualquer":          "any",
	"erro":              "error",

	// Literais
	"verdadeiro": "true",
	"falso":      "false",
	"nulo":       "nil",
	"iota":       "iota",

	// Built-ins
	"criar":      "make",
	"novo":       "new",
	"adicionar":  "append",
	"anexar":     "append",
	"tamanho":    "len",
	"capacidade": "cap",
	"fechar":     "close",
	"copiar":     "copy",
	"deletar":    "delete",
	"excluir":    "delete",
	"panico":     "panic",
	"recuperar":  "recover",
	"imprimir":   "print",
	"imprimirln": "println",
}

// ReverseTypesAndBuiltinsMap mapeia tipos e built-ins do Go para VAMOS.
var ReverseTypesAndBuiltinsMap = map[string]string{
	"string":     "texto",
	"int":        "inteiro",
	"int8":       "inteiro8",
	"int16":      "inteiro16",
	"int32":      "inteiro32",
	"int64":      "inteiro64",
	"uint":       "uinteiro",
	"uint8":      "uinteiro8",
	"uint16":     "uinteiro16",
	"uint32":     "uinteiro32",
	"uint64":     "uinteiro64",
	"uintptr":    "uinteiro_ponteiro",
	"float32":    "decimal32",
	"float64":    "decimal64",
	"complex64":  "complexo64",
	"complex128": "complexo128",
	"bool":       "booleano",
	"rune":       "runa",
	"any":        "qualquer",
	"error":      "erro",
	"true":       "verdadeiro",
	"false":      "falso",
	"nil":        "nulo",
	"make":       "criar",
	"new":        "novo",
	"append":     "adicionar",
	"len":        "tamanho",
	"cap":        "capacidade",
	"close":      "fechar",
	"copy":       "copiar",
	"delete":     "excluir",
	"panic":      "panico",
	"recover":    "recuperar",
}

// StdlibImportPathsMap mapeia nomes de pacotes PT-BR para import paths no Go.
var StdlibImportPathsMap = map[string]string{
	"formatar":      "fmt",
	"tempo":         "time",
	"matematica":    "math",
	"cordas":        "strings",
	"textos":        "strings",
	"so":            "os",
	"sistema":       "os",
	"sincronizar":   "sync",
	"erros":         "errors",
	"io":            "io",
	"leitor":        "bufio",
	"json":          "encoding/json",
	"servidor_http": "net/http",
	"http":          "net/http",
	"contexto":      "context",
	"registrador":   "log",
	"caminho":       "path/filepath",
	"ordenacao":     "sort",
	"teste":         "testing",
	"testing":       "testing",
}

// ReverseStdlibImportsMap mapeia import paths do Go para pacotes PT-BR.
var ReverseStdlibImportsMap = map[string]string{
	"fmt":           "formatar",
	"time":          "tempo",
	"math":          "matematica",
	"strings":       "cordas",
	"os":            "so",
	"sync":          "sincronizar",
	"errors":        "erros",
	"bufio":         "leitor",
	"encoding/json": "json",
	"net/http":      "servidor_http",
	"context":       "contexto",
	"log":           "registrador",
	"path/filepath": "caminho",
	"sort":          "ordenacao",
	"testing":       "teste",
}

// StdlibIdentMap mapeia identificadores de pacotes usados no código.
var StdlibIdentMap = map[string]string{
	"formatar":      "fmt",
	"fmt":           "fmt",
	"tempo":         "time",
	"time":          "time",
	"matematica":    "math",
	"math":          "math",
	"cordas":        "strings",
	"textos":        "strings",
	"strings":       "strings",
	"so":            "os",
	"sistema":       "os",
	"os":            "os",
	"sincronizar":   "sync",
	"sync":          "sync",
	"erros":         "errors",
	"errors":        "errors",
	"io":            "io",
	"leitor":        "bufio",
	"bufio":         "bufio",
	"json":          "json",
	"servidor_http": "http",
	"http":          "http",
	"contexto":      "context",
	"context":       "context",
	"registrador":   "log",
	"log":           "log",
	"caminho":       "filepath",
	"filepath":      "filepath",
	"ordenacao":     "sort",
	"sort":          "sort",
	"teste":         "testing",
	"testing":       "testing",
}

// StdlibMethodsMap mapeia métodos conhecidos de pacotes da biblioteca padrão.
var StdlibMethodsMap = map[string]map[string]string{
	"fmt": {
		"ImprimirLinha":      "Println",
		"ImprimirFormatado":  "Printf",
		"Imprimir":           "Print",
		"FormatarLinha":      "Sprintln",
		"Formatar":           "Sprintf",
		"FormatarSimples":    "Sprint",
		"CriarErro":          "Errorf",
		"ErroFormatado":      "Errorf",
		"Escanear":           "Scan",
		"EscanearLinha":      "Scanln",
		"EscanearFormatado":  "Scanf",
		"Fimprimir":          "Fprint",
		"FimprimirLinha":     "Fprintln",
		"FimprimirFormatado": "Fprintf",
	},
	"time": {
		"Dormir":              "Sleep",
		"Agora":               "Now",
		"Segundo":             "Second",
		"Milissegundo":        "Millisecond",
		"Microssegundo":       "Microsecond",
		"Nanosegundo":         "Nanosecond",
		"Minuto":              "Minute",
		"Hora":                "Hour",
		"Duracao":             "Duration",
		"Desde":               "Since",
		"Apos":                "After",
		"NovoTemporizador":    "NewTimer",
		"NovoTique":           "NewTicker",
		"Tique":               "Tick",
		"Data":                "Date",
		"CarregarLocalizacao": "LoadLocation",
		"FormatoUnix":         "Unix",
	},
	"math": {
		"Pi":        "Pi",
		"Raiz":      "Sqrt",
		"Potencia":  "Pow",
		"Absoluto":  "Abs",
		"Piso":      "Floor",
		"Teto":      "Ceil",
		"Maximo":    "Max",
		"Minimo":    "Min",
		"Seno":      "Sin",
		"Cosseno":   "Cos",
		"Tangente":  "Tan",
		"Logaritmo": "Log",
		"Mod":       "Mod",
	},
	"strings": {
		"Contem":         "Contains",
		"Dividir":        "Split",
		"Juntar":         "Join",
		"ParaMaiusculas": "ToUpper",
		"ParaMinusculas": "ToLower",
		"Substituir":     "Replace",
		"SubstituirTudo": "ReplaceAll",
		"TemPrefixo":     "HasPrefix",
		"TemSufixo":      "HasSuffix",
		"Aparar":         "Trim",
		"ApararEspaco":   "TrimSpace",
		"Indice":         "Index",
		"Repetir":        "Repeat",
		"Contar":         "Count",
		"Comparar":       "Compare",
		"IgualDobra":     "EqualFold",
	},
	"sync": {
		"GrupoEspera":         "WaitGroup",
		"Trava":               "Mutex",
		"TravaLeituraEscrita": "RWMutex",
		"UmaVez":              "Once",
		"Condicao":            "Cond",
		"NovoCondicao":        "NewCond",
	},
	"errors": {
		"Novo":         "New",
		"Desempacotar": "Unwrap",
		"E":            "Is",
		"Como":         "As",
		"Juntar":       "Join",
	},
	"os": {
		"Sair":                    "Exit",
		"Argumentos":              "Args",
		"LerArquivo":              "ReadFile",
		"EscreverArquivo":         "WriteFile",
		"ObterVariavelAmbiente":   "Getenv",
		"DefinirVariavelAmbiente": "Setenv",
		"CriarDiretorio":          "Mkdir",
		"CriarTodosDiretorios":    "MkdirAll",
		"Remover":                 "Remove",
		"RemoverTodos":            "RemoveAll",
		"Abrir":                   "Open",
		"Criar":                   "Create",
		"Arquivo":                 "File",
		"EntradaPadrao":           "Stdin",
		"SaidaPadrao":             "Stdout",
		"ErroPadrao":              "Stderr",
		"ObterDiretorioAtual":     "Getwd",
		"MudarDiretorio":          "Chdir",
		"NomeHost":                "Hostname",
	},
	"log": {
		"ImprimirLinha":     "Println",
		"ImprimirFormatado": "Printf",
		"Imprimir":          "Print",
		"Fatal":             "Fatal",
		"FatalFormatado":    "Fatalf",
		"FatalLinha":        "Fatalln",
		"Panico":            "Panic",
		"PanicoFormatado":   "Panicf",
		"PanicoLinha":       "Panicln",
	},
	"http": {
		"OuvirEServir":            "ListenAndServe",
		"Manipulador":             "HandleFunc",
		"NovoRequisicao":          "NewRequest",
		"Cliente":                 "Client",
		"Resposta":                "Response",
		"Requisicao":              "Request",
		"EscritorResposta":        "ResponseWriter",
		"StatusOk":                "StatusOK",
		"StatusCriado":            "StatusCreated",
		"StatusNaoEncontrado":     "StatusNotFound",
		"StatusErroInterno":       "StatusInternalServerError",
		"StatusRequisicaoInvalida": "StatusBadRequest",
		"Obter":                   "Get",
		"Postar":                  "Post",
	},
	"json": {
		"Serializar":         "Marshal",
		"SerializarIdentado": "MarshalIndent",
		"Deserializar":       "Unmarshal",
		"NovoCodificador":    "NewEncoder",
		"NovoDecodificador":  "NewDecoder",
	},
	"testing": {
		"Erro":           "Error",
		"ErroFormatado":  "Errorf",
		"Fatal":          "Fatal",
		"FatalFormatado": "Fatalf",
		"Falhar":         "Fail",
		"FalhaAgora":     "FailNow",
		"Log":            "Log",
		"LogFormatado":   "Logf",
		"Executar":       "Run",
		"Pular":          "Skip",
		"PularAgora":     "SkipNow",
	},
}

// MethodAliasesMap mapeia métodos chamados em instâncias de objetos e structs.
var MethodAliasesMap = map[string]string{
	"Adicionar":      "Add",
	"Concluido":      "Done",
	"Feito":          "Done",
	"Esperar":        "Wait",
	"Bloquear":       "Lock",
	"Desbloquear":    "Unlock",
	"Fechar":         "Close",
	"Ler":            "Read",
	"Escrever":       "Write",
	"Texto":          "String",
	"Erro":           "Error",
	"Formato":        "Format",
	"Formatar":       "Format",
	"Unix":           "Unix",
	"Segundos":       "Seconds",
	"Milissegundos":  "Milliseconds",
	"Microssegundos": "Microseconds",
	"Nanosegundos":   "Nanoseconds",
	"Minutos":        "Minutes",
	"Horas":          "Hours",
}

// ReverseMethodAliasesMap mapeia métodos Go para PT-BR.
var ReverseMethodAliasesMap = map[string]string{
	"Add":          "Adicionar",
	"Done":         "Concluido",
	"Wait":         "Esperar",
	"Lock":         "Bloquear",
	"Unlock":       "Desbloquear",
	"Close":        "Fechar",
	"Read":         "Ler",
	"Write":        "Escrever",
	"String":       "Texto",
	"Error":        "Erro",
	"Format":       "Formato",
	"Milliseconds": "Milissegundos",
	"Seconds":      "Segundos",
}

// ============================================================================
// 3. ANALISADOR LÉXICO (LEXER)
// ============================================================================

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenWhitespace
	TokenNewline
	TokenComment
	TokenIdentifier
	TokenKeyword
	TokenString
	TokenRawString
	TokenRune
	TokenNumber
	TokenOperator
	TokenDelimiter
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	source []rune
	pos    int
	line   int
	col    int
	length int
}

func NewLexer(input string) *Lexer {
	runes := []rune(input)
	return &Lexer{
		source: runes,
		pos:    0,
		line:   1,
		col:    1,
		length: len(runes),
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= l.length {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) peekNext() rune {
	if l.pos+1 >= l.length {
		return 0
	}
	return l.source[l.pos+1]
}

func (l *Lexer) next() rune {
	if l.pos >= l.length {
		return 0
	}
	r := l.source[l.pos]
	l.pos++
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isIdentPart(r rune) bool {
	return isLetter(r) || isDigit(r)
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens, nil
}

func (l *Lexer) NextToken() (Token, error) {
	if l.pos >= l.length {
		return Token{Type: TokenEOF, Value: "", Line: l.line, Column: l.col}, nil
	}

	r := l.peek()
	curLine := l.line
	curCol := l.col

	// 1. Espaços em branco (não quebra de linha)
	if r == ' ' || r == '\t' || r == '\r' {
		var val []rune
		for l.pos < l.length && (l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\r') {
			val = append(val, l.next())
		}
		return Token{Type: TokenWhitespace, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 2. Quebras de linha (preservação estrita 1:1)
	if r == '\n' {
		l.next()
		return Token{Type: TokenNewline, Value: "\n", Line: curLine, Column: curCol}, nil
	}

	// 3. Comentários de linha e de bloco
	if r == '/' && l.peekNext() == '/' {
		var val []rune
		for l.pos < l.length && l.peek() != '\n' {
			val = append(val, l.next())
		}
		return Token{Type: TokenComment, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	if r == '/' && l.peekNext() == '*' {
		var val []rune
		val = append(val, l.next()) // '/'
		val = append(val, l.next()) // '*'
		for l.pos < l.length {
			if l.peek() == '*' && l.peekNext() == '/' {
				val = append(val, l.next()) // '*'
				val = append(val, l.next()) // '/'
				break
			}
			val = append(val, l.next())
		}
		return Token{Type: TokenComment, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 4. Strings Literais Interpretadas ("...")
	if r == '"' {
		var val []rune
		val = append(val, l.next())
		escaped := false
		for l.pos < l.length {
			ch := l.next()
			val = append(val, ch)
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				break
			}
		}
		return Token{Type: TokenString, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 5. Strings Literais Brutas (`...`)
	if r == '`' {
		var val []rune
		val = append(val, l.next())
		for l.pos < l.length {
			ch := l.next()
			val = append(val, ch)
			if ch == '`' {
				break
			}
		}
		return Token{Type: TokenRawString, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 6. Literais de Caractere / Runa ('...')
	if r == '\'' {
		var val []rune
		val = append(val, l.next())
		escaped := false
		for l.pos < l.length {
			ch := l.next()
			val = append(val, ch)
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '\'' {
				break
			}
		}
		return Token{Type: TokenRune, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 7. Literais Numéricos
	if isDigit(r) || (r == '.' && isDigit(l.peekNext())) {
		var val []rune
		if r == '0' && (l.peekNext() == 'x' || l.peekNext() == 'X' || l.peekNext() == 'o' || l.peekNext() == 'O' || l.peekNext() == 'b' || l.peekNext() == 'B') {
			val = append(val, l.next())
			val = append(val, l.next())
			for l.pos < l.length && (isIdentPart(l.peek()) || l.peek() == '_') {
				val = append(val, l.next())
			}
			return Token{Type: TokenNumber, Value: string(val), Line: curLine, Column: curCol}, nil
		}

		hasDot := false
		hasExp := false
		for l.pos < l.length {
			ch := l.peek()
			if ch == '.' && !hasDot && isDigit(l.peekNext()) {
				hasDot = true
				val = append(val, l.next())
			} else if (ch == 'e' || ch == 'E') && !hasExp {
				hasExp = true
				val = append(val, l.next())
				if l.peek() == '+' || l.peek() == '-' {
					val = append(val, l.next())
				}
			} else if isDigit(ch) || ch == '_' {
				val = append(val, l.next())
			} else if ch == 'i' {
				val = append(val, l.next())
				break
			} else {
				break
			}
		}
		return Token{Type: TokenNumber, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 8. Identificadores e Palavras Reservadas
	if isLetter(r) {
		var val []rune
		for l.pos < l.length && isIdentPart(l.peek()) {
			val = append(val, l.next())
		}
		strVal := string(val)

		if _, ok := KeywordsMap[strVal]; ok {
			return Token{Type: TokenKeyword, Value: strVal, Line: curLine, Column: curCol}, nil
		}
		if _, ok := ReverseKeywordsMap[strVal]; ok {
			return Token{Type: TokenKeyword, Value: strVal, Line: curLine, Column: curCol}, nil
		}
		return Token{Type: TokenIdentifier, Value: strVal, Line: curLine, Column: curCol}, nil
	}

	// 9. Operadores e Delimitadores com 3 ou 2 caracteres
	if l.pos+2 < l.length {
		tri := string(l.source[l.pos : l.pos+3])
		if tri == "..." || tri == "<<=" || tri == ">>=" || tri == "&^=" {
			l.next()
			l.next()
			l.next()
			return Token{Type: TokenOperator, Value: tri, Line: curLine, Column: curCol}, nil
		}
	}

	if l.pos+1 < l.length {
		bi := string(l.source[l.pos : l.pos+2])
		switch bi {
		case ":=", "==", "!=", "<=", ">=", "&&", "||", "<-", "++", "--",
			"+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>", "&^":
			l.next()
			l.next()
			return Token{Type: TokenOperator, Value: bi, Line: curLine, Column: curCol}, nil
		}
	}

	// 10. Delimitadores e Operadores de 1 caractere
	single := l.next()
	strSingle := string(single)
	switch single {
	case '(', ')', '[', ']', '{', '}', ';', ',', ':', '.':
		return Token{Type: TokenDelimiter, Value: strSingle, Line: curLine, Column: curCol}, nil
	case '+', '-', '*', '/', '%', '&', '|', '^', '<', '>', '=', '!':
		return Token{Type: TokenOperator, Value: strSingle, Line: curLine, Column: curCol}, nil
	default:
		return Token{Type: TokenError, Value: strSingle, Line: curLine, Column: curCol},
			fmt.Errorf("caractere inesperado '%c' na linha %d, coluna %d", single, curLine, curCol)
	}
}

// ============================================================================
// 4. MOTOR DE TRANSPILAÇÃO (TRANSPILER PT-BR -> GO E REVERSO GO -> PT-BR)
// ============================================================================

type Transpiler struct {
	tokens []Token
	pos    int
	length int
}

func NewTranspiler(tokens []Token) *Transpiler {
	return &Transpiler{
		tokens: tokens,
		pos:    0,
		length: len(tokens),
	}
}

// TranspileSource transpila VAMOS-LANG para Go.
func TranspileSource(source string) (string, error) {
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return "", fmt.Errorf("erro léxico: %w", err)
	}

	t := NewTranspiler(tokens)
	return t.Transpile()
}

// TranspileGoToVamos converte código Go nativo de volta para VAMOS-LANG (PT-BR).
func TranspileGoToVamos(source string) (string, error) {
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return "", fmt.Errorf("erro léxico: %w", err)
	}

	var sb strings.Builder
	pos := 0
	length := len(tokens)

	for pos < length {
		tok := tokens[pos]

		if tok.Type == TokenEOF {
			break
		}

		if tok.Type == TokenWhitespace || tok.Type == TokenNewline || tok.Type == TokenComment || tok.Type == TokenRawString || tok.Type == TokenRune {
			sb.WriteString(tok.Value)
			pos++
			continue
		}

		// Importações Go -> VAMOS
		if tok.Type == TokenKeyword && tok.Value == "import" {
			pos++
			sb.WriteString("importar")
			// process reverse import
			inParen := false
			for pos < length {
				t := tokens[pos]
				if t.Type == TokenEOF {
					break
				}
				if t.Type == TokenWhitespace || t.Type == TokenComment {
					sb.WriteString(t.Value)
					pos++
					continue
				}
				if t.Type == TokenDelimiter && t.Value == "(" {
					inParen = true
					sb.WriteString(t.Value)
					pos++
					continue
				}
				if t.Type == TokenDelimiter && t.Value == ")" {
					inParen = false
					sb.WriteString(t.Value)
					pos++
					break
				}
				if t.Type == TokenNewline {
					sb.WriteString(t.Value)
					pos++
					if !inParen {
						break
					}
					continue
				}
				if t.Type == TokenString {
					raw := t.Value
					if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
						unq := raw[1 : len(raw)-1]
						if vamosPkg, ok := ReverseStdlibImportsMap[unq]; ok {
							sb.WriteString(fmt.Sprintf("%q", vamosPkg))
							pos++
							if !inParen {
								break
							}
							continue
						}
					}
					sb.WriteString(t.Value)
					pos++
					if !inParen {
						break
					}
					continue
				}
				sb.WriteString(t.Value)
				pos++
			}
			continue
		}

		// Package Go -> VAMOS
		if tok.Type == TokenKeyword && tok.Value == "package" {
			pos++
			sb.WriteString("pacote")
			continue
		}

		// Palavras reservadas Go -> VAMOS
		if tok.Type == TokenKeyword {
			if vamosKw, ok := ReverseKeywordsMap[tok.Value]; ok {
				sb.WriteString(vamosKw)
			} else {
				sb.WriteString(tok.Value)
			}
			pos++
			continue
		}

		if tok.Type == TokenString || tok.Type == TokenNumber {
			sb.WriteString(tok.Value)
			pos++
			continue
		}

		// Identificadores (membros como fmt.Println -> formatar.ImprimirLinha)
		if tok.Type == TokenIdentifier {
			if pos+2 < length && tokens[pos+1].Type == TokenDelimiter && tokens[pos+1].Value == "." && tokens[pos+2].Type == TokenIdentifier {
				pkgName := tok.Value
				methodName := tokens[pos+2].Value

				// Mapeia pkgName se for pacote stdlib
				mappedPkg := pkgName
				if ptPkg, ok := ReverseStdlibImportsMap[pkgName]; ok {
					mappedPkg = ptPkg
				}

				// Mapeia método
				mappedMethod := methodName
				if mmap, ok := StdlibMethodsMap[pkgName]; ok {
					for ptMethod, goMethod := range mmap {
						if goMethod == methodName {
							mappedMethod = ptMethod
							break
						}
					}
				} else if alias, ok := ReverseMethodAliasesMap[methodName]; ok {
					mappedMethod = alias
				}

				sb.WriteString(mappedPkg)
				sb.WriteString(".")
				sb.WriteString(mappedMethod)
				pos += 3
				continue
			}

			// Tipos / builtins Go -> VAMOS
			if vamosType, ok := ReverseTypesAndBuiltinsMap[tok.Value]; ok {
				sb.WriteString(vamosType)
			} else {
				sb.WriteString(tok.Value)
			}
			pos++
			continue
		}

		sb.WriteString(tok.Value)
		pos++
	}

	return sb.String(), nil
}

func (t *Transpiler) peek(offset int) *Token {
	idx := t.pos + offset
	if idx >= t.length {
		return &Token{Type: TokenEOF}
	}
	return &t.tokens[idx]
}

func (t *Transpiler) next() Token {
	if t.pos >= t.length {
		return Token{Type: TokenEOF}
	}
	tok := t.tokens[t.pos]
	t.pos++
	return tok
}

func (t *Transpiler) Transpile() (string, error) {
	var sb strings.Builder

	for t.pos < t.length {
		tok := t.peek(0)

		if tok.Type == TokenEOF {
			break
		}

		// Preservar whitespace, comentários e literais intactos
		if tok.Type == TokenWhitespace || tok.Type == TokenNewline || tok.Type == TokenComment || tok.Type == TokenRawString || tok.Type == TokenRune {
			sb.WriteString(t.next().Value)
			continue
		}

		// Importações
		if tok.Type == TokenKeyword && tok.Value == "importar" {
			t.next()
			sb.WriteString("import")
			t.processImportClause(&sb)
			continue
		}

		// Pacote
		if tok.Type == TokenKeyword && tok.Value == "pacote" {
			t.next()
			sb.WriteString("package")
			continue
		}

		// Palavras-chave regulares
		if tok.Type == TokenKeyword {
			val := t.next().Value
			if goKw, ok := KeywordsMap[val]; ok {
				sb.WriteString(goKw)
			} else {
				sb.WriteString(val)
			}
			continue
		}

		// Strings literais (fora de importações)
		if tok.Type == TokenString {
			sb.WriteString(t.next().Value)
			continue
		}

		// Números
		if tok.Type == TokenNumber {
			sb.WriteString(t.next().Value)
			continue
		}

		// Identificadores e acesso a membros (pacote.Metodo / objeto.Metodo)
		if tok.Type == TokenIdentifier {
			nextTok := t.peek(1)
			if nextTok.Type == TokenDelimiter && nextTok.Value == "." {
				nextNextTok := t.peek(2)
				if nextNextTok.Type == TokenIdentifier {
					pkgName := tok.Value
					methodName := nextNextTok.Value

					// Pacote padrão mapeado
					goIdent, isStdPkg := StdlibIdentMap[pkgName]
					if isStdPkg {
						if methodMap, hasMethods := StdlibMethodsMap[goIdent]; hasMethods {
							if goMethod, hasMethod := methodMap[methodName]; hasMethod {
								t.next() // consome pkgName
								t.next() // consome '.'
								t.next() // consome methodName
								sb.WriteString(goIdent)
								sb.WriteString(".")
								sb.WriteString(goMethod)
								continue
							}
						}
						if alias, hasAlias := MethodAliasesMap[methodName]; hasAlias {
							t.next()
							t.next()
							t.next()
							sb.WriteString(goIdent)
							sb.WriteString(".")
							sb.WriteString(alias)
							continue
						}

						t.next()
						t.next()
						sb.WriteString(goIdent)
						sb.WriteString(".")
						continue
					}

					// Pacote Go nativo
					if methodMap, hasMethods := StdlibMethodsMap[pkgName]; hasMethods {
						if goMethod, hasMethod := methodMap[methodName]; hasMethod {
							t.next() // consome pkgName
							t.next() // consome '.'
							t.next() // consome methodName
							sb.WriteString(pkgName)
							sb.WriteString(".")
							sb.WriteString(goMethod)
							continue
						}
					}

					// Métodos em objetos normais
					if alias, hasAlias := MethodAliasesMap[methodName]; hasAlias {
						identTok := t.next()
						t.next()
						t.next()
						mappedIdent := identTok.Value
						if mapped, ok := TypesAndBuiltinsMap[mappedIdent]; ok {
							mappedIdent = mapped
						}
						sb.WriteString(mappedIdent)
						sb.WriteString(".")
						sb.WriteString(alias)
						continue
					}
				}
			}

			// Tipos primitivos, built-ins ou identificadores normais
			val := t.next().Value
			if goIdent, isStdPkg := StdlibIdentMap[val]; isStdPkg {
				sb.WriteString(goIdent)
			} else if goType, ok := TypesAndBuiltinsMap[val]; ok {
				sb.WriteString(goType)
			} else {
				sb.WriteString(val)
			}
			continue
		}

		// Delimitador '.' seguido de método com alias
		if tok.Type == TokenDelimiter && tok.Value == "." {
			nextTok := t.peek(1)
			if nextTok.Type == TokenIdentifier {
				if alias, hasAlias := MethodAliasesMap[nextTok.Value]; hasAlias {
					t.next()
					t.next()
					sb.WriteString(".")
					sb.WriteString(alias)
					continue
				}
			}
		}

		// Operadores e delimitadores
		sb.WriteString(t.next().Value)
	}

	return sb.String(), nil
}

func (t *Transpiler) processImportClause(sb *strings.Builder) {
	inParen := false

	for t.pos < t.length {
		tok := t.peek(0)

		if tok.Type == TokenEOF {
			break
		}

		if tok.Type == TokenWhitespace || tok.Type == TokenComment {
			sb.WriteString(t.next().Value)
			continue
		}

		if tok.Type == TokenDelimiter && tok.Value == "(" {
			inParen = true
			sb.WriteString(t.next().Value)
			continue
		}

		if tok.Type == TokenDelimiter && tok.Value == ")" {
			inParen = false
			sb.WriteString(t.next().Value)
			break
		}

		if tok.Type == TokenNewline {
			sb.WriteString(t.next().Value)
			if !inParen {
				break
			}
			continue
		}

		if tok.Type == TokenDelimiter && tok.Value == ";" {
			sb.WriteString(t.next().Value)
			if !inParen {
				break
			}
			continue
		}

		if tok.Type == TokenString {
			rawPath := tok.Value
			if len(rawPath) >= 2 && rawPath[0] == '"' && rawPath[len(rawPath)-1] == '"' {
				unquoted := rawPath[1 : len(rawPath)-1]
				if mappedPkg, ok := StdlibImportPathsMap[unquoted]; ok {
					sb.WriteString(fmt.Sprintf("%q", mappedPkg))
					t.next()
					if !inParen {
						for t.pos < t.length && (t.peek(0).Type == TokenWhitespace || t.peek(0).Type == TokenComment) {
							sb.WriteString(t.next().Value)
						}
						if t.pos < t.length && t.peek(0).Type == TokenNewline {
							sb.WriteString(t.next().Value)
						}
						break
					}
					continue
				}
			}
			sb.WriteString(t.next().Value)
			if !inParen {
				break
			}
			continue
		}

		if tok.Type == TokenIdentifier {
			sb.WriteString(t.next().Value)
			continue
		}

		sb.WriteString(t.next().Value)
	}
}

// ============================================================================
// 5. FORMATADOR DE CÓDIGO (VAMOS FMT)
// ============================================================================

// FormatSource formata o código-fonte VAMOS-LANG com indentação consistente e limpeza.
func FormatSource(source string) (string, error) {
	lines := strings.Split(source, "\n")
	var formatted []string
	indent := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if len(formatted) > 0 && formatted[len(formatted)-1] != "" {
				formatted = append(formatted, "")
			}
			continue
		}

		// Ajusta indentação ao fechar bloco
		if strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, ")") || strings.HasPrefix(trimmed, "]") {
			if indent > 0 {
				indent--
			}
		}

		prefix := strings.Repeat("\t", indent)
		formatted = append(formatted, prefix+trimmed)

		// Incrementa indentação ao abrir bloco
		if strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, "(") || strings.HasSuffix(trimmed, "[") {
			indent++
		}
	}

	result := strings.Join(formatted, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}

// ============================================================================
// 6. LINTER E ANALISADOR ESTÁTICO (VAMOS LINT)
// ============================================================================

type LintIssue struct {
	Line    int
	Message string
}

func LintSource(source string) []LintIssue {
	var issues []LintIssue
	lines := strings.Split(source, "\n")

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Ignora comentários
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// Alerta 1: Descarte de erro sem verificação
		if strings.Contains(trimmed, "_ = err") || strings.Contains(trimmed, "_, _ =") || strings.Contains(trimmed, "_ = erro") {
			issues = append(issues, LintIssue{
				Line:    lineNum,
				Message: "Aviso: Tratamento de erro ignorado explicitamente com '_'",
			})
		}

		// Alerta 2: Uso incorreto de pacote principal em letras maiúsculas
		if strings.HasPrefix(trimmed, "pacote Principal") {
			issues = append(issues, LintIssue{
				Line:    lineNum,
				Message: "Convenção: 'pacote principal' deve ser todo em minúsculas",
			})
		}
	}

	return issues
}

// ============================================================================
// 7. AUXILIARES DE WORKSPACE, ARQUIVOS E PROJETOS
// ============================================================================

func prepareWorkspace(rootDir string) (string, func(), error) {
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
			if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" {
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

			goCode, err := TranspileSource(string(content))
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

		if d.Name() == "go.mod" || d.Name() == "go.sum" || strings.HasSuffix(path, ".go") {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			return copyFile(path, targetPath)
		}

		return nil
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

func copyFile(src, dst string) error {
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

func rewriteErrors(output, tempDir, rootDir string) string {
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

// ============================================================================
// 8. COMANDOS DA CLI
// ============================================================================

// cmdRun executa 'vamos run <arquivo.vamos> [args...]'
func cmdRun(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos informado para execução.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos run <arquivo.vamos> [argumentos...]\n")
		return 1
	}

	srcFile := args[0]
	progArgs := args[1:]

	absSrc, err := filepath.Abs(srcFile)
	if err != nil {
		absSrc = srcFile
	}

	content, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcFile, err)
		return 1
	}

	goCode, err := TranspileSource(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcFile, err)
		return 1
	}

	tempFile, err := os.CreateTemp("", "vamos_run_*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo temporário: %v\n", err)
		return 1
	}
	tempGoPath := tempFile.Name()
	defer os.Remove(tempGoPath)

	if _, err := tempFile.WriteString(goCode); err != nil {
		tempFile.Close()
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo temporário: %v\n", err)
		return 1
	}
	tempFile.Close()

	runArgs := append([]string{"run", tempGoPath}, progArgs...)
	cmd := exec.Command("go", runArgs...)
	cmd.Dir = filepath.Dir(absSrc)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderrPipe io.ReadCloser
	stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		cmd.Stderr = os.Stderr
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar o processo: %v\n", err)
		return 1
	}

	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()

	if stderrPipe != nil {
		var errBuf bytes.Buffer
		_, _ = io.Copy(&errBuf, stderrPipe)
		filtered := strings.ReplaceAll(errBuf.String(), tempGoPath, srcFile)
		if len(filtered) > 0 {
			fmt.Fprint(os.Stderr, filtered)
		}
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

// cmdBuild executa 'vamos build [flags...] [arquivo.vamos]'
func cmdBuild(args []string) int {
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

		content, err := os.ReadFile(srcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
			return 1
		}

		goCode, err := TranspileSource(string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcPath, err)
			return 1
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

		tempFile, err := os.CreateTemp("", "vamos_build_*.go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo temporário: %v\n", err)
			return 1
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
		cmd.Dir = filepath.Dir(absSrcPath)

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		execErr := cmd.Run()

		if stdoutBuf.Len() > 0 {
			fmt.Print(strings.ReplaceAll(stdoutBuf.String(), tempGoPath, srcPath))
		}
		if stderrBuf.Len() > 0 {
			fmt.Fprint(os.Stderr, strings.ReplaceAll(stderrBuf.String(), tempGoPath, srcPath))
		}

		if execErr != nil {
			fmt.Fprintf(os.Stderr, "Falha na compilação do binário.\n")
			return 1
		}

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

	tempDir, cleanup, err := prepareWorkspace(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao preparar workspace: %v\n", err)
		return 1
	}
	defer cleanup()

	buildArgs := []string{"build", "-o", absOutFile}
	buildArgs = append(buildArgs, forwardArgs...)
	buildArgs = append(buildArgs, ".")

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()

	if stdoutBuf.Len() > 0 {
		fmt.Print(rewriteErrors(stdoutBuf.String(), tempDir, cwd))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprint(os.Stderr, rewriteErrors(stderrBuf.String(), tempDir, cwd))
	}

	if execErr != nil {
		fmt.Fprintf(os.Stderr, "Falha na compilação.\n")
		return 1
	}

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
	return 0
}

// cmdTest executa 'vamos test [flags...]'
func cmdTest(args []string) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter diretório atual: %v\n", err)
		return 1
	}

	tempDir, cleanup, err := prepareWorkspace(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao preparar workspace: %v\n", err)
		return 1
	}
	defer cleanup()

	testArgs := append([]string{"test"}, args...)
	if len(args) == 0 {
		testArgs = append(testArgs, "./...")
	}

	cmd := exec.Command("go", testArgs...)
	cmd.Dir = tempDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()

	if stdoutBuf.Len() > 0 {
		fmt.Print(rewriteErrors(stdoutBuf.String(), tempDir, cwd))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprint(os.Stderr, rewriteErrors(stderrBuf.String(), tempDir, cwd))
	}

	if execErr != nil {
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

// cmdTranspile executa 'vamos transpile <origem.vamos> [-o <destino.go>]'
func cmdTranspile(args []string) int {
	var srcPath string
	var outFile string
	toStdout := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Erro: flag '%s' requer um caminho de saída.\n", arg)
				return 1
			}
		} else if strings.HasPrefix(arg, "-o=") {
			outFile = strings.TrimPrefix(arg, "-o=")
		} else if strings.HasPrefix(arg, "--output=") {
			outFile = strings.TrimPrefix(arg, "--output=")
		} else if arg == "--stdout" || arg == "-stdout" {
			toStdout = true
		} else if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Erro: flag desconhecida '%s'\n", arg)
			return 1
		} else if srcPath == "" {
			srcPath = arg
		} else {
			fmt.Fprintf(os.Stderr, "Erro: argumento inesperado '%s'\n", arg)
			return 1
		}
	}

	if srcPath == "" {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .vamos especificado.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos transpile <origem.vamos> [-o <destino.go>] [--stdout]\n")
		return 1
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
		return 1
	}

	goCode, err := TranspileSource(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro de transpilação em '%s':\n%v\n", srcPath, err)
		return 1
	}

	if toStdout {
		fmt.Print(goCode)
		return 0
	}

	target := outFile
	if target == "" {
		ext := filepath.Ext(srcPath)
		base := strings.TrimSuffix(srcPath, ext)
		target = base + ".go"
	}

	if err := os.WriteFile(target, []byte(goCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar arquivo gerado '%s': %v\n", target, err)
		return 1
	}

	fmt.Printf("✓ Arquivo transpilado com sucesso: %s -> %s\n", srcPath, target)
	return 0
}

// cmdGoToVamos executa 'vamos go2vamos <origem.go> [-o <destino.vamos>]'
func cmdGoToVamos(args []string) int {
	var srcPath string
	var outFile string
	toStdout := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" || arg == "--output" {
			if i+1 < len(args) {
				outFile = args[i+1]
				i++
			}
		} else if arg == "--stdout" || arg == "-stdout" {
			toStdout = true
		} else if srcPath == "" && !strings.HasPrefix(arg, "-") {
			srcPath = arg
		}
	}

	if srcPath == "" {
		fmt.Fprintf(os.Stderr, "Erro: nenhum arquivo .go especificado para conversão reversa.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos go2vamos <arquivo.go> [-o <destino.vamos>] [--stdout]\n")
		return 1
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler arquivo '%s': %v\n", srcPath, err)
		return 1
	}

	vamosCode, err := TranspileGoToVamos(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter Go para VAMOS: %v\n", err)
		return 1
	}

	if toStdout {
		fmt.Print(vamosCode)
		return 0
	}

	target := outFile
	if target == "" {
		ext := filepath.Ext(srcPath)
		base := strings.TrimSuffix(srcPath, ext)
		target = base + ".vamos"
	}

	if err := os.WriteFile(target, []byte(vamosCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar arquivo .vamos: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Código Go convertido para VAMOS-LANG com sucesso: %s -> %s\n", srcPath, target)
	return 0
}

// cmdFmt executa 'vamos fmt [caminho...]'
func cmdFmt(args []string) int {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	count := 0
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			formatted, err := FormatSource(string(content))
			if err != nil {
				return err
			}
			if string(content) != formatted {
				if err := os.WriteFile(path, []byte(formatted), 0644); err != nil {
					return err
				}
				fmt.Printf("✓ Formatado: %s\n", path)
				count++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao formatar arquivos: %v\n", err)
		return 1
	}

	if count == 0 {
		fmt.Println("✓ Todos os arquivos .vamos já estão devidamente formatados.")
	} else {
		fmt.Printf("✓ Total de %d arquivo(s) formatado(s) com sucesso.\n", count)
	}
	return 0
}

// cmdLint executa 'vamos lint [caminho...]'
func cmdLint(args []string) int {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	totalIssues := 0
	_ = filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".vamos") {
			content, err := os.ReadFile(path)
			if err == nil {
				issues := LintSource(string(content))
				for _, iss := range issues {
					fmt.Printf("%s:%d: %s\n", path, iss.Line, iss.Message)
					totalIssues++
				}
			}
		}
		return nil
	})

	if totalIssues == 0 {
		fmt.Println("✓ Análise concluída: nenhum aviso ou problema de estilo encontrado.")
		return 0
	}
	fmt.Printf("\nTotal de %d aviso(s) encontrado(s).\n", totalIssues)
	return 0
}

// cmdInit executa 'vamos init <nome-modulo>'
func cmdInit(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: informe o nome do módulo para o novo projeto.\n")
		fmt.Fprintf(os.Stderr, "Uso: vamos init <nome-do-modulo>\n")
		return 1
	}

	modName := args[0]
	projectDir := modName
	if projectDir == "." {
		cwd, _ := os.Getwd()
		projectDir = filepath.Base(cwd)
	}

	fmt.Printf("==> Criando novo projeto VAMOS-LANG: %s\n", modName)

	_ = os.MkdirAll(filepath.Join(projectDir, "cmd"), 0755)
	_ = os.MkdirAll(filepath.Join(projectDir, "pacotes", "matematica"), 0755)

	// go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.20\n", modName)
	_ = os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644)

	// cmd/principal.vamos
	mainContent := `// Ponto de entrada da aplicação
pacote principal

importar (
	"formatar"
	"tempo"
)

funcao principal() {
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirLinha("   Novo Projeto VAMOS-LANG Iniciado!    ")
	formatar.ImprimirLinha("========================================")
	formatar.ImprimirFormatado("Executado em: %s\n", tempo.Agora().Formato("02/01/2006 15:04:05"))
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "cmd", "principal.vamos"), []byte(mainContent), 0644)

	// pacotes/matematica/matematica.vamos
	mathContent := `pacote matematica

funcao Somar(a inteiro, b inteiro) inteiro {
	retornar a + b
}
`
	_ = os.WriteFile(filepath.Join(projectDir, "pacotes", "matematica", "matematica.vamos"), []byte(mathContent), 0644)

	// .gitignore
	gitIgnore := `# Binários
bin/
*.exe
*.out
app
`
	_ = os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitIgnore), 0644)

	fmt.Println("✓ Estrutura de pastas gerada com sucesso!")
	fmt.Printf("\nPara começar a desenvolver:\n")
	fmt.Printf("  cd %s\n", projectDir)
	fmt.Printf("  vamos run cmd/principal.vamos\n\n")
	return 0
}

// cmdREPL inicia o shell interativo do VAMOS-LANG
func cmdREPL() int {
	fmt.Println("==========================================================")
	fmt.Printf(" VAMOS-LANG REPL v%s (%s)\n", Version, Codename)
	fmt.Println(" Digite expressões ou declarações. Digite :ajuda para comandos.")
	fmt.Println("==========================================================")

	scanner := bufio.NewScanner(os.Stdin)
	var history []string

	for {
		fmt.Print("vamos> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if line == ":sair" || line == ":exit" || line == ":q" {
			fmt.Println("Até logo!")
			break
		}

		if line == ":ajuda" || line == ":help" {
			fmt.Println("Comandos do REPL:")
			fmt.Println("  :ajuda      Exibe esta ajuda")
			fmt.Println("  :codigo     Exibe o histórico de código acumulado na sessão")
			fmt.Println("  :limpar     Limpa a sessão atual do REPL")
			fmt.Println("  :sair       Encerra o REPL")
			continue
		}

		if line == ":limpar" || line == ":clear" {
			history = nil
			fmt.Println("✓ Sessão limpa com sucesso.")
			continue
		}

		if line == ":codigo" || line == ":source" {
			fmt.Println("// Código acumulado da sessão:")
			for _, h := range history {
				fmt.Println(h)
			}
			continue
		}

		// Monta o programa de teste para execução
		history = append(history, line)

		var progBuilder strings.Builder
		progBuilder.WriteString("pacote principal\n\nimportar (\n\t\"formatar\"\n\t\"tempo\"\n\t\"matematica\"\n)\n\nfuncao principal() {\n")
		for _, h := range history {
			// Se for apenas uma expressão de cálculo ou chamada, envolve em formatar.ImprimirLinha se aplicável
			if !strings.HasPrefix(h, "variavel ") && !strings.HasPrefix(h, "constante ") && !strings.HasPrefix(h, "tipo ") && !strings.HasPrefix(h, "funcao ") && !strings.Contains(h, ":=") && !strings.HasPrefix(h, "formatar.") && !strings.HasPrefix(h, "se ") && !strings.HasPrefix(h, "para ") {
				progBuilder.WriteString(fmt.Sprintf("\tformatar.ImprimirLinha(%s)\n", h))
			} else {
				progBuilder.WriteString(fmt.Sprintf("\t%s\n", h))
			}
		}
		progBuilder.WriteString("}\n")

		goCode, err := TranspileSource(progBuilder.String())
		if err != nil {
			fmt.Printf("Erro de sintaxe: %v\n", err)
			history = history[:len(history)-1] // desfaz a última linha
			continue
		}

		tempFile, err := os.CreateTemp("", "vamos_repl_*.go")
		if err != nil {
			fmt.Printf("Erro de I/O: %v\n", err)
			continue
		}
		tempGoPath := tempFile.Name()
		_ = os.WriteFile(tempGoPath, []byte(goCode), 0644)
		tempFile.Close()

		cmd := exec.Command("go", "run", tempGoPath)
		out, err := cmd.CombinedOutput()
		_ = os.Remove(tempGoPath)

		if err != nil {
			errMsg := string(out)
			lines := strings.Split(errMsg, "\n")
			if len(lines) > 0 {
				fmt.Printf("Erro: %s\n", lines[0])
			}
			history = history[:len(history)-1]
		} else {
			res := strings.TrimSpace(string(out))
			if res != "" {
				fmt.Println(res)
			}
		}
	}

	return 0
}

// cmdWeb inicia o Playground web
func cmdWeb(args []string) int {
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

// cmdVersion exibe informações de versão
func cmdVersion() {
	goVer := "go desconhecido"
	cmd := exec.Command("go", "version")
	if out, err := cmd.Output(); err == nil {
		goVer = strings.TrimSpace(string(out))
	}

	fmt.Printf("┌──────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│               VAMOS-LANG - Versão %-22s │\n", Version)
	fmt.Printf("│   Linguagem Go 100%% em Português do Brasil (PT-BR)       │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│ Codinome:     %-42s │\n", Codename)
	fmt.Printf("│ Plataforma:   %-42s │\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	fmt.Printf("│ Compilador:   %-42s │\n", goVer)
	fmt.Printf("│ Autor:        %-42s │\n", Author)
	fmt.Printf("└──────────────────────────────────────────────────────────┘\n")
}

// cmdHelp exibe o menu de ajuda
func cmdHelp() {
	fmt.Print(banner)
	fmt.Println("VAMOS-LANG: A linguagem de programação Go em Português do Brasil (PT-BR)")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  vamos <comando> [opções] [argumentos]")
	fmt.Println()
	fmt.Println("Comandos Principais (suporte a aliases em PT e EN):")
	fmt.Println("  run | rodar | executar <arquivo.vamos> [args...]   Transpila e executa diretamente")
	fmt.Println("  build | compilar | construir [flags...]            Transpila o projeto e compila binário")
	fmt.Println("  test | testar | teste [flags...]                   Transpila e executa os testes (go test)")
	fmt.Println("  transpile | converter | transpilar <origem> [-o]   Converte .vamos para Go (.go)")
	fmt.Println("  go2vamos | descompilar <arquivo.go> [-o]           Converte Go (.go) de volta para VAMOS")
	fmt.Println("  fmt | formatar [caminho...]                        Formata arquivos .vamos automaticamente")
	fmt.Println("  lint | checar | verificar [caminho...]             Verifica problemas de estilo e boas práticas")
	fmt.Println("  init | iniciar <nome-modulo>                       Cria um novo projeto com estrutura padrão")
	fmt.Println("  repl | interativo                                  Inicia o terminal interativo (REPL)")
	fmt.Println("  playground | web [porta]                           Inicia o Playground web interativo")
	fmt.Println("  version | versao                                   Exibe informações de versão")
	fmt.Println("  help | ajuda                                       Exibe esta tela de ajuda")
	fmt.Println()
	fmt.Println("Exemplos:")
	fmt.Println("  vamos run teste.vamos")
	fmt.Println("  vamos build -o meu_app")
	fmt.Println("  vamos fmt .")
	fmt.Println("  vamos repl")
	fmt.Println("  vamos init meu_novo_app")
	fmt.Println("  vamos go2vamos main.go -o main.vamos")
	fmt.Println("  vamos web 8080")
	fmt.Println()
}

// ============================================================================
// 9. FUNÇÃO PRINCIPAL (MAIN)
// ============================================================================

func main() {
	if len(os.Args) < 2 {
		cmdHelp()
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch command {
	case "run", "rodar", "executar":
		os.Exit(cmdRun(args))
	case "build", "compilar", "construir":
		os.Exit(cmdBuild(args))
	case "test", "testar", "teste":
		os.Exit(cmdTest(args))
	case "transpile", "converter", "transpilar":
		os.Exit(cmdTranspile(args))
	case "go2vamos", "descompilar":
		os.Exit(cmdGoToVamos(args))
	case "fmt", "formatar":
		os.Exit(cmdFmt(args))
	case "lint", "checar", "verificar":
		os.Exit(cmdLint(args))
	case "init", "iniciar":
		os.Exit(cmdInit(args))
	case "repl", "interativo":
		os.Exit(cmdREPL())
	case "playground", "web":
		os.Exit(cmdWeb(args))
	case "version", "versao", "-v", "--version":
		cmdVersion()
		os.Exit(0)
	case "help", "ajuda", "-h", "--help":
		cmdHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Erro: comando desconhecido '%s'\n\n", command)
		cmdHelp()
		os.Exit(1)
	}
}
