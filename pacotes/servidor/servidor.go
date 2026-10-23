package servidor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Contexto encapsula a requisição HTTP e fornece métodos facilitadores em Português.
type Contexto struct {
	Resposta   http.ResponseWriter
	Requisicao *http.Request
	Metodo     string
	Caminho    string
	params     map[string]string
}

// ManipuladorFunc é a assinatura padrão dos controladores de rota.
type ManipuladorFunc func(ctx *Contexto)

// MiddlewareFunc permite interceptar e tratar requisições antes do manipulador final.
type MiddlewareFunc func(proximo ManipuladorFunc) ManipuladorFunc

// Status define o código de resposta HTTP.
func (c *Contexto) Status(codigo int) {
	c.Resposta.WriteHeader(codigo)
}

// JSON serializa e envia dados no formato JSON com Content-Type apropriado.
func (c *Contexto) JSON(codigo int, dados interface{}) error {
	c.Resposta.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Resposta.WriteHeader(codigo)
	return json.NewEncoder(c.Resposta).Encode(dados)
}

// Texto envia uma resposta em texto puro.
func (c *Contexto) Texto(codigo int, mensagem string) error {
	c.Resposta.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Resposta.WriteHeader(codigo)
	_, err := c.Resposta.Write([]byte(mensagem))
	return err
}

// EscanearJSON decodifica o corpo da requisição JSON para uma estrutura.
func (c *Contexto) EscanearJSON(alvo interface{}) error {
	corpo, err := io.ReadAll(c.Requisicao.Body)
	if err != nil {
		return err
	}
	defer c.Requisicao.Body.Close()
	return json.Unmarshal(corpo, alvo)
}

// Parametro obtém um parâmetro da rota (ex: /usuarios/:id).
func (c *Contexto) Parametro(nome string) string {
	if c.params == nil {
		return ""
	}
	return c.params[nome]
}

// Consulta obtém um parâmetro da query string (?chave=valor).
func (c *Contexto) Consulta(chave string) string {
	return c.Requisicao.URL.Query().Get(chave)
}

// Servidor é a estrutura central da micro-framework web em VAMOS-LANG.
type Servidor struct {
	rotas       map[string]map[string]ManipuladorFunc // METODO -> ROTA -> HANDLER
	middlewares []MiddlewareFunc
}

// Novo cria e configura uma nova instância do Servidor web.
func Novo() *Servidor {
	return &Servidor{
		rotas: make(map[string]map[string]ManipuladorFunc),
	}
}

// Usar adiciona um middleware global que será executado em todas as rotas.
func (s *Servidor) Usar(m MiddlewareFunc) {
	s.middlewares = append(s.middlewares, m)
}

// RegistrarRota adiciona uma rota com suporte a verbos HTTP.
func (s *Servidor) RegistrarRota(metodo, caminho string, h ManipuladorFunc) {
	metodo = strings.ToUpper(metodo)
	if s.rotas[metodo] == nil {
		s.rotas[metodo] = make(map[string]ManipuladorFunc)
	}
	s.rotas[metodo][caminho] = h
}

// Obter registra uma rota HTTP GET.
func (s *Servidor) Obter(caminho string, h ManipuladorFunc) {
	s.RegistrarRota("GET", caminho, h)
}

// Postar registra uma rota HTTP POST.
func (s *Servidor) Postar(caminho string, h ManipuladorFunc) {
	s.RegistrarRota("POST", caminho, h)
}

// Colocar registra uma rota HTTP PUT.
func (s *Servidor) Colocar(caminho string, h ManipuladorFunc) {
	s.RegistrarRota("PUT", caminho, h)
}

// Deletar registra uma rota HTTP DELETE.
func (s *Servidor) Deletar(caminho string, h ManipuladorFunc) {
	s.RegistrarRota("DELETE", caminho, h)
}

// Ouvir inicia o servidor HTTP na porta e endereço especificados.
func (s *Servidor) Ouvir(endereco string) error {
	if !strings.HasPrefix(endereco, ":") && !strings.Contains(endereco, ":") {
		endereco = ":" + endereco
	}

	fmt.Printf("🚀 Servidor VAMOS-LANG ouvindo em http://localhost%s\n", endereco)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metodo := r.Method
		caminho := r.URL.Path

		rotasMetodo := s.rotas[metodo]
		var handlerFinal ManipuladorFunc
		params := make(map[string]string)

		// 1. Busca exata de rota
		if h, ok := rotasMetodo[caminho]; ok {
			handlerFinal = h
		} else {
			// 2. Busca por rotas parametrizadas (ex: /usuarios/:id)
			for rotaCadastrada, h := range rotasMetodo {
				if match, p := matchRoute(rotaCadastrada, caminho); match {
					handlerFinal = h
					params = p
					break
				}
			}
		}

		if handlerFinal == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"erro":     "Rota não encontrada",
				"caminho":  caminho,
				"mensagem": "Verifique o método HTTP e o endereço requisitado",
			})
			return
		}

		ctx := &Contexto{
			Resposta:   w,
			Requisicao: r,
			Metodo:     metodo,
			Caminho:    caminho,
			params:     params,
		}

		// Encadeia os middlewares
		exec := handlerFinal
		for i := len(s.middlewares) - 1; i >= 0; i-- {
			exec = s.middlewares[i](exec)
		}

		exec(ctx)
	})

	return http.ListenAndServe(endereco, handler)
}

func matchRoute(pattern, path string) (bool, map[string]string) {
	patParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i := 0; i < len(patParts); i++ {
		if strings.HasPrefix(patParts[i], ":") {
			paramName := strings.TrimPrefix(patParts[i], ":")
			params[paramName] = pathParts[i]
		} else if patParts[i] != pathParts[i] {
			return false, nil
		}
	}
	return true, params
}

// MiddlewareRegistrador cria um middleware que loga requisições com tempo de resposta.
func MiddlewareRegistrador() MiddlewareFunc {
	return func(proximo ManipuladorFunc) ManipuladorFunc {
		return func(ctx *Contexto) {
			inicio := time.Now()
			proximo(ctx)
			duracao := time.Since(inicio)
			fmt.Printf(" [LOG] %-7s %-25s | Duração: %v\n", ctx.Metodo, ctx.Caminho, duracao)
		}
	}
}

// MiddlewareCORS libera cabeçalhos para requisições cross-origin.
func MiddlewareCORS() MiddlewareFunc {
	return func(proximo ManipuladorFunc) ManipuladorFunc {
		return func(ctx *Contexto) {
			ctx.Resposta.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Resposta.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Resposta.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if ctx.Requisicao.Method == "OPTIONS" {
				ctx.Resposta.WriteHeader(http.StatusOK)
				return
			}
			proximo(ctx)
		}
	}
}
