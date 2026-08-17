package transpiler

// TokenType representa o tipo de token léxico identificado.
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

// Token representa uma unidade léxica com sua posição no código-fonte.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// KeywordsMap mapeia as palavras reservadas do VAMOS-LANG para Golang.
var KeywordsMap = map[string]string{
	"pacote":      "package",
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

// TypesAndBuiltinsMap mapeia tipos primitivos, valores literais e funções embutidas.
var TypesAndBuiltinsMap = map[string]string{
	// Tipos Primitivos
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

	// Constantes e Valores Literais
	"verdadeiro": "true",
	"falso":      "false",
	"nulo":       "nil",
	"iota":       "iota",

	// Funções Nativas Embutidas (Built-ins)
	"criar":      "make",
	"novo":       "new",
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

	// Ponto de entrada padrão
	"principal": "main",
}

// MethodAliasesMap mapeia nomes de métodos de objetos e interfaces em PT-BR para Go.
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

// StdlibImportPathsMap mapeia as strings de importação em PT-BR para caminhos Go no import ("...").
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
}

// StdlibIdentMap mapeia os identificadores de pacotes em código (ex: servidor_http. -> http., formatar. -> fmt.).
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
}

// StdlibMethodsMap mapeia métodos/funções conhecidos para o identificador do pacote Go.
var StdlibMethodsMap = map[string]map[string]string{
	"fmt": {
		"ImprimirLinha":      "Println",
		"ImprimirFormatado":  "Printf",
		"Imprimir":           "Print",
		"FormatarLinha":      "Sprintln",
		"Formatar":           "Sprintf",
		"FormatarSimples":    "Sprint",
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
}
