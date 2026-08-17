.PHONY: all build test run-examples clean install

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

clean:
	@echo "==> Limpando arquivos compilados..."
	rm -f $(BINARY_NAME)
	rm -f bin_*
	rm -f exemplos/*.go

install:
	@echo "==> Instalando 'vamos' no GOPATH/bin..."
	go install .
