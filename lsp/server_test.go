package lsp

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func makeRPCRequest(id, method string, params interface{}) []byte {
	var pBytes []byte
	if params != nil {
		pBytes, _ = json.Marshal(params)
	}
	var reqStr string
	if id != "" {
		reqStr = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"method":"%s","params":%s}`, id, method, string(pBytes))
	} else {
		reqStr = fmt.Sprintf(`{"jsonrpc":"2.0","method":"%s","params":%s}`, method, string(pBytes))
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(reqStr))
	return []byte(header + reqStr)
}

func TestLSPInitialize(t *testing.T) {
	srv := NewServer(nil, nil)
	req := makeRPCRequest("1", "initialize", map[string]interface{}{})
	out, err := srv.ProcessSingleMessage(req)
	if err != nil {
		t.Fatalf("Erro ao processar initialize: %v", err)
	}

	outStr := string(out)
	if !strings.Contains(outStr, "vamos-lsp") {
		t.Errorf("Esperado 'vamos-lsp' na resposta do initialize, obteve: %s", outStr)
	}
	if !strings.Contains(outStr, "completionProvider") {
		t.Errorf("Esperado 'completionProvider' na resposta de capabilities")
	}
}

func TestLSPHoverAndCompletion(t *testing.T) {
	srv := NewServer(nil, nil)
	uri := "file:///teste.vamos"
	code := "pacote principal\nfuncao principal() {\n\tformatar.\n}"

	// 1. Simula didOpen
	didOpenReq := makeRPCRequest("", "textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":  uri,
			"text": code,
		},
	})
	_, _ = srv.ProcessSingleMessage(didOpenReq)

	// 2. Simula hover sobre 'funcao'
	hoverReq := makeRPCRequest("2", "textDocument/hover", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     position{Line: 1, Character: 2},
	})
	hOut, err := srv.ProcessSingleMessage(hoverReq)
	if err != nil {
		t.Fatalf("Erro no hover: %v", err)
	}
	if !strings.Contains(string(hOut), "Palavra-chave: `funcao`") {
		t.Errorf("Esperado hover de 'funcao', obteve: %s", string(hOut))
	}

	// 3. Simula completion após 'formatar.'
	compReq := makeRPCRequest("3", "textDocument/completion", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     position{Line: 2, Character: 10},
	})
	cOut, err := srv.ProcessSingleMessage(compReq)
	if err != nil {
		t.Fatalf("Erro no completion: %v", err)
	}
	if !strings.Contains(string(cOut), "ImprimirLinha") {
		t.Errorf("Esperado 'ImprimirLinha' no completion de 'formatar.', obteve: %s", string(cOut))
	}

	// 4. Simula completion de pacote nativo 'texto_util.'
	codeTexto := "pacote principal\nfuncao principal() {\n\ttexto_util.\n}"
	srv.ProcessSingleMessage(makeRPCRequest("", "textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri, "text": codeTexto},
	}))
	compTextoReq := makeRPCRequest("5", "textDocument/completion", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     position{Line: 2, Character: 12},
	})
	cTextoOut, errTexto := srv.ProcessSingleMessage(compTextoReq)
	if errTexto != nil || !strings.Contains(string(cTextoOut), "Dividir") || !strings.Contains(string(cTextoOut), "Inverter") {
		t.Errorf("Esperado 'Dividir' e 'Inverter' no completion de 'texto_util.', obteve: %s", string(cTextoOut))
	}
}

func TestLSPFormatting(t *testing.T) {
	srv := NewServer(nil, nil)
	uri := "file:///teste.vamos"
	unformatted := "pacote principal\nfuncao principal(){\nformatar.ImprimirLinha(\"oi\")\n}"

	didOpenReq := makeRPCRequest("", "textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":  uri,
			"text": unformatted,
		},
	})
	_, _ = srv.ProcessSingleMessage(didOpenReq)

	fmtReq := makeRPCRequest("4", "textDocument/formatting", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
	})
	fOut, err := srv.ProcessSingleMessage(fmtReq)
	if err != nil {
		t.Fatalf("Erro no formatting: %v", err)
	}
	if !strings.Contains(string(fOut), "\\tformatar.ImprimirLinha") {
		t.Errorf("Esperado código formatado com indentação tab, obteve: %s", string(fOut))
	}
}

func TestLSPDefinition(t *testing.T) {
	srv := NewServer(nil, nil)
	uri := "file:///definicao.vamos"
	code := `pacote principal

funcao somar(a inteiro, b inteiro) inteiro {
	retornar a + b
}

tipo Pessoa estrutura {
	Nome texto
}

funcao principal() {
	resultado := somar(10, 20)
	p := Pessoa{Nome: "Maria"}
}`

	srv.ProcessSingleMessage(makeRPCRequest("", "textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri, "text": code},
	}))

	// Procura definição de 'somar' na linha da chamada (linha 11)
	defReq := makeRPCRequest("10", "textDocument/definition", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     position{Line: 11, Character: 15},
	})
	dOut, err := srv.ProcessSingleMessage(defReq)
	if err != nil {
		t.Fatalf("Erro no definition: %v", err)
	}
	outStr := string(dOut)
	if !strings.Contains(outStr, `"line":2`) {
		t.Errorf("Esperado definição de 'somar' na linha 2, obteve: %s", outStr)
	}

	// Procura definição de 'Pessoa' na linha de instanciação (linha 12)
	defTipoReq := makeRPCRequest("11", "textDocument/definition", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     position{Line: 12, Character: 7},
	})
	dTipoOut, err := srv.ProcessSingleMessage(defTipoReq)
	if err != nil {
		t.Fatalf("Erro no definition de tipo: %v", err)
	}
	if !strings.Contains(string(dTipoOut), `"line":6`) {
		t.Errorf("Esperado definição de 'Pessoa' na linha 6, obteve: %s", string(dTipoOut))
	}
}

func TestLSPDocumentSymbol(t *testing.T) {
	srv := NewServer(nil, nil)
	uri := "file:///simbolos.vamos"
	code := `pacote principal

constante VersaoApp texto = "1.0"
variavel Contador inteiro = 0

tipo Servico interface {
	Executar()
}

tipo Configuracao estrutura {
	Porta inteiro
}

funcao Iniciar() {
}

funcao (c Configuracao) ObterPorta() inteiro {
	retornar c.Porta
}`

	srv.ProcessSingleMessage(makeRPCRequest("", "textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri, "text": code},
	}))

	symReq := makeRPCRequest("20", "textDocument/documentSymbol", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
	})
	sOut, err := srv.ProcessSingleMessage(symReq)
	if err != nil {
		t.Fatalf("Erro no documentSymbol: %v", err)
	}
	outStr := string(sOut)
	expectedSymbols := []string{"VersaoApp", "Contador", "Servico", "Configuracao", "Iniciar", "ObterPorta"}
	for _, sym := range expectedSymbols {
		if !strings.Contains(outStr, sym) {
			t.Errorf("Esperado símbolo '%s' na lista do documento, obteve: %s", sym, outStr)
		}
	}
}
