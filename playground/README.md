# 🌐 VAMOS-LANG Web Playground

Ambiente interativo web para edição, teste e transpilação em tempo real de código **VAMOS-LANG** para **Go puro (Golang)**.

---

## ✨ Funcionalidades

1. **Transpilação em Tempo Real**: À medida que você digita código VAMOS-LANG no painel esquerdo, o painel direito exibe instantaneamente o código Go gerado.
2. **Exemplos Integrados**:
   - `1. Olá, Mundo!` (Variáveis, constantes e formatação)
   - `2. Concorrência` (`disparar`, `canal`, `adiar` e `sincronizar.GrupoEspera`)
   - `3. Estruturas e Interfaces` (`tipo`, `estrutura`, `interface` e polimorfismo)
   - `4. Calculadora e Erros` (`se`/`senao`, `escolha`/`caso`, slices e `erro`)
   - `5. Servidor Web HTTP` (`servidor_http` e rotas)
3. **Simulador de Execução**: Permite visualizar a saída esperada do código no terminal integrado.
4. **Inspetor Léxico (Tokens)**: Aba dedicada para visualizar cada token identificado pelo scanner (tipo, linha, coluna e valor léxico).
5. **Tabela de Sintaxe Interativa**: Modal de busca rápida de palavras-chave, tipos e métodos com recurso de clique para inserir no editor.
6. **Download e Cópia**: Exporte o código Go com um clique ou copie ambos os formatos para a área de transferência.

---

## 🚀 Como Executar

### Opção 1: Abrir diretamente no Navegador
Basta abrir o arquivo `index.html` em qualquer navegador web moderno:
```bash
# No Linux
xdg-open playground/index.html

# No macOS
open playground/index.html

# No Windows
start playground/index.html
```

### Opção 2: Servidor Local via Go ou Python
Se preferir servir por HTTP localmente:

```bash
# Usando o comando nativo do VAMOS-LANG CLI (se compilado):
./vamos web

# Ou via Python 3:
python3 -m http.server 8080 -d playground
```
Acesse em: `http://localhost:8080`

---

## 📂 Estrutura de Arquivos

- `index.html`: Interface visual estruturada e acessível.
- `style.css`: Folha de estilos moderna com design system em tons de ardósia (HSL), tipografia `Inter` e `JetBrains Mono`.
- `transpiler.js`: Motor léxico e transpilador escrito em JavaScript puro para execução no cliente.
- `app.js`: Controlador da interface, gerenciamento de snippets, simulação e área de transferência.
