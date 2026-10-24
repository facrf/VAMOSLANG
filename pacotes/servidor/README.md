# 🌐 Pacote `servidor` (Micro-framework Web VAMOS-LANG)

O pacote **`servidor`** é um micro-framework HTTP nativo em Português do Brasil para desenvolvimento ágil de APIs REST e aplicações web em VAMOS-LANG.

---

## ⚡ Principais Recursos

- **Roteamento Expresso**: Métodos `GET`, `POST`, `PUT`, `DELETE` e `Grupo`.
- **Parâmetros de Rota**: Suporte a parâmetros dinâmicos `:param` (ex: `/api/tarefas/:id`).
- **Middlewares**: Encadeamento flexível (`Usar`) com middleware oficial de log estruturado (`MiddlewareRegistrador`).
- **Respostas Facilitadas**: `ctx.JSON(status, dados)`, `ctx.Texto(status, msg)`, `ctx.Status(codigo)`.
- **Decodificação Ágil**: `ctx.EscanearJSON(&objeto)` e `ctx.Parametro("chave")`.

---

## 🚀 Exemplo de API REST

```vamos
pacote principal

importar (
	"servidor"
)

tipo Tarefa estrutura {
	ID        inteiro `json:"id"`
	Descricao texto   `json:"descricao"`
	Feita     booleano `json:"feita"`
}

funcao principal() {
	app := servidor.Novo()
	app.Usar(servidor.MiddlewareRegistrador)

	app.GET("/", funcao(ctx *servidor.Contexto) {
		ctx.JSON(200, mapa[texto]texto{
			"mensagem": "API REST em VAMOS-LANG ativa!",
		})
	})

	app.GET("/api/tarefas/:id", funcao(ctx *servidor.Contexto) {
		idStr := ctx.Parametro("id")
		t := Tarefa{ID: 1, Descricao: "Aprender VAMOS-LANG", Feita: verdadeiro}
		ctx.JSON(200, t)
	})

	app.Ouvir(":8080")
}
```
