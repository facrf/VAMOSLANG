.PHONY: all build test run-examples run-games check wasm clean install docker-build bootstrap

BINARY_NAME=vamos

all: build test

build:
	@echo "==> Compilando o transpilador VAMOS-LANG..."
	go build -o $(BINARY_NAME) .

test:
	@echo "==> Executando testes unitários..."
	go test -v ./...

run-examples: build
	@echo "==> Testando exemplo 'ola_mundo.vamos':"
	./$(BINARY_NAME) run exemplos/ola_mundo.vamos
	@echo ""
	@echo "==> Testando exemplo 'concorrencia.vamos':"
	./$(BINARY_NAME) run exemplos/concorrencia.vamos
	@echo ""
	@echo "==> Testando exemplo 'estruturas_interfaces.vamos':"
	./$(BINARY_NAME) run exemplos/estruturas_interfaces.vamos
	@echo ""
	@echo "==> Testando exemplo 'calculadora.vamos':"
	./$(BINARY_NAME) run exemplos/calculadora.vamos
	@echo ""
	@echo "==> Testando exemplo 'cripto_e_conversao.vamos':"
	./$(BINARY_NAME) run exemplos/cripto_e_conversao.vamos
	@echo ""
	@echo "==> Testando exemplo 'nano.vamos' (Editor de Texto):"
	./$(BINARY_NAME) run exemplos/nano.vamos --demo
	@echo ""
	@echo "==> Testando exemplo 'api_rest.vamos' (Micro-framework REST):"
	./$(BINARY_NAME) run exemplos/api_rest.vamos --demo
	@echo ""
	@echo "==> Testando exemplo 'uuid_externo.vamos' (Pacotes Go Externos):"
	./$(BINARY_NAME) run exemplos/uuid_externo.vamos
	@echo ""
	@echo "==> Testando exemplo 'genericos.vamos' (Tipos Genéricos e Restrições):"
	./$(BINARY_NAME) run exemplos/genericos.vamos
	@echo ""
	@echo "==> Testando exemplo 'maquina_virtual.vamos' (Execução Direta VVM):"
	./$(BINARY_NAME) run exemplos/maquina_virtual.vamos
	@echo ""
	@echo "==> Testando exemplo 'wasm_modulo.vamos' (Compilação Nativa WebAssembly):"
	./$(BINARY_NAME) wasm exemplos/wasm_modulo.vamos -o /tmp/modulo.wasm
	@rm -f /tmp/modulo.wasm
	@echo ""
	@echo "==> Testando exemplo 'c_nativo.vamos' (Backend C99 & GCC Nativo):"
	./$(BINARY_NAME) compilar-c exemplos/c_nativo.vamos -o /tmp/c_demo
	/tmp/c_demo
	@rm -f /tmp/c_demo
	@echo ""
	@echo "==> Testando exemplo 'pacotes_nativos.vamos' (Biblioteca Padrão em VAMOS puro):"
	./$(BINARY_NAME) run exemplos/pacotes_nativos.vamos

check: build
	@echo "==> 1. Verificando formatação de código..."
	./$(BINARY_NAME) fmt .
	@echo "==> 2. Executando linter estático e validação sintática..."
	./$(BINARY_NAME) lint --strict .
	@echo "==> 3. Rodando suíte completa de testes..."
	go test -v ./...
	@echo "==> 4. Gerando documentação oficial em HTML..."
	./$(BINARY_NAME) doc --html docs.html
	@echo "✓ Verificação completa de qualidade aprovada com êxito!"

docker-build:
	@echo "==> Construindo imagem Docker minimalista de produção..."
	docker build -t vamos-app .

wasm:
	@echo "==> Compilando motor WebAssembly (vamos.wasm)..."
	GOOS=js GOARCH=wasm go build -o playground/vamos.wasm ./wasm

run-games: build
	@echo "==> Executando demonstração do Jogo da Cobrinha:"
	./$(BINARY_NAME) run jogos/cobrinha.vamos --demo
	@echo ""
	@echo "==> Executando demonstração do Pacman:"
	./$(BINARY_NAME) run jogos/pacman.vamos --demo
	@echo ""
	@echo "==> Executando demonstração do Jogo de Damas:"
	./$(BINARY_NAME) run jogos/damas.vamos --demo

bootstrap: build
	@echo "==> 1. Compilando G1 'compilador/vamosc' via G0 (vamos)..."
	./$(BINARY_NAME) build -o compilador/vamosc compilador/main.vamos
	@echo "==> 2. Auto-teste de integridade da G1 em VAMOS-LANG..."
	./compilador/vamosc --auto-teste
	@echo "==> 3. Gerando AST com o parser nativo em VAMOS-LANG..."
	./compilador/vamosc ast exemplos/ola_mundo.vamos
	@echo "==> 4. Testando formatador de código em VAMOS-LANG..."
	./compilador/vamosc fmt exemplos/ola_mundo.vamos
	@echo "==> 5. Prova de VVM: Executando bytecode e Máquina Virtual nativa..."
	./compilador/vamosc bytecode exemplos/maquina_virtual.vamos -o compilador/mv.vamosb
	./compilador/vamosc vm compilador/mv.vamosb
	@rm -f compilador/mv.vamosb
	@echo "==> 6. Prova de Ouro: G1 compilando G2 (vamosc-g2)..."
	./compilador/vamosc compilar compilador/main.vamos -o compilador/vamosc-g2
	@echo "==> 7. Prova de Ouro: G2 compilando G3 (vamosc-g3)..."
	./compilador/vamosc-g2 compilar compilador/main.vamos -o compilador/vamosc-g3
	@echo "==> 8. Executando auto-teste com G3..."
	./compilador/vamosc-g3 --auto-teste
	@echo "==> 9. Prova de WebAssembly Nativo: Compilando binário Wasm 1.0 com G3..."
	./compilador/vamosc-g3 wasm exemplos/wasm_modulo.vamos -o compilador/modulo.wasm
	@echo "==> 10. Prova de Backend C Nativo: Compilando binário nativo via GCC com G3..."
	./compilador/vamosc-g3 compilar-c exemplos/c_nativo.vamos -o compilador/c_app
	./compilador/c_app
	@rm -f compilador/vamosc-g2 compilador/vamosc-g3 compilador/modulo.wasm compilador/c_app
	@echo "✓ Ciclo de Ouro de Bootstrapping (G0 -> G1 -> G2 -> G3 + Wasm + C/GCC) CONCLUÍDO COM ÊXITO!"

clean:
	@echo "==> Limpando arquivos compilados..."
	rm -f $(BINARY_NAME)
	rm -f compilador/vamosc*
	rm -f bin_*
	rm -f exemplos/*.go

install:
	@echo "==> Instalando 'vamos' no GOPATH/bin..."
	go install .
