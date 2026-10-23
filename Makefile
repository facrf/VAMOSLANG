.PHONY: all build test run-examples run-games check wasm clean install docker-build

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

clean:
	@echo "==> Limpando arquivos compilados..."
	rm -f $(BINARY_NAME)
	rm -f bin_*
	rm -f exemplos/*.go

install:
	@echo "==> Instalando 'vamos' no GOPATH/bin..."
	go install .
