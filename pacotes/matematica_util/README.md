# 🧮 Pacote `matematica_util` (Biblioteca Padrão VAMOS-LANG)

O pacote **`matematica_util`** fornece algoritmos e utilitários matemáticos clássicos desenvolvidos integralmente na sintaxe VAMOS-LANG puro.

---

## 📐 Funções Disponíveis

| Função | Assinatura | Descrição |
| :--- | :--- | :--- |
| **`Absoluto`** | `Absoluto(n inteiro) inteiro` | Retorna o valor absoluto \|n\| do inteiro informado. |
| **`Maximo`** | `Maximo(a, b inteiro) inteiro` | Retorna o maior entre dois números inteiros. |
| **`Minimo`** | `Minimo(a, b inteiro) inteiro` | Retorna o menor entre dois números inteiros. |
| **`Potencia`** | `Potencia(base, expoente inteiro) inteiro` | Calcula base elevado ao expoente ($base^{expoente}$). |
| **`Fatorial`** | `Fatorial(n inteiro) inteiro` | Calcula o fatorial de $n!$ ($n \times (n-1) \dots 1$). |
| **`Fibonacci`** | `Fibonacci(n inteiro) inteiro` | Retorna o n-ésimo termo da sequência de Fibonacci de forma iterativa $O(n)$. |
| **`MDC`** | `MDC(a, b inteiro) inteiro` | Máximo Divisor Comum via Algoritmo de Euclides. |
| **`MMC`** | `MMC(a, b inteiro) inteiro` | Mínimo Múltiplo Comum ($|a \times b| / MDC(a,b)$). |
| **`EhPrimo`** | `EhPrimo(n inteiro) booleano` | Valida se o número é primo via teste de divisibilidade até $\sqrt{n}$. |
| **`Clampar`** | `Clampar(val, min, max inteiro) inteiro` | Restringe um valor ao intervalo delimitado $[min, max]$. |

---

## 🚀 Exemplo de Uso

```vamos
pacote principal

importar (
	"formatar"
	"matematica_util"
)

funcao principal() {
	formatar.ImprimirFormatado("Fatorial de 6: %d\n", matematica_util.Fatorial(6)) // 720
	formatar.ImprimirFormatado("Fibonacci(10): %d\n", matematica_util.Fibonacci(10)) // 55
	formatar.ImprimirFormatado("MDC(48, 18): %d\n", matematica_util.MDC(48, 18)) // 6
	formatar.ImprimirFormatado("MMC(4, 6): %d\n", matematica_util.MMC(4, 6)) // 12
	formatar.ImprimirFormatado("29 é primo? %v\n", matematica_util.EhPrimo(29)) // true
	formatar.ImprimirFormatado("Valor clampado: %d\n", matematica_util.Clampar(150, 0, 100)) // 100
}
```
