package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"vamos/transpiler"
)

// Server representa o Language Server Protocol (LSP) para VAMOS-LANG.
type Server struct {
	reader io.Reader
	writer io.Writer
	mu     sync.Mutex
	docs   map[string]string // URI -> conteúdo do arquivo
}

// NewServer inicializa uma nova instância do LSP Server.
func NewServer(r io.Reader, w io.Writer) *Server {
	return &Server{
		reader: r,
		writer: w,
		docs:   make(map[string]string),
	}
}

// JSON-RPC tipos
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcNotification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LSP estruturas
type position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type lspRange struct {
	Start position `json:"start"`
	End   position `json:"end"`
}

type diagnostic struct {
	Range    lspRange `json:"range"`
	Severity int      `json:"severity"` // 1: Error, 2: Warning, 3: Info
	Source   string   `json:"source"`
	Message  string   `json:"message"`
}

type completionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"` // 14: Keyword, 7: Class/Struct, 3: Function, 6: Variable
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

type hover struct {
	Contents struct {
		Kind  string `json:"kind"` // "markdown" ou "plaintext"
		Value string `json:"value"`
	} `json:"contents"`
}

type textEdit struct {
	Range   lspRange `json:"range"`
	NewText string   `json:"newText"`
}

type lspLocation struct {
	URI   string   `json:"uri"`
	Range lspRange `json:"range"`
}

type documentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           int              `json:"kind"`
	Range          lspRange         `json:"range"`
	SelectionRange lspRange         `json:"selectionRange"`
	Children       []documentSymbol `json:"children,omitempty"`
}

// Run inicia o loop de processamento do protocolo LSP.
func (s *Server) Run() error {
	reader := bufio.NewReader(s.reader)

	for {
		// Lê cabeçalhos HTTP-like: "Content-Length: <n>\r\n\r\n"
		var contentLength int
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			if strings.HasPrefix(line, "Content-Length:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					cl, err := strconv.Atoi(strings.TrimSpace(parts[1]))
					if err == nil {
						contentLength = cl
					}
				}
			}
		}

		if contentLength == 0 {
			continue
		}

		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return err
		}

		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			continue
		}

		s.handleRequest(&req)
	}
}

func (s *Server) sendResponse(id json.RawMessage, result interface{}) {
	resp := rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.send(resp)
}

func (s *Server) sendNotification(method string, params interface{}) {
	notif := rpcNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	s.send(notif)
}

func (s *Server) send(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	_, _ = s.writer.Write([]byte(header))
	_, _ = s.writer.Write(data)
}

func (s *Server) setDoc(uri, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[uri] = text
}

func (s *Server) getDoc(uri string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.docs[uri]
}

func (s *Server) deleteDoc(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.docs, uri)
}

func (s *Server) handleRequest(req *rpcRequest) {
	switch req.Method {
	case "initialize":
		result := map[string]interface{}{
			"capabilities": map[string]interface{}{
				"textDocumentSync": 1, // Full sync
				"completionProvider": map[string]interface{}{
					"resolveProvider":   false,
					"triggerCharacters": []string{".", " "},
				},
				"hoverProvider":              true,
				"documentFormattingProvider": true,
				"definitionProvider":          true,
				"documentSymbolProvider":      true,
			},
			"serverInfo": map[string]interface{}{
				"name":    "vamos-lsp",
				"version": "1.0.0",
			},
		}
		s.sendResponse(req.ID, result)

	case "initialized":
		// Confirmação de inicialização do cliente

	case "shutdown":
		s.sendResponse(req.ID, nil)

	case "textDocument/didOpen":
		var params struct {
			TextDocument struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			s.setDoc(params.TextDocument.URI, params.TextDocument.Text)
			s.publishDiagnostics(params.TextDocument.URI, params.TextDocument.Text)
		}

	case "textDocument/didChange":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil && len(params.ContentChanges) > 0 {
			newText := params.ContentChanges[len(params.ContentChanges)-1].Text
			s.setDoc(params.TextDocument.URI, newText)
			s.publishDiagnostics(params.TextDocument.URI, newText)
		}

	case "textDocument/didSave":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			if doc != "" {
				s.publishDiagnostics(params.TextDocument.URI, doc)
			}
		}

	case "textDocument/didClose":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			s.deleteDoc(params.TextDocument.URI)
		}

	case "textDocument/completion":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			Position position `json:"position"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			items := s.getCompletions(doc, params.Position)
			s.sendResponse(req.ID, items)
		} else {
			s.sendResponse(req.ID, []completionItem{})
		}

	case "textDocument/hover":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			Position position `json:"position"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			h := s.getHover(doc, params.Position)
			s.sendResponse(req.ID, h)
		} else {
			s.sendResponse(req.ID, nil)
		}

	case "textDocument/formatting":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			formatted, err := transpiler.FormatSource(doc)
			if err == nil {
				lines := strings.Split(doc, "\n")
				lastLine := len(lines) - 1
				if lastLine < 0 {
					lastLine = 0
				}
				lastCol := 0
				if len(lines) > 0 {
					lastCol = len(lines[lastLine])
				}
				edits := []textEdit{
					{
						Range: lspRange{
							Start: position{Line: 0, Character: 0},
							End:   position{Line: lastLine, Character: lastCol},
						},
						NewText: formatted,
					},
				}
				s.sendResponse(req.ID, edits)
				return
			}
		}
		s.sendResponse(req.ID, []textEdit{})

	case "textDocument/definition":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			Position position `json:"position"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			loc := s.getDefinition(params.TextDocument.URI, doc, params.Position)
			s.sendResponse(req.ID, loc)
		} else {
			s.sendResponse(req.ID, nil)
		}

	case "textDocument/documentSymbol":
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err == nil {
			doc := s.getDoc(params.TextDocument.URI)
			syms := s.getDocumentSymbols(doc)
			s.sendResponse(req.ID, syms)
		} else {
			s.sendResponse(req.ID, []documentSymbol{})
		}

	default:
		if len(req.ID) > 0 {
			s.sendResponse(req.ID, nil)
		}
	}
}

func (s *Server) publishDiagnostics(uri, source string) {
	diagList, _ := transpiler.ValidateSource(source)
	var lspDiags []diagnostic

	for _, d := range diagList {
		severity := 1 // Error
		if d.Severity == transpiler.SeverityWarning {
			severity = 2 // Warning
		} else if d.Severity == transpiler.SeverityInfo {
			severity = 3
		}

		lineIdx := d.Line - 1
		if lineIdx < 0 {
			lineIdx = 0
		}
		colIdx := d.Column - 1
		if colIdx < 0 {
			colIdx = 0
		}

		lspDiags = append(lspDiags, diagnostic{
			Range: lspRange{
				Start: position{Line: lineIdx, Character: colIdx},
				End:   position{Line: lineIdx, Character: colIdx + 5},
			},
			Severity: severity,
			Source:   "vamos",
			Message:  d.Message,
		})
	}

	s.sendNotification("textDocument/publishDiagnostics", map[string]interface{}{
		"uri":         uri,
		"diagnostics": lspDiags,
	})
}

func (s *Server) getCompletions(doc string, pos position) []completionItem {
	var items []completionItem

	lines := strings.Split(doc, "\n")
	prefix := ""
	if pos.Line < len(lines) {
		line := lines[pos.Line]
		col := pos.Character
		if col > len(line) {
			col = len(line)
		}
		prefix = line[:col]
	}

	// 1. Se o usuário acabou de digitar "identificador."
	trimmedPrefix := strings.TrimSpace(prefix)
	if dotIdx := strings.LastIndex(trimmedPrefix, "."); dotIdx != -1 {
		ident := strings.TrimSpace(trimmedPrefix[:dotIdx])
		lastSpace := strings.LastIndexAny(ident, " \t(,+-*/%<>=!")
		if lastSpace != -1 {
			ident = ident[lastSpace+1:]
		}

		// Procura métodos conhecidos na Stdlib
		goIdent := ident
		if mapped, ok := transpiler.StdlibIdentMap[ident]; ok {
			goIdent = mapped
		}
		if methods, ok := transpiler.StdlibMethodsMap[goIdent]; ok {
			for ptMethod, goMethod := range methods {
				items = append(items, completionItem{
					Label:         ptMethod,
					Kind:          3, // Function
					Detail:        fmt.Sprintf("%s.%s (Go: %s.%s)", ident, ptMethod, goIdent, goMethod),
					Documentation: fmt.Sprintf("Método da biblioteca padrão VAMOS (%s.%s)", ident, ptMethod),
					InsertText:    ptMethod,
				})
			}
			return items
		}
	}

	// 2. Palavras-chave
	for vamosKw, goKw := range transpiler.KeywordsMap {
		items = append(items, completionItem{
			Label:         vamosKw,
			Kind:          14, // Keyword
			Detail:        fmt.Sprintf("Palavra-chave (Go: %s)", goKw),
			Documentation: fmt.Sprintf("Instrução reservada da linguagem VAMOS-LANG equivalente a '%s' no Go.", goKw),
			InsertText:    vamosKw,
		})
	}

	// 3. Tipos e Funções Embutidas (Builtins)
	for vamosType, goType := range transpiler.TypesAndBuiltinsMap {
		kind := 7 // Class / Type
		if strings.Contains("criar novo adicionar anexar tamanho capacidade fechar copiar excluir panico recuperar imprimir imprimirln", vamosType) {
			kind = 3 // Function
		}
		items = append(items, completionItem{
			Label:         vamosType,
			Kind:          kind,
			Detail:        fmt.Sprintf("Tipo / Builtin (Go: %s)", goType),
			Documentation: fmt.Sprintf("Identificador pré-definido equivalente a '%s' no Go.", goType),
			InsertText:    vamosType,
		})
	}

	// 4. Pacotes da Biblioteca Padrão
	for vamosPkg, goPkg := range transpiler.StdlibImportPathsMap {
		items = append(items, completionItem{
			Label:         vamosPkg,
			Kind:          9, // Module
			Detail:        fmt.Sprintf("Pacote Stdlib (Go: %s)", goPkg),
			Documentation: fmt.Sprintf("Pacote padrão VAMOS-LANG equivalente a '%s'.", goPkg),
			InsertText:    vamosPkg,
		})
	}

	return items
}

func (s *Server) getHover(doc string, pos position) *hover {
	lines := strings.Split(doc, "\n")
	if pos.Line >= len(lines) {
		return nil
	}

	line := lines[pos.Line]
	col := pos.Character
	if col > len(line) {
		col = len(line)
	}

	// Localiza os limites da palavra sob o cursor
	start := col
	for start > 0 && isWordChar(rune(line[start-1])) {
		start--
	}
	end := col
	for end < len(line) && isWordChar(rune(line[end])) {
		end++
	}

	if start >= end {
		return nil
	}

	word := line[start:end]

	var docStr string
	// Palavras-chave
	if goKw, ok := transpiler.KeywordsMap[word]; ok {
		docStr = fmt.Sprintf("### Palavra-chave: `%s`\n\n**Equivalente Go:** `%s`\n\nInstrução de controle ou declaração oficial da linguagem VAMOS-LANG.", word, goKw)
	} else if goType, ok := transpiler.TypesAndBuiltinsMap[word]; ok {
		docStr = fmt.Sprintf("### Tipo / Função Nativa: `%s`\n\n**Equivalente Go:** `%s`\n\nTipo primitivo ou função embutida (builtin) padrão.", word, goType)
	} else if goPkg, ok := transpiler.StdlibImportPathsMap[word]; ok {
		docStr = fmt.Sprintf("### Pacote Padrão: `%s`\n\n**Caminho Go:** `%s`\n\nBiblioteca padrão traduzida para Português do Brasil.", word, goPkg)
	} else if alias, ok := transpiler.MethodAliasesMap[word]; ok {
		docStr = fmt.Sprintf("### Método / Ação: `%s`\n\n**Equivalente Go:** `.%s()`\n\nMétodo comum de objeto ou interface mapeado.", word, alias)
	}

	if docStr == "" {
		return nil
	}

	h := &hover{}
	h.Contents.Kind = "markdown"
	h.Contents.Value = docStr
	return h
}

func isWordChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (s *Server) getDefinition(uri, doc string, pos position) *lspLocation {
	lines := strings.Split(doc, "\n")
	if pos.Line >= len(lines) {
		return nil
	}

	line := lines[pos.Line]
	col := pos.Character
	if col > len(line) {
		col = len(line)
	}

	start := col
	for start > 0 && isWordChar(rune(line[start-1])) {
		start--
	}
	end := col
	for end < len(line) && isWordChar(rune(line[end])) {
		end++
	}

	if start >= end {
		return nil
	}

	word := line[start:end]

	// Procura a declaração do símbolo no documento
	for lineIdx, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// 1. Função ou Método: funcao Nome(...) ou funcao (r Tipo) Nome(...)
		if strings.HasPrefix(trimmed, "funcao ") {
			afterFn := strings.TrimSpace(strings.TrimPrefix(trimmed, "funcao "))
			if strings.HasPrefix(afterFn, "(") {
				if idxClose := strings.Index(afterFn, ")"); idxClose != -1 {
					afterRec := strings.TrimSpace(afterFn[idxClose+1:])
					if idxParen := strings.Index(afterRec, "("); idxParen != -1 {
						fnName := strings.TrimSpace(afterRec[:idxParen])
						if fnName == word {
							charCol := strings.Index(l, fnName)
							return &lspLocation{
								URI: uri,
								Range: lspRange{
									Start: position{Line: lineIdx, Character: charCol},
									End:   position{Line: lineIdx, Character: charCol + len(fnName)},
								},
							}
						}
					}
				}
			} else {
				if idxParen := strings.Index(afterFn, "("); idxParen != -1 {
					fnName := strings.TrimSpace(afterFn[:idxParen])
					if idxBracket := strings.Index(fnName, "["); idxBracket != -1 {
						fnName = strings.TrimSpace(fnName[:idxBracket])
					}
					if fnName == word {
						charCol := strings.Index(l, fnName)
						return &lspLocation{
							URI: uri,
							Range: lspRange{
								Start: position{Line: lineIdx, Character: charCol},
								End:   position{Line: lineIdx, Character: charCol + len(fnName)},
							},
						}
					}
				}
			}
		}

		// 2. Tipo / Estrutura / Interface
		if strings.HasPrefix(trimmed, "tipo ") {
			afterTipo := strings.TrimSpace(strings.TrimPrefix(trimmed, "tipo "))
			parts := strings.Fields(afterTipo)
			if len(parts) > 0 {
				typeName := parts[0]
				if idxBracket := strings.Index(typeName, "["); idxBracket != -1 {
					typeName = typeName[:idxBracket]
				}
				if typeName == word {
					charCol := strings.Index(l, typeName)
					return &lspLocation{
						URI: uri,
						Range: lspRange{
							Start: position{Line: lineIdx, Character: charCol},
							End:   position{Line: lineIdx, Character: charCol + len(typeName)},
						},
					}
				}
			}
		}

		// 3. Constante ou Variável
		if strings.HasPrefix(trimmed, "constante ") || strings.HasPrefix(trimmed, "variavel ") {
			fields := strings.Fields(trimmed)
			if len(fields) >= 2 && fields[1] == word {
				charCol := strings.Index(l, word)
				return &lspLocation{
					URI: uri,
					Range: lspRange{
						Start: position{Line: lineIdx, Character: charCol},
						End:   position{Line: lineIdx, Character: charCol + len(word)},
					},
				}
			}
		}

		// 4. Declaração rápida: x := ...
		if strings.Contains(trimmed, ":=") {
			idxDecl := strings.Index(trimmed, ":=")
			left := strings.TrimSpace(trimmed[:idxDecl])
			vars := strings.Split(left, ",")
			for _, v := range vars {
				if strings.TrimSpace(v) == word {
					charCol := strings.Index(l, word)
					return &lspLocation{
						URI: uri,
						Range: lspRange{
							Start: position{Line: lineIdx, Character: charCol},
							End:   position{Line: lineIdx, Character: charCol + len(word)},
						},
					}
				}
			}
		}
	}

	return nil
}

func (s *Server) getDocumentSymbols(doc string) []documentSymbol {
	var symbols []documentSymbol
	lines := strings.Split(doc, "\n")

	for lineIdx, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// Função ou Método
		if strings.HasPrefix(trimmed, "funcao ") {
			afterFn := strings.TrimSpace(strings.TrimPrefix(trimmed, "funcao "))
			if strings.HasPrefix(afterFn, "(") {
				// Método com receptor: funcao (r Receptor) Nome(...)
				if idxClose := strings.Index(afterFn, ")"); idxClose != -1 {
					recStr := afterFn[1:idxClose]
					afterRec := strings.TrimSpace(afterFn[idxClose+1:])
					if idxParen := strings.Index(afterRec, "("); idxParen != -1 {
						fnName := strings.TrimSpace(afterRec[:idxParen])
						charCol := strings.Index(l, fnName)
						symbols = append(symbols, documentSymbol{
							Name:   fnName,
							Detail: fmt.Sprintf("(%s)", recStr),
							Kind:   6, // Method
							Range: lspRange{
								Start: position{Line: lineIdx, Character: 0},
								End:   position{Line: lineIdx, Character: len(l)},
							},
							SelectionRange: lspRange{
								Start: position{Line: lineIdx, Character: charCol},
								End:   position{Line: lineIdx, Character: charCol + len(fnName)},
							},
						})
					}
				}
			} else {
				// Função global
				if idxParen := strings.Index(afterFn, "("); idxParen != -1 {
					fnName := strings.TrimSpace(afterFn[:idxParen])
					detail := ""
					if idxBracket := strings.Index(fnName, "["); idxBracket != -1 {
						detail = fnName[idxBracket:]
						fnName = strings.TrimSpace(fnName[:idxBracket])
					}
					charCol := strings.Index(l, fnName)
					symbols = append(symbols, documentSymbol{
						Name:   fnName,
						Detail: detail,
						Kind:   12, // Function
						Range: lspRange{
							Start: position{Line: lineIdx, Character: 0},
							End:   position{Line: lineIdx, Character: len(l)},
						},
						SelectionRange: lspRange{
							Start: position{Line: lineIdx, Character: charCol},
							End:   position{Line: lineIdx, Character: charCol + len(fnName)},
						},
					})
				}
			}
			continue
		}

		// Tipos (Estruturas e Interfaces)
		if strings.HasPrefix(trimmed, "tipo ") {
			afterTipo := strings.TrimSpace(strings.TrimPrefix(trimmed, "tipo "))
			parts := strings.Fields(afterTipo)
			if len(parts) >= 2 {
				typeName := parts[0]
				kindStr := parts[1]
				kind := 23 // Struct
				if strings.HasPrefix(kindStr, "interface") {
					kind = 11 // Interface
				}
				charCol := strings.Index(l, typeName)
				symbols = append(symbols, documentSymbol{
					Name:   typeName,
					Detail: kindStr,
					Kind:   kind,
					Range: lspRange{
						Start: position{Line: lineIdx, Character: 0},
						End:   position{Line: lineIdx, Character: len(l)},
					},
					SelectionRange: lspRange{
						Start: position{Line: lineIdx, Character: charCol},
						End:   position{Line: lineIdx, Character: charCol + len(typeName)},
					},
				})
			}
			continue
		}

		// Constantes Globais
		if strings.HasPrefix(trimmed, "constante ") {
			afterConst := strings.TrimSpace(strings.TrimPrefix(trimmed, "constante "))
			parts := strings.Fields(afterConst)
			if len(parts) >= 1 {
				constName := parts[0]
				charCol := strings.Index(l, constName)
				symbols = append(symbols, documentSymbol{
					Name: constName,
					Kind: 14, // Constant
					Range: lspRange{
						Start: position{Line: lineIdx, Character: 0},
						End:   position{Line: lineIdx, Character: len(l)},
					},
					SelectionRange: lspRange{
						Start: position{Line: lineIdx, Character: charCol},
						End:   position{Line: lineIdx, Character: charCol + len(constName)},
					},
				})
			}
			continue
		}

		// Variáveis Globais
		if strings.HasPrefix(trimmed, "variavel ") {
			afterVar := strings.TrimSpace(strings.TrimPrefix(trimmed, "variavel "))
			parts := strings.Fields(afterVar)
			if len(parts) >= 1 {
				varName := parts[0]
				charCol := strings.Index(l, varName)
				symbols = append(symbols, documentSymbol{
					Name: varName,
					Kind: 13, // Variable
					Range: lspRange{
						Start: position{Line: lineIdx, Character: 0},
						End:   position{Line: lineIdx, Character: len(l)},
					},
					SelectionRange: lspRange{
						Start: position{Line: lineIdx, Character: charCol},
						End:   position{Line: lineIdx, Character: charCol + len(varName)},
					},
				})
			}
			continue
		}
	}

	return symbols
}

// ProcessSingleMessage processa uma única mensagem JSON-RPC síncrona (útil para testes unitários).
func (s *Server) ProcessSingleMessage(input []byte) ([]byte, error) {
	s.reader = bytes.NewReader(input)
	var out bytes.Buffer
	s.writer = &out

	reader := bufio.NewReader(s.reader)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	var contentLength int
	for {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				contentLength, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		}
		line, _ = reader.ReadString('\n')
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, err
	}

	var req rpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	s.handleRequest(&req)
	return out.Bytes(), nil
}
