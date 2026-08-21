package transpiler

import (
	"fmt"
	"unicode"
)

// Lexer realiza a análise léxica de código-fonte VAMOS-LANG.
type Lexer struct {
	source []rune
	pos    int
	line   int
	col    int
	length int
}

// NewLexer cria uma nova instância de Lexer para o código fornecido.
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

// peek retorna o caractere atual sem avançar o cursor.
func (l *Lexer) peek() rune {
	if l.pos >= l.length {
		return 0
	}
	return l.source[l.pos]
}

// peekNext retorna o próximo caractere sem avançar o cursor.
func (l *Lexer) peekNext() rune {
	if l.pos+1 >= l.length {
		return 0
	}
	return l.source[l.pos+1]
}

// next avança o cursor e retorna o caractere lido, atualizando linha e coluna.
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

// isLetter verifica se o rune é uma letra válida para identificadores (incluindo acentos PT-BR).
func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// isDigit verifica se o rune é um dígito decimal.
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// isIdentPart verifica se o rune pode compor o corpo de um identificador.
func isIdentPart(r rune) bool {
	return isLetter(r) || isDigit(r)
}

// Tokenize processa todo o código-fonte e retorna a lista completa de tokens.
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

// NextToken consome e retorna o próximo token léxico.
func (l *Lexer) NextToken() (Token, error) {
	if l.pos >= l.length {
		return Token{Type: TokenEOF, Value: "", Line: l.line, Column: l.col}, nil
	}

	r := l.peek()
	curLine := l.line
	curCol := l.col

	// 1. Tratamento de Whitespace (sem quebra de linha)
	if r == ' ' || r == '\t' || r == '\r' {
		var val []rune
		for l.pos < l.length && (l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\r') {
			val = append(val, l.next())
		}
		return Token{Type: TokenWhitespace, Value: string(val), Line: curLine, Column: curCol}, nil
	}

	// 2. Tratamento de Nova Linha (preservação estrita)
	if r == '\n' {
		l.next()
		return Token{Type: TokenNewline, Value: "\n", Line: curLine, Column: curCol}, nil
	}

	// 3. Comentários (linha '//' ou bloco '/* */')
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
		val = append(val, l.next()) // abre aspas
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
		val = append(val, l.next()) // abre crase
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
		val = append(val, l.next()) // abre aspas simples
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

	// 7. Literais Numéricos (Decimais, Hexadecimais, Flutuantes, etc.)
	if isDigit(r) || (r == '.' && isDigit(l.peekNext())) {
		var val []rune
		// Hex, Octal, Binário
		if r == '0' && (l.peekNext() == 'x' || l.peekNext() == 'X' || l.peekNext() == 'o' || l.peekNext() == 'O' || l.peekNext() == 'b' || l.peekNext() == 'B') {
			val = append(val, l.next()) // '0'
			val = append(val, l.next()) // base prefix
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
			} else if ch == 'i' { // Imaginário
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

		// Verifica se é uma palavra-chave reservada (VAMOS-LANG ou Go)
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
