package transpiler

import (
	"fmt"
	"strings"
)

// Transpiler é o motor responsável por converter código VAMOS-LANG em Go puro.
type Transpiler struct {
	tokens []Token
	pos    int
	length int
}

// NewTranspiler cria uma nova instância do transpilador a partir de uma lista de tokens.
func NewTranspiler(tokens []Token) *Transpiler {
	return &Transpiler{
		tokens: tokens,
		pos:    0,
		length: len(tokens),
	}
}

// TranspileSource recebe o código-fonte em VAMOS-LANG e retorna o código Go equivalente.
func TranspileSource(source string) (string, error) {
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return "", fmt.Errorf("erro léxico: %w", err)
	}

	t := NewTranspiler(tokens)
	return t.Transpile()
}

// peek retorna o token no deslocamento relativo offset (0 para atual, 1 para próximo, etc.)
func (t *Transpiler) peek(offset int) *Token {
	idx := t.pos + offset
	if idx >= t.length {
		return &Token{Type: TokenEOF}
	}
	return &t.tokens[idx]
}

// next consome e retorna o token atual.
func (t *Transpiler) next() Token {
	if t.pos >= t.length {
		return Token{Type: TokenEOF}
	}
	tok := t.tokens[t.pos]
	t.pos++
	return tok
}

// Transpile processa todos os tokens e gera o código Go correspondente.
func (t *Transpiler) Transpile() (string, error) {
	var sb strings.Builder

	for t.pos < t.length {
		tok := t.peek(0)

		if tok.Type == TokenEOF {
			break
		}

		// 1. Whitespace, Newline, Comentários, Raw Strings e Runes são preservados byte-a-byte
		if tok.Type == TokenWhitespace || tok.Type == TokenNewline || tok.Type == TokenComment || tok.Type == TokenRawString || tok.Type == TokenRune {
			sb.WriteString(t.next().Value)
			continue
		}

		// 2. Seção de Importação: mapear caminhos de pacotes padrão dentro de strings
		if tok.Type == TokenKeyword && tok.Value == "importar" {
			t.next()
			sb.WriteString("import")
			t.processImportClause(&sb)
			continue
		}

		// 3. Seção de Pacote: 'pacote principal' -> 'package main'
		if tok.Type == TokenKeyword && tok.Value == "pacote" {
			t.next()
			sb.WriteString("package")
			continue
		}

		// 4. Palavras-chave regulares
		if tok.Type == TokenKeyword {
			val := t.next().Value
			if goKw, ok := KeywordsMap[val]; ok {
				sb.WriteString(goKw)
			} else {
				sb.WriteString(val)
			}
			continue
		}

		// 5. Strings literais (fora de import) são preservadas
		if tok.Type == TokenString {
			sb.WriteString(t.next().Value)
			continue
		}

		// 6. Literais Numéricos
		if tok.Type == TokenNumber {
			sb.WriteString(t.next().Value)
			continue
		}

		// 7. Chamada qualificada de membro: pacote.Metodo ou objeto.Metodo
		if tok.Type == TokenIdentifier {
			nextTok := t.peek(1)
			if nextTok.Type == TokenDelimiter && nextTok.Value == "." {
				nextNextTok := t.peek(2)
				if nextNextTok.Type == TokenIdentifier {
					pkgName := tok.Value
					methodName := nextNextTok.Value

					// Identificador de pacote Go (ex: servidor_http -> http, formatar -> fmt)
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

						t.next() // consome pkgName
						t.next() // consome '.'
						sb.WriteString(goIdent)
						sb.WriteString(".")
						continue
					}

					// Se o pacote já está em Go nativo (ex: http.Manipulador)
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
						identTok := t.next() // consome objeto
						t.next()             // consome '.'
						t.next()             // consome metodo
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

			// 8. Tipos primitivos, built-ins, literais ou identificadores de pacotes isolados
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

		// 9. Delimitador '.' seguido de identificador com alias
		if tok.Type == TokenDelimiter && tok.Value == "." {
			nextTok := t.peek(1)
			if nextTok.Type == TokenIdentifier {
				if alias, hasAlias := MethodAliasesMap[nextTok.Value]; hasAlias {
					t.next() // consome '.'
					t.next() // consome ident
					sb.WriteString(".")
					sb.WriteString(alias)
					continue
				}
			}
		}

		// 10. Operadores e Delimitadores
		sb.WriteString(t.next().Value)
	}

	return sb.String(), nil
}

// processImportClause processa os tokens após uma declaração 'import' para traduzir caminhos de importação.
func (t *Transpiler) processImportClause(sb *strings.Builder) {
	inParen := false

	for t.pos < t.length {
		tok := t.peek(0)

		if tok.Type == TokenEOF {
			break
		}

		// Preservar espaços e comentários
		if tok.Type == TokenWhitespace || tok.Type == TokenComment {
			sb.WriteString(t.next().Value)
			continue
		}

		// Abertura de parênteses para bloco de imports
		if tok.Type == TokenDelimiter && tok.Value == "(" {
			inParen = true
			sb.WriteString(t.next().Value)
			continue
		}

		// Fechamento de parênteses
		if tok.Type == TokenDelimiter && tok.Value == ")" {
			inParen = false
			sb.WriteString(t.next().Value)
			break
		}

		// Nova linha: se não estiver em parênteses, encerra o import de linha única
		if tok.Type == TokenNewline {
			sb.WriteString(t.next().Value)
			if !inParen {
				break
			}
			continue
		}

		// Ponto e vírgula encerra import de linha única
		if tok.Type == TokenDelimiter && tok.Value == ";" {
			sb.WriteString(t.next().Value)
			if !inParen {
				break
			}
			continue
		}

		// Strings dentro do bloco de importação: traduzir se for pacote padrão
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

		// Se for identificador (ex: alias de import como `f "formatar"`)
		if tok.Type == TokenIdentifier {
			sb.WriteString(t.next().Value)
			continue
		}

		// Qualquer outro token
		sb.WriteString(t.next().Value)
	}
}
