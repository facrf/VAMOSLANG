package transpiler

import (
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"strings"
)

// DiagnosticSeverity indica a gravidade de um diagnóstico.
type DiagnosticSeverity string

const (
	SeverityError   DiagnosticSeverity = "error"
	SeverityWarning DiagnosticSeverity = "warning"
	SeverityInfo    DiagnosticSeverity = "info"
)

// SyntaxDiagnostic representa um erro de sintaxe ou aviso com sua localização no código VAMOS-LANG.
type SyntaxDiagnostic struct {
	Line     int                `json:"line"`
	Column   int                `json:"column"`
	Message  string             `json:"message"`
	Severity DiagnosticSeverity `json:"severity"`
}

// translateParserError traduz mensagens técnicas de erro sintático do parser Go para Português amigável.
func translateParserError(msg string) string {
	m := msg
	m = strings.ReplaceAll(m, "expected ';', found", "esperado ponto-e-vírgula ou quebra de linha, encontrado")
	m = strings.ReplaceAll(m, "expected '}', found 'EOF'", "bloco não encerrado: esperado fechamento de chave '}' antes do final do arquivo")
	m = strings.ReplaceAll(m, "expected ')', found", "esperado fechamento de parêntese ')', encontrado")
	m = strings.ReplaceAll(m, "expected ']', found", "esperado fechamento de colchete ']', encontrado")
	m = strings.ReplaceAll(m, "expected declaration, found", "esperada declaração (ex: funcao, variavel, constante), encontrado")
	m = strings.ReplaceAll(m, "expected operand, found", "esperado operando ou valor, encontrado")
	m = strings.ReplaceAll(m, "expected type, found", "esperado tipo de dado, encontrado")
	m = strings.ReplaceAll(m, "missing ',' in argument list", "vírgula ausente na lista de argumentos")
	m = strings.ReplaceAll(m, "missing ',' in parameter list", "vírgula ausente na lista de parâmetros")
	m = strings.ReplaceAll(m, "missing ',' before newline in composite literal", "vírgula obrigatória antes da quebra de linha em literal composto")
	m = strings.ReplaceAll(m, "syntax error:", "erro de sintaxe:")
	return m
}

// ValidateSource analisa o código VAMOS-LANG em múltiplas camadas (léxica, sintática Go e estilo linter).
func ValidateSource(source string) ([]SyntaxDiagnostic, error) {
	var diagnostics []SyntaxDiagnostic

	// 1. Camada Léxica
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		diagnostics = append(diagnostics, SyntaxDiagnostic{
			Line:     1,
			Column:   1,
			Message:  fmt.Sprintf("Erro léxico: %v", err),
			Severity: SeverityError,
		})
		return diagnostics, nil
	}

	// 2. Camada de Transpilação
	t := NewTranspiler(tokens)
	goCode, err := t.Transpile()
	if err != nil {
		diagnostics = append(diagnostics, SyntaxDiagnostic{
			Line:     1,
			Column:   1,
			Message:  fmt.Sprintf("Erro durante transpilação: %v", err),
			Severity: SeverityError,
		})
		return diagnostics, nil
	}

	// 3. Camada Sintática com Go Parser
	fset := token.NewFileSet()
	_, parseErr := parser.ParseFile(fset, "fonte.vamos", goCode, parser.AllErrors)
	if parseErr != nil {
		if errList, ok := parseErr.(scanner.ErrorList); ok {
			for _, item := range errList {
				diagnostics = append(diagnostics, SyntaxDiagnostic{
					Line:     item.Pos.Line,
					Column:   item.Pos.Column,
					Message:  translateParserError(item.Msg),
					Severity: SeverityError,
				})
			}
		} else {
			diagnostics = append(diagnostics, SyntaxDiagnostic{
				Line:     1,
				Column:   1,
				Message:  translateParserError(parseErr.Error()),
				Severity: SeverityError,
			})
		}
	}

	// 4. Camada de Linter (Boas Práticas e Convenções)
	lintIssues := LintSource(source)
	for _, issue := range lintIssues {
		diagnostics = append(diagnostics, SyntaxDiagnostic{
			Line:     issue.Line,
			Column:   1,
			Message:  issue.Message,
			Severity: SeverityWarning,
		})
	}

	return diagnostics, nil
}
