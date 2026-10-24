// Pacote matematica_util - Utilitários numéricos e matemáticos em VAMOS-LANG
package matematica_util

// Absoluto retorna o valor absoluto de um inteiro
func Absoluto(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Maximo retorna o maior entre dois inteiros
func Maximo(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Minimo retorna o menor entre dois inteiros
func Minimo(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Potencia calcula base elevada ao expoente
func Potencia(base int, expoente int) int {
	if expoente == 0 {
		return 1
	}
	resultado := 1
	for i := 0; i < expoente; i++ {
		resultado = resultado * base
	}
	return resultado
}

// Fatorial calcula o fatorial de n
func Fatorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Fatorial(n - 1)
}

// Fibonacci calcula o n-ésimo termo da sequência de Fibonacci
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a := 0
	b := 1
	for i := 2; i <= n; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}

// MDC calcula o Máximo Divisor Comum pelo Algoritmo de Euclides
func MDC(a int, b int) int {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

// MMC calcula o Mínimo Múltiplo Comum
func MMC(a int, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return (Absoluto(a) * Absoluto(b)) / MDC(a, b)
}

// EhPrimo verifica se um número inteiro é primo
func EhPrimo(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n % 2 == 0 || n % 3 == 0 {
		return false
	}
	for i := 5; i * i <= n; i = i + 6 {
		if n % i == 0 || n % (i + 2) == 0 {
			return false
		}
	}
	return true
}

// Clampar restringe um valor entre um mínimo e um máximo
func Clampar(val int, min int, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
