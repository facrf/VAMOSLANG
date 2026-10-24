# 🔤 Pacote `texto_util` (Biblioteca Padrão VAMOS-LANG)

O pacote **`texto_util`** fornece operações canônicas e amigáveis para manipulação de textos e cadeias de caracteres em Português do Brasil com suporte total a **Unicode UTF-8**.

---

## 📝 Funções Disponíveis

| Função | Assinatura | Descrição |
| :--- | :--- | :--- |
| **`Dividir`** | `Dividir(s, separador texto) []texto` | Divide a string em fatias utilizando o separador informado. |
| **`Juntar`** | `Juntar(elementos []texto, separador texto) texto` | Concatena os elementos intercalados pelo separador. |
| **`Contem`** | `Contem(s, sub texto) booleano` | Verifica se a string engloba a substring indicada. |
| **`ComecaCom`** | `ComecaCom(s, prefixo texto) booleano` | Retorna se o texto inicia com o prefixo. |
| **`TerminaCom`** | `TerminaCom(s, sufixo texto) booleano` | Retorna se o texto finaliza com o sufixo. |
| **`ParaMaiusculas`**| `ParaMaiusculas(s texto) texto` | Converte o texto para caixa alta. |
| **`ParaMinusculas`**| `ParaMinusculas(s texto) texto` | Converte o texto para caixa baixa. |
| **`Aparar`** | `Aparar(s texto) texto` | Remove espaços em branco das extremidades (Trim). |
| **`Substituir`** | `Substituir(s, antigo, novo texto) texto` | Substitui todas as ocorrências de um trecho antigo por um novo. |
| **`Repetir`** | `Repetir(s texto, vezes inteiro) texto` | Gera uma string repetindo o conteúdo $N$ vezes. |
| **`IndiceDe`** | `IndiceDe(s, sub texto) inteiro` | Retorna o índice inicial da substring, ou -1 se inexistente. |
| **`Inverter`** | `Inverter(s texto) texto` | Inverte a ordem das runas/caracteres preservando acentuação UTF-8. |

---

## 🚀 Exemplo de Uso

```vamos
pacote principal

importar (
	"formatar"
	"texto_util"
)

funcao principal() {
	partes := texto_util.Dividir("VAMOS,LANG,Compilador", ",")
	formatar.ImprimirLinha(texto_util.Juntar(partes, " -> "))
	
	msg := "   Olá desenvolvedor VAMOS   "
	formatar.ImprimirLinha(texto_util.Aparar(msg))
	formatar.ImprimirLinha(texto_util.ParaMaiusculas("vamos-lang"))
	
	palavra := "café"
	formatar.ImprimirFormatado("'%s' invertido: '%s'\n", palavra, texto_util.Inverter(palavra)) // 'éfac'
}
```
