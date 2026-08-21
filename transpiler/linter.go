package transpiler

import (
	"strings"
)

// LintIssue representa um aviso ou problema de estilo identificado no código VAMOS-LANG.
type LintIssue struct {
	Line    int
	Message string
}

// LintSource analisa estaticamente o código VAMOS-LANG em busca de más práticas e avisos de estilo.
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

		// Alerta 1: Descarte explícito de erro sem verificação
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
