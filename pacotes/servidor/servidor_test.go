package servidor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServidorRotasERespontas(t *testing.T) {
	app := Novo()

	app.Obter("/api/ola", func(ctx *Contexto) {
		_ = ctx.JSON(200, map[string]string{"mensagem": "Olá mundo"})
	})

	app.Obter("/api/usuarios/:id", func(ctx *Contexto) {
		id := ctx.Parametro("id")
		_ = ctx.JSON(200, map[string]string{"id": id})
	})

	app.Postar("/api/eco", func(ctx *Contexto) {
		var corpo map[string]string
		_ = ctx.EscanearJSON(&corpo)
		_ = ctx.JSON(201, corpo)
	})

	// Teste 1: Rota simples
	req1 := httptest.NewRequest("GET", "/api/ola", nil)
	w1 := httptest.NewRecorder()
	handler := createTestHandler(app)
	handler.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("Esperado status 200, obteve %d", w1.Code)
	}
	if !strings.Contains(w1.Body.String(), "Olá mundo") {
		t.Errorf("Resposta inesperada: %s", w1.Body.String())
	}

	// Teste 2: Parâmetro dinâmico
	req2 := httptest.NewRequest("GET", "/api/usuarios/42", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("Esperado status 200, obteve %d", w2.Code)
	}
	var res2 map[string]string
	_ = json.Unmarshal(w2.Body.Bytes(), &res2)
	if res2["id"] != "42" {
		t.Errorf("Esperado ID '42', obteve %s", res2["id"])
	}

	// Teste 3: Rota não encontrada
	req3 := httptest.NewRequest("GET", "/rota/inexistente", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code != 404 {
		t.Errorf("Esperado status 404 para rota inexistente, obteve %d", w3.Code)
	}
}

func createTestHandler(s *Servidor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metodo := r.Method
		caminho := r.URL.Path

		rotasMetodo := s.rotas[metodo]
		var handlerFinal ManipuladorFunc
		params := make(map[string]string)

		if h, ok := rotasMetodo[caminho]; ok {
			handlerFinal = h
		} else {
			for rotaCadastrada, h := range rotasMetodo {
				if match, p := matchRoute(rotaCadastrada, caminho); match {
					handlerFinal = h
					params = p
					break
				}
			}
		}

		if handlerFinal == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := &Contexto{
			Resposta:   w,
			Requisicao: r,
			Metodo:     metodo,
			Caminho:    caminho,
			params:     params,
		}

		handlerFinal(ctx)
	})
}
