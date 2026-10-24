// Pacote texto_util - Manipulação nativa de cadeias de caracteres em VAMOS-LANG
// Oferece funções completas em Português para operações com textos/strings.

package texto_util

import (
	"strings"
)

// Dividir separa uma string em partes utilizando um delimitador especificado
func Dividir(s string, separador string) []string {
	return strings.Split(s, separador)
}

// Juntar concatena uma lista de strings intercalando pelo separador indicado
func Juntar(elementos []string, separador string) string {
	return strings.Join(elementos, separador)
}

// Contem verifica se o texto principal engloba a substring informada
func Contem(s string, sub string) bool {
	return strings.Contains(s, sub)
}

// ComecaCom valida se o texto inicia com o prefixo informado
func ComecaCom(s string, prefixo string) bool {
	return strings.HasPrefix(s, prefixo)
}

// TerminaCom valida se o texto termina com o sufixo informado
func TerminaCom(s string, sufixo string) bool {
	return strings.HasSuffix(s, sufixo)
}

// ParaMaiusculas converte todos os caracteres para caixa alta
func ParaMaiusculas(s string) string {
	return strings.ToUpper(s)
}

// ParaMinusculas converte todos os caracteres para caixa baixa
func ParaMinusculas(s string) string {
	return strings.ToLower(s)
}

// Aparar remove espaços em branco do início e do final da cadeia de caracteres
func Aparar(s string) string {
	return strings.TrimSpace(s)
}

// Substituir troca ocorrências de um trecho antigo por um novo
func Substituir(s string, antigo string, new string) string {
	return strings.ReplaceAll(s, antigo, new)
}

// Repetir produz uma nova string repetindo a original N vezes
func Repetir(s string, vezes int) string {
	if vezes <= 0 {
		return ""
	}
	return strings.Repeat(s, vezes)
}

// IndiceDe retorna a posição inicial da primeira ocorrência da substring, ou -1 se não existir
func IndiceDe(s string, sub string) int {
	return strings.Index(s, sub)
}

// Inverter inverte a ordem dos caracteres do texto informado
func Inverter(s string) string {
	runas := []rune(s)
	n := len(runas)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		temporario := runas[i]
		runas[i] = runas[j]
		runas[j] = temporario
	}
	return string(runas)
}
