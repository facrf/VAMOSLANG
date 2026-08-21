package transpiler

import (
	"strings"
)

// FormatSource formata o código-fonte VAMOS-LANG com indentação consistente e preservação de raw strings e comentários.
func FormatSource(source string) (string, error) {
	// Primeiro analisa léxicamente para identificar linhas dentro de strings brutas (backticks)
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		// Se houver erro léxico grave, faz formatação segura básica
		return formatBasic(source), nil
	}

	// Mapeia linhas que estão integralmente ou parcialmente dentro de raw strings com quebra de linha
	rawStringLines := make(map[int]bool)
	for _, tok := range tokens {
		if tok.Type == TokenRawString && strings.Contains(tok.Value, "\n") {
			startLine := tok.Line
			lineCount := strings.Count(tok.Value, "\n")
			for l := startLine + 1; l <= startLine+lineCount; l++ {
				rawStringLines[l] = true
			}
		}
	}

	lines := strings.Split(source, "\n")
	var formatted []string
	indent := 0

	for i, line := range lines {
		lineNum := i + 1

		// Linhas dentro de raw strings não devem ter espaços ou tabs alterados
		if rawStringLines[lineNum] {
			formatted = append(formatted, line)
			continue
		}

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

		// Incrementa indentação ao abrir bloco no final da linha (ignorando comentários)
		codeWithoutComment := trimmed
		if idx := strings.Index(trimmed, "//"); idx != -1 {
			codeWithoutComment = strings.TrimSpace(trimmed[:idx])
		}

		if strings.HasSuffix(codeWithoutComment, "{") || strings.HasSuffix(codeWithoutComment, "(") || strings.HasSuffix(codeWithoutComment, "[") {
			indent++
		}
	}

	result := strings.Join(formatted, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}

func formatBasic(source string) string {
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

		if strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, ")") || strings.HasPrefix(trimmed, "]") {
			if indent > 0 {
				indent--
			}
		}

		prefix := strings.Repeat("\t", indent)
		formatted = append(formatted, prefix+trimmed)

		if strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, "(") || strings.HasSuffix(trimmed, "[") {
			indent++
		}
	}

	result := strings.Join(formatted, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result
}
