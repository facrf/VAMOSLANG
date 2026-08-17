/**
 * VAMOS-LANG Transpiler para JavaScript (Web Playground)
 * Converte código VAMOS-LANG (PT-BR) em código Go puro (Golang) diretamente no navegador.
 */

const VAMOS_KEYWORDS = {
  "pacote": "package",
  "principal": "main",
  "importar": "import",
  "funcao": "func",
  "variavel": "var",
  "constante": "const",
  "tipo": "type",
  "estrutura": "struct",
  "interface": "interface",
  "mapa": "map",
  "canal": "chan",
  "se": "if",
  "senao": "else",
  "para": "for",
  "intervalo": "range",
  "escolha": "switch",
  "caso": "case",
  "padrao": "default",
  "prosseguir": "fallthrough",
  "interromper": "break",
  "continuar": "continue",
  "ir_para": "goto",
  "retornar": "return",
  "disparar": "go",
  "adiar": "defer",
  "selecionar": "select"
};

const VAMOS_TYPES_BUILTINS = {
  "texto": "string",
  "inteiro": "int",
  "inteiro8": "int8",
  "inteiro16": "int16",
  "inteiro32": "int32",
  "inteiro64": "int64",
  "uinteiro": "uint",
  "uinteiro8": "uint8",
  "uinteiro16": "uint16",
  "uinteiro32": "uint32",
  "uinteiro64": "uint64",
  "uinteiro_ponteiro": "uintptr",
  "decimal": "float64",
  "decimal32": "float32",
  "decimal64": "float64",
  "complexo64": "complex64",
  "complexo128": "complex128",
  "booleano": "bool",
  "byte": "byte",
  "runa": "rune",
  "caractere": "rune",
  "qualquer": "any",
  "erro": "error",
  "verdadeiro": "true",
  "falso": "false",
  "nulo": "nil",
  "iota": "iota",
  "criar": "make",
  "novo": "new",
  "adicionar": "append",
  "anexar": "append",
  "tamanho": "len",
  "capacidade": "cap",
  "fechar": "close",
  "copiar": "copy",
  "deletar": "delete",
  "excluir": "delete",
  "panico": "panic",
  "recuperar": "recover",
  "imprimir": "print",
  "imprimirln": "println"
};

const VAMOS_STDLIB_IMPORTS = {
  "formatar": "fmt",
  "tempo": "time",
  "matematica": "math",
  "cordas": "strings",
  "textos": "strings",
  "so": "os",
  "sistema": "os",
  "sincronizar": "sync",
  "erros": "errors",
  "io": "io",
  "leitor": "bufio",
  "json": "encoding/json",
  "servidor_http": "net/http",
  "http": "net/http",
  "contexto": "context",
  "registrador": "log",
  "caminho": "path/filepath",
  "ordenacao": "sort",
  "teste": "testing"
};

const VAMOS_STDLIB_IDENTS = {
  "formatar": "fmt",
  "fmt": "fmt",
  "tempo": "time",
  "time": "time",
  "matematica": "math",
  "math": "math",
  "cordas": "strings",
  "textos": "strings",
  "strings": "strings",
  "so": "os",
  "sistema": "os",
  "os": "os",
  "sincronizar": "sync",
  "sync": "sync",
  "erros": "errors",
  "errors": "errors",
  "io": "io",
  "leitor": "bufio",
  "bufio": "bufio",
  "json": "json",
  "servidor_http": "http",
  "http": "http",
  "contexto": "context",
  "context": "context",
  "registrador": "log",
  "log": "log",
  "caminho": "filepath",
  "filepath": "filepath",
  "ordenacao": "sort",
  "sort": "sort",
  "teste": "testing",
  "testing": "testing"
};

const VAMOS_STDLIB_METHODS = {
  "fmt": {
    "ImprimirLinha": "Println",
    "ImprimirFormatado": "Printf",
    "Imprimir": "Print",
    "FormatarLinha": "Sprintln",
    "Formatar": "Sprintf",
    "FormatarSimples": "Sprint",
    "CriarErro": "Errorf",
    "ErroFormatado": "Errorf",
    "Escanear": "Scan",
    "EscanearLinha": "Scanln",
    "EscanearFormatado": "Scanf",
    "Fimprimir": "Fprint",
    "FimprimirLinha": "Fprintln",
    "FimprimirFormatado": "Fprintf"
  },
  "time": {
    "Dormir": "Sleep",
    "Agora": "Now",
    "Segundo": "Second",
    "Milissegundo": "Millisecond",
    "Microssegundo": "Microsecond",
    "Nanosegundo": "Nanosecond",
    "Minuto": "Minute",
    "Hora": "Hour",
    "Duracao": "Duration",
    "Desde": "Since",
    "Apos": "After",
    "NovoTemporizador": "NewTimer",
    "NovoTique": "NewTicker",
    "Tique": "Tick",
    "Data": "Date",
    "CarregarLocalizacao": "LoadLocation",
    "FormatoUnix": "Unix"
  },
  "math": {
    "Pi": "Pi",
    "Raiz": "Sqrt",
    "Potencia": "Pow",
    "Absoluto": "Abs",
    "Piso": "Floor",
    "Teto": "Ceil",
    "Maximo": "Max",
    "Minimo": "Min",
    "Seno": "Sin",
    "Cosseno": "Cos",
    "Tangente": "Tan",
    "Logaritmo": "Log",
    "Mod": "Mod"
  },
  "strings": {
    "Contem": "Contains",
    "Dividir": "Split",
    "Juntar": "Join",
    "ParaMaiusculas": "ToUpper",
    "ParaMinusculas": "ToLower",
    "Substituir": "Replace",
    "SubstituirTudo": "ReplaceAll",
    "TemPrefixo": "HasPrefix",
    "TemSufixo": "HasSuffix",
    "Aparar": "Trim",
    "ApararEspaco": "TrimSpace",
    "Indice": "Index",
    "Repetir": "Repeat",
    "Contar": "Count",
    "Comparar": "Compare"
  },
  "sync": {
    "GrupoEspera": "WaitGroup",
    "Trava": "Mutex",
    "TravaLeituraEscrita": "RWMutex",
    "UmaVez": "Once"
  },
  "errors": {
    "Novo": "New",
    "Desempacotar": "Unwrap",
    "E": "Is",
    "Como": "As"
  },
  "os": {
    "Sair": "Exit",
    "Argumentos": "Args",
    "LerArquivo": "ReadFile",
    "EscreverArquivo": "WriteFile",
    "ObterVariavelAmbiente": "Getenv",
    "DefinirVariavelAmbiente": "Setenv",
    "CriarDiretorio": "Mkdir"
  },
  "http": {
    "OuvirEServir": "ListenAndServe",
    "Manipulador": "HandleFunc",
    "NovoRequisicao": "NewRequest",
    "Cliente": "Client",
    "Resposta": "Response",
    "Requisicao": "Request",
    "EscritorResposta": "ResponseWriter",
    "StatusOk": "StatusOK"
  }
};

const VAMOS_METHOD_ALIASES = {
  "Adicionar": "Add",
  "Concluido": "Done",
  "Feito": "Done",
  "Esperar": "Wait",
  "Bloquear": "Lock",
  "Desbloquear": "Unlock",
  "Fechar": "Close",
  "Ler": "Read",
  "Escrever": "Write",
  "Texto": "String",
  "Erro": "Error",
  "Formato": "Format",
  "Formatar": "Format",
  "Unix": "Unix",
  "Segundos": "Seconds",
  "Milissegundos": "Milliseconds",
  "Microssegundos": "Microseconds",
  "Nanosegundos": "Nanoseconds",
  "Minutos": "Minutes",
  "Horas": "Hours"
};

/**
 * Tokenizador Léxico em JavaScript
 */
class VamosLexer {
  constructor(input) {
    this.src = input;
    this.pos = 0;
    this.len = input.length;
    this.line = 1;
    this.col = 1;
  }

  peek(offset = 0) {
    const idx = this.pos + offset;
    return idx < this.len ? this.src[idx] : "";
  }

  next() {
    if (this.pos >= this.len) return "";
    const ch = this.src[this.pos++];
    if (ch === "\n") {
      this.line++;
      this.col = 1;
    } else {
      this.col++;
    }
    return ch;
  }

  isLetter(ch) {
    return ch === "_" || (ch >= "a" && ch <= "z") || (ch >= "A" && ch <= "Z") || ch.charCodeAt(0) > 127;
  }

  isDigit(ch) {
    return ch >= "0" && ch <= "9";
  }

  isIdentPart(ch) {
    return this.isLetter(ch) || this.isDigit(ch);
  }

  tokenize() {
    const tokens = [];
    while (this.pos < this.len) {
      const curLine = this.line;
      const curCol = this.col;
      const ch = this.peek();

      // Whitespace
      if (ch === " " || ch === "\t" || ch === "\r") {
        let val = "";
        while (this.pos < this.len && (this.peek() === " " || this.peek() === "\t" || this.peek() === "\r")) {
          val += this.next();
        }
        tokens.push({ type: "WS", val, line: curLine, col: curCol });
        continue;
      }

      // Newline
      if (ch === "\n") {
        this.next();
        tokens.push({ type: "NEWLINE", val: "\n", line: curLine, col: curCol });
        continue;
      }

      // Single-line Comment
      if (ch === "/" && this.peek(1) === "/") {
        let val = "";
        while (this.pos < this.len && this.peek() !== "\n") {
          val += this.next();
        }
        tokens.push({ type: "COMMENT", val, line: curLine, col: curCol });
        continue;
      }

      // Multi-line Comment
      if (ch === "/" && this.peek(1) === "*") {
        let val = this.next() + this.next();
        while (this.pos < this.len) {
          if (this.peek() === "*" && this.peek(1) === "/") {
            val += this.next() + this.next();
            break;
          }
          val += this.next();
        }
        tokens.push({ type: "COMMENT", val, line: curLine, col: curCol });
        continue;
      }

      // Interpreted String Literal
      if (ch === '"') {
        let val = this.next();
        let escaped = false;
        while (this.pos < this.len) {
          const c = this.next();
          val += c;
          if (escaped) {
            escaped = false;
          } else if (c === "\\") {
            escaped = true;
          } else if (c === '"') {
            break;
          }
        }
        tokens.push({ type: "STRING", val, line: curLine, col: curCol });
        continue;
      }

      // Raw String Literal
      if (ch === "`") {
        let val = this.next();
        while (this.pos < this.len) {
          const c = this.next();
          val += c;
          if (c === "`") break;
        }
        tokens.push({ type: "RAW_STRING", val, line: curLine, col: curCol });
        continue;
      }

      // Rune Literal
      if (ch === "'") {
        let val = this.next();
        let escaped = false;
        while (this.pos < this.len) {
          const c = this.next();
          val += c;
          if (escaped) {
            escaped = false;
          } else if (c === "\\") {
            escaped = true;
          } else if (c === "'") {
            break;
          }
        }
        tokens.push({ type: "RUNE", val, line: curLine, col: curCol });
        continue;
      }

      // Numbers
      if (this.isDigit(ch) || (ch === "." && this.isDigit(this.peek(1)))) {
        let val = "";
        if (ch === "0" && ["x", "X", "o", "O", "b", "B"].includes(this.peek(1))) {
          val += this.next() + this.next();
          while (this.pos < this.len && (this.isIdentPart(this.peek()) || this.peek() === "_")) {
            val += this.next();
          }
          tokens.push({ type: "NUM", val, line: curLine, col: curCol });
          continue;
        }

        let hasDot = false;
        while (this.pos < this.len) {
          const c = this.peek();
          if (c === "." && !hasDot && this.isDigit(this.peek(1))) {
            hasDot = true;
            val += this.next();
          } else if (this.isDigit(c) || c === "_" || c === "e" || c === "E" || c === "+" || c === "-") {
            val += this.next();
          } else {
            break;
          }
        }
        tokens.push({ type: "NUM", val, line: curLine, col: curCol });
        continue;
      }

      // Identifiers & Keywords
      if (this.isLetter(ch)) {
        let val = "";
        while (this.pos < this.len && this.isIdentPart(this.peek())) {
          val += this.next();
        }
        if (VAMOS_KEYWORDS[val]) {
          tokens.push({ type: "KEYWORD", val, line: curLine, col: curCol });
        } else {
          tokens.push({ type: "IDENT", val, line: curLine, col: curCol });
        }
        continue;
      }

      // Operators and Delimiters
      const tri = this.src.substring(this.pos, this.pos + 3);
      if (["...", "<<=", ">>=", "&^="].includes(tri)) {
        this.next(); this.next(); this.next();
        tokens.push({ type: "OP", val: tri, line: curLine, col: curCol });
        continue;
      }

      const bi = this.src.substring(this.pos, this.pos + 2);
      if ([":=", "==", "!=", "<=", ">=", "&&", "||", "<-", "++", "--", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>", "&^"].includes(bi)) {
        this.next(); this.next();
        tokens.push({ type: "OP", val: bi, line: curLine, col: curCol });
        continue;
      }

      const single = this.next();
      if ("()[]{};,:.".includes(single)) {
        tokens.push({ type: "DELIM", val: single, line: curLine, col: curCol });
      } else {
        tokens.push({ type: "OP", val: single, line: curLine, col: curCol });
      }
    }

    tokens.push({ type: "EOF", val: "", line: this.line, col: this.col });
    return tokens;
  }
}

/**
 * Transpilador VAMOS-LANG para Go
 */
class VamosTranspiler {
  constructor(tokens) {
    this.tokens = tokens;
    this.pos = 0;
    this.len = tokens.length;
  }

  peek(offset = 0) {
    const idx = this.pos + offset;
    return idx < this.len ? this.tokens[idx] : { type: "EOF", val: "" };
  }

  next() {
    return this.pos < this.len ? this.tokens[this.pos++] : { type: "EOF", val: "" };
  }

  transpile() {
    let out = "";

    while (this.pos < this.len) {
      const tok = this.peek();
      if (tok.type === "EOF") break;

      // Preserva whitespace, comentários, literais brutos e caracteres
      if (["WS", "NEWLINE", "COMMENT", "RAW_STRING", "RUNE"].includes(tok.type)) {
        out += this.next().val;
        continue;
      }

      // Cláusula de Importação
      if (tok.type === "KEYWORD" && tok.val === "importar") {
        this.next();
        out += "import";
        out += this.processImportClause();
        continue;
      }

      // Pacote
      if (tok.type === "KEYWORD" && tok.val === "pacote") {
        this.next();
        out += "package";
        continue;
      }

      // Palavras-chave
      if (tok.type === "KEYWORD") {
        const val = this.next().val;
        out += VAMOS_KEYWORDS[val] || val;
        continue;
      }

      // Strings fora de imports
      if (tok.type === "STRING") {
        out += this.next().val;
        continue;
      }

      // Números
      if (tok.type === "NUM") {
        out += this.next().val;
        continue;
      }

      // Identificadores (membro pacote.Metodo ou objeto.Metodo)
      if (tok.type === "IDENT") {
        const nextTok = this.peek(1);
        if (nextTok.type === "DELIM" && nextTok.val === ".") {
          const nextNext = this.peek(2);
          if (nextNext.type === "IDENT") {
            const pkgName = tok.val;
            const methodName = nextNext.val;

            // Pacote padrão mapeado
            const goIdent = VAMOS_STDLIB_IDENTS[pkgName];
            if (goIdent) {
              const goMethod = (VAMOS_STDLIB_METHODS[goIdent] && VAMOS_STDLIB_METHODS[goIdent][methodName]) ||
                               VAMOS_METHOD_ALIASES[methodName];
              if (goMethod) {
                this.next(); this.next(); this.next();
                out += `${goIdent}.${goMethod}`;
                continue;
              }

              this.next(); this.next();
              out += `${goIdent}.`;
              continue;
            }

            // Método em objeto
            const alias = VAMOS_METHOD_ALIASES[methodName];
            if (alias) {
              const identTok = this.next();
              this.next(); this.next();
              const mappedIdent = VAMOS_TYPES_BUILTINS[identTok.val] || identTok.val;
              out += `${mappedIdent}.${alias}`;
              continue;
            }
          }
        }

        const val = this.next().val;
        if (VAMOS_STDLIB_IDENTS[val]) {
          out += VAMOS_STDLIB_IDENTS[val];
        } else if (VAMOS_TYPES_BUILTINS[val]) {
          out += VAMOS_TYPES_BUILTINS[val];
        } else {
          out += val;
        }
        continue;
      }

      // Delimitador '.' seguido de método com alias (ex: obj.Adicionar())
      if (tok.type === "DELIM" && tok.val === ".") {
        const nextTok = this.peek(1);
        if (nextTok.type === "IDENT" && VAMOS_METHOD_ALIASES[nextTok.val]) {
          this.next(); this.next();
          out += `.${VAMOS_METHOD_ALIASES[nextTok.val]}`;
          continue;
        }
      }

      // Demais operadores / delimitadores
      out += this.next().val;
    }

    return out;
  }

  processImportClause() {
    let out = "";
    let inParen = false;

    while (this.pos < this.len) {
      const tok = this.peek();
      if (tok.type === "EOF") break;

      if (tok.type === "WS" || tok.type === "COMMENT") {
        out += this.next().val;
        continue;
      }

      if (tok.type === "DELIM" && tok.val === "(") {
        inParen = true;
        out += this.next().val;
        continue;
      }

      if (tok.type === "DELIM" && tok.val === ")") {
        inParen = false;
        out += this.next().val;
        break;
      }

      if (tok.type === "NEWLINE") {
        out += this.next().val;
        if (!inParen) break;
        continue;
      }

      if (tok.type === "DELIM" && tok.val === ";") {
        out += this.next().val;
        if (!inParen) break;
        continue;
      }

      if (tok.type === "STRING") {
        const raw = tok.val;
        if (raw.startsWith('"') && raw.endsWith('"')) {
          const unq = raw.substring(1, raw.length - 1);
          if (VAMOS_STDLIB_IMPORTS[unq]) {
            out += `"${VAMOS_STDLIB_IMPORTS[unq]}"`;
            this.next();
            if (!inParen) {
              while (this.pos < this.len && (this.peek().type === "WS" || this.peek().type === "COMMENT")) {
                out += this.next().val;
              }
              if (this.pos < this.len && this.peek().type === "NEWLINE") {
                out += this.next().val;
              }
              break;
            }
            continue;
          }
        }
        out += this.next().val;
        if (!inParen) break;
        continue;
      }

      if (tok.type === "IDENT") {
        out += this.next().val;
        continue;
      }

      out += this.next().val;
    }

    return out;
  }
}

/**
 * Transpilador Reverso: Go -> VAMOS-LANG
 */
function transpileGoToVamos(goCode) {
  try {
    const lexer = new VamosLexer(goCode);
    const tokens = lexer.tokenize();
    let out = "";
    let pos = 0;
    const len = tokens.length;

    const reverseKw = Object.fromEntries(Object.entries(VAMOS_KEYWORDS).map(([k, v]) => [v, k]));
    const reverseTypes = Object.fromEntries(Object.entries(VAMOS_TYPES_BUILTINS).map(([k, v]) => [v, k]));
    const reverseImports = Object.fromEntries(Object.entries(VAMOS_STDLIB_IMPORTS).map(([k, v]) => [v, k]));

    while (pos < len) {
      const tok = tokens[pos];
      if (tok.type === "EOF") break;

      if (["WS", "NEWLINE", "COMMENT", "RAW_STRING", "RUNE"].includes(tok.type)) {
        out += tok.val;
        pos++;
        continue;
      }

      if (tok.type === "KEYWORD" && tok.val === "import") {
        pos++;
        out += "importar";
        let inParen = false;
        while (pos < len) {
          const t = tokens[pos];
          if (t.type === "EOF") break;
          if (t.type === "WS" || t.type === "COMMENT") {
            out += t.val;
            pos++;
            continue;
          }
          if (t.type === "DELIM" && t.val === "(") {
            inParen = true;
            out += t.val;
            pos++;
            continue;
          }
          if (t.type === "DELIM" && t.val === ")") {
            inParen = false;
            out += t.val;
            pos++;
            break;
          }
          if (t.type === "NEWLINE") {
            out += t.val;
            pos++;
            if (!inParen) break;
            continue;
          }
          if (t.type === "STRING") {
            const raw = t.val;
            if (raw.startsWith('"') && raw.endsWith('"')) {
              const unq = raw.substring(1, raw.length - 1);
              if (reverseImports[unq]) {
                out += `"${reverseImports[unq]}"`;
                pos++;
                if (!inParen) break;
                continue;
              }
            }
            out += t.val;
            pos++;
            if (!inParen) break;
            continue;
          }
          out += t.val;
          pos++;
        }
        continue;
      }

      if (tok.type === "KEYWORD" && tok.val === "package") {
        pos++;
        out += "pacote";
        continue;
      }

      if (tok.type === "KEYWORD") {
        out += reverseKw[tok.val] || tok.val;
        pos++;
        continue;
      }

      if (tok.type === "STRING" || tok.type === "NUM") {
        out += tok.val;
        pos++;
        continue;
      }

      if (tok.type === "IDENT") {
        if (pos + 2 < len && tokens[pos+1].type === "DELIM" && tokens[pos+1].val === "." && tokens[pos+2].type === "IDENT") {
          const pkg = tok.val;
          const method = tokens[pos+2].val;
          let mappedPkg = reverseImports[pkg] || pkg;
          let mappedMethod = method;

          if (VAMOS_STDLIB_METHODS[pkg]) {
            for (const [ptM, goM] of Object.entries(VAMOS_STDLIB_METHODS[pkg])) {
              if (goM === method) {
                mappedMethod = ptM;
                break;
              }
            }
          } else {
            for (const [ptM, goM] of Object.entries(VAMOS_METHOD_ALIASES)) {
              if (goM === method) {
                mappedMethod = ptM;
                break;
              }
            }
          }

          out += `${mappedPkg}.${mappedMethod}`;
          pos += 3;
          continue;
        }

        if (reverseTypes[tok.val]) {
          out += reverseTypes[tok.val];
        } else {
          out += tok.val;
        }
        pos++;
        continue;
      }

      out += tok.val;
      pos++;
    }

    return { success: true, vamosCode: out };
  } catch (err) {
    return { success: false, error: err.message || String(err) };
  }
}

/**
 * Função pública para formatar código VAMOS-LANG
 */
function formatVamosSource(source) {
  const lines = source.split("\n");
  const formatted = [];
  let indent = 0;

  for (const line of lines) {
    const trimmed = line.trim();
    if (trimmed === "") {
      if (formatted.length > 0 && formatted[formatted.length - 1] !== "") {
        formatted.push("");
      }
      continue;
    }

    if (trimmed.startsWith("}") || trimmed.startsWith(")") || trimmed.startsWith("]")) {
      if (indent > 0) indent--;
    }

    const prefix = "\t".repeat(indent);
    formatted.push(prefix + trimmed);

    if (trimmed.endsWith("{") || trimmed.endsWith("(") || trimmed.endsWith("[")) {
      indent++;
    }
  }

  let res = formatted.join("\n");
  if (!res.endsWith("\n")) res += "\n";
  return res;
}

/**
 * Função pública para transpilar código VAMOS-LANG
 */
function transpileVamos(code) {
  try {
    const lexer = new VamosLexer(code);
    const tokens = lexer.tokenize();
    const transpiler = new VamosTranspiler(tokens);
    return {
      success: true,
      goCode: transpiler.transpile(),
      tokens: tokens
    };
  } catch (err) {
    return {
      success: false,
      error: err.message || String(err)
    };
  }
}
