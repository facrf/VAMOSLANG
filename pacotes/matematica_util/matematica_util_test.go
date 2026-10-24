package matematica_util

import (
	"testing"
)

func TestMatematicaUtil(t *testing.T) {
	if Absoluto(-42) != 42 || Absoluto(42) != 42 {
		t.Errorf("Absoluto falhou")
	}

	if Maximo(10, 20) != 20 || Minimo(10, 20) != 10 {
		t.Errorf("Maximo/Minimo falhou")
	}

	if Potencia(2, 8) != 256 {
		t.Errorf("Potencia(2, 8) esperada 256, obteve %d", Potencia(2, 8))
	}

	if Fatorial(5) != 120 {
		t.Errorf("Fatorial(5) esperado 120, obteve %d", Fatorial(5))
	}

	if Fibonacci(10) != 55 {
		t.Errorf("Fibonacci(10) esperado 55, obteve %d", Fibonacci(10))
	}

	if MDC(48, 18) != 6 {
		t.Errorf("MDC(48, 18) esperado 6, obteve %d", MDC(48, 18))
	}

	if MMC(4, 6) != 12 {
		t.Errorf("MMC(4, 6) esperado 12, obteve %d", MMC(4, 6))
	}

	if !EhPrimo(17) || !EhPrimo(2) || EhPrimo(4) || EhPrimo(1) {
		t.Errorf("EhPrimo falhou")
	}

	if Clampar(150, 0, 100) != 100 || Clampar(-10, 0, 100) != 0 || Clampar(50, 0, 100) != 50 {
		t.Errorf("Clampar falhou")
	}
}
