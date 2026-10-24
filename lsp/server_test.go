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
