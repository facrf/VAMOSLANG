# 🇧🇷 VAMOS-LANG - Extensão para VS Code

Extensão oficial para o **Visual Studio Code** e **Cursor**, adicionando suporte completo de sintaxe e snippets para a linguagem de programação **VAMOS-LANG** (Golang 100% em Português).

---

## ✨ Recursos

- 🎨 **Realce de Sintaxe Oficial**: Cores precisas para palavras-chave (`pacote`, `funcao`, `se`, `disparar`, `canal`, `adiar`), tipos (`texto`, `inteiro`, `decimal64`, `booleano`), literais e métodos.
- ⚡ **Snippets Inteligentes**:
  - `pkgm` ➔ Estrutura inicial de programa com `pacote principal` e `funcao principal()`.
  - `funcao` ➔ Declaração de função.
  - `se` / `sesenao` ➔ Condicionais.
  - `para` / `intervalo` ➔ Laços e iteração em coleções.
  - `estrutura` ➔ Definição de structs.
  - `interface` ➔ Definição de interfaces.
  - `disparar` ➔ Goroutine anônima.
  - `imp` / `impf` ➔ `formatar.ImprimirLinha` e `formatar.ImprimirFormatado`.
- 🗂️ **Configuração de Linguagem**: Fechamento automático de chaves, parênteses, aspas e indentação automática em blocos de código.

---

## 📦 Como Instalar Localmente

### Opção 1: Copiar para a pasta de extensões do VS Code
```bash
# No Linux / macOS
cp -r /storage/www/projetos/VAMOSLANG/vscode-extension ~/.vscode/extensions/vamos-lang

# Ou no VS Code Insiders / Cursor
cp -r /storage/www/projetos/VAMOSLANG/vscode-extension ~/.cursor/extensions/vamos-lang
```

### Opção 2: Empacotar como arquivo .vsix
Se possuir o `vsce` instalado (`npm install -g @vscode/vsce`):
```bash
cd /storage/www/projetos/VAMOSLANG/vscode-extension
vsce package
code --install-extension vamos-lang-1.0.0.vsix
```

Reinicie o VS Code e abra qualquer arquivo `.vamos` para aproveitar a experiência completa!
