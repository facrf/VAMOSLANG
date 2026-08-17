/**
 * VAMOS-LANG Playground - Application Controller (app.js)
 */

const SNIPPETS = {
  "ola_mundo": `pacote principal

importar (
	"formatar"
	"tempo"
)

funcao principal() {
	variavel saudacao texto = "Olá, desenvolvedor VAMOS-LANG!"
	constante versao decimal64 = 1.0
	anoAtual := 2026

	formatar.ImprimirLinha("========================================")
	formatar.ImprimirLinha("  Bem-vindo ao Ecossistema VAMOS-LANG!  ")
	formatar.ImprimirLinha("========================================")

	formatar.ImprimirFormatado("Mensagem: %s\\n", saudacao)
	formatar.ImprimirFormatado("Versão: %.1f | Ano: %d\\n", versao, anoAtual)
	formatar.ImprimirFormatado("Executado em: %s\\n", tempo.Agora().Formato("02/01/2006 15:04:05"))
}`,

  "concorrencia": `pacote principal

importar (
	"formatar"
	"sincronizar"
	"tempo"
)

tipo Tarefa estrutura {
	Id   inteiro
	Nome texto
}

funcao processar(t Tarefa, canalSaida canal texto, wg *sincronizar.GrupoEspera) {
	adiar wg.Concluido()
	tempo.Dormir(tempo.Milissegundo * tempo.Duracao(100*t.Id))
	canalSaida <- formatar.Formatar("Tarefa #%d (%s) concluída!", t.Id, t.Nome)
}

funcao principal() {
	canalMensagens := criar(canal texto, 3)
	variavel wg sincronizar.GrupoEspera

	tarefas := []Tarefa{
		{Id: 1, Nome: "Download de Dados"},
		{Id: 2, Nome: "Processamento de Imagem"},
		{Id: 3, Nome: "Envio de Relatório"},
	}

	para _, tarefa := intervalo tarefas {
		wg.Adicionar(1)
		disparar processar(tarefa, canalMensagens, &wg)
	}

	disparar funcao() {
		wg.Esperar()
		fechar(canalMensagens)
	}()

	para msg := intervalo canalMensagens {
		formatar.ImprimirLinha("[RECEBIDO]", msg)
	}
}`,

  "estruturas": `pacote principal

importar (
	"formatar"
	"matematica"
)

tipo Forma interface {
	Area() decimal64
	Nome() texto
}

tipo Circulo estrutura {
	Raio decimal64
}

funcao (c Circulo) Area() decimal64 {
	retornar matematica.Pi * matematica.Potencia(c.Raio, 2)
}

funcao (c Circulo) Nome() texto {
	retornar "Círculo"
}

funcao exibir(f Forma) {
	formatar.ImprimirFormatado("Forma: %s | Área: %.2f cm²\\n", f.Nome(), f.Area())
}

funcao principal() {
	c := Circulo{Raio: 4.5}
	exibir(c)
}`,

  "calculadora": `pacote principal

importar (
	"erros"
	"formatar"
)

funcao dividir(a decimal64, b decimal64) (decimal64, erro) {
	se b == 0 {
		retornar 0, erros.Novo("não é possível dividir por zero")
	}
	retornar a / b, nulo
}

funcao principal() {
	valores := [][]decimal64{
		{10, 2},
		{50, 5},
		{100, 0},
	}

	para _, par := intervalo valores {
		res, err := dividir(par[0], par[1])
		se err != nulo {
			formatar.ImprimirLinha("[ERRO]", par[0], "/", par[1], "->", err)
		} senao {
			formatar.ImprimirLinha("[OK]  ", par[0], "/", par[1], "=", res)
		}
	}
}`,

  "servidor_web": `pacote principal

importar (
	"formatar"
	"servidor_http"
)

funcao manipularHome(w servidor_http.EscritorResposta, r *servidor_http.Requisicao) {
	formatar.FimprimirLinha(w, "Olá! Servidor HTTP 100% em VAMOS-LANG!")
}

funcao principal() {
	formatar.ImprimirLinha("Configurando rotas HTTP...")
	servidor_http.Manipulador("/", manipularHome)
	formatar.ImprimirLinha("Servidor pronto na porta :8080")
}`
};

// DOM Elements
const editorTextarea = document.getElementById("editorTextarea");
const lineNumbers = document.getElementById("lineNumbers");
const outputCode = document.getElementById("outputCode");
const snippetSelect = document.getElementById("snippetSelect");
const runBtn = document.getElementById("runBtn");
const copyVamosBtn = document.getElementById("copyVamosBtn");
const copyGoBtn = document.getElementById("copyGoBtn");
const downloadGoBtn = document.getElementById("downloadGoBtn");
const formatBtn = document.getElementById("formatBtn");
const shareBtn = document.getElementById("shareBtn");
const reverseBtn = document.getElementById("reverseBtn");
const clearBtn = document.getElementById("clearBtn");
const cheatSheetBtn = document.getElementById("cheatSheetBtn");
const cheatSheetModal = document.getElementById("cheatSheetModal");
const closeModalBtn = document.getElementById("closeModalBtn");
const searchKeyword = document.getElementById("searchKeyword");
const refGrid = document.getElementById("refGrid");
const terminalOutput = document.getElementById("terminalOutput");
const tokenTableBody = document.getElementById("tokenTableBody");
const tabBtns = document.querySelectorAll(".tab-btn");
const tabContents = document.querySelectorAll(".tab-content");
const toast = document.getElementById("toast");

let debounceTimer = null;

// Inicialização
document.addEventListener("DOMContentLoaded", () => {
  // Verifica se há código passado na URL via hash
  if (window.location.hash.startsWith("#code=")) {
    try {
      const b64 = window.location.hash.substring(6);
      const decoded = decodeURIComponent(escape(atob(b64)));
      editorTextarea.value = decoded;
      updateLineNumbers();
      runTranspilation();
      showToast("Código carregado a partir do link compartilhado!");
    } catch (e) {
      loadSnippet("ola_mundo");
    }
  } else {
    loadSnippet("ola_mundo");
  }

  populateCheatSheet();
  setupEventListeners();
});

function setupEventListeners() {
  // Evento de digitação no editor
  editorTextarea.addEventListener("input", () => {
    updateLineNumbers();
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(runTranspilation, 150);
  });

  // Sincronização de scroll entre números de linha e textarea
  editorTextarea.addEventListener("scroll", () => {
    lineNumbers.scrollTop = editorTextarea.scrollTop;
  });

  // Suporte a indentação com Tab no editor
  editorTextarea.addEventListener("keydown", (e) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const start = editorTextarea.selectionStart;
      const end = editorTextarea.selectionEnd;
      editorTextarea.value = editorTextarea.value.substring(0, start) + "\t" + editorTextarea.value.substring(end);
      editorTextarea.selectionStart = editorTextarea.selectionEnd = start + 1;
      updateLineNumbers();
      runTranspilation();
    }
  });

  // Seletor de Snippet
  snippetSelect.addEventListener("change", (e) => {
    loadSnippet(e.target.value);
  });

  // Botão Executar (Simulação no Console)
  runBtn.addEventListener("click", simulateExecution);

  // Formatar Código
  if (formatBtn) {
    formatBtn.addEventListener("click", () => {
      editorTextarea.value = formatVamosSource(editorTextarea.value);
      updateLineNumbers();
      runTranspilation();
      showToast("Código formatado com sucesso!");
    });
  }

  // Compartilhar Link com Código
  if (shareBtn) {
    shareBtn.addEventListener("click", () => {
      try {
        const b64 = btoa(unescape(encodeURIComponent(editorTextarea.value)));
        const shareUrl = `${window.location.origin}${window.location.pathname}#code=${b64}`;
        copyToClipboard(shareUrl, "Link compartilhável copiado para a área de transferência!");
      } catch (e) {
        showToast("Erro ao gerar link de compartilhamento.");
      }
    });
  }

  // Conversão Reversa Go -> VAMOS
  if (reverseBtn) {
    reverseBtn.addEventListener("click", () => {
      const goCode = prompt("Cole seu código Go nativo para converter em VAMOS-LANG:");
      if (goCode) {
        const res = transpileGoToVamos(goCode);
        if (res.success) {
          editorTextarea.value = res.vamosCode;
          updateLineNumbers();
          runTranspilation();
          showToast("Código Go convertido para VAMOS-LANG!");
        } else {
          alert(`Erro na conversão: ${res.error}`);
        }
      }
    });
  }

  // Copiar Códigos
  copyVamosBtn.addEventListener("click", () => {
    copyToClipboard(editorTextarea.value, "Código VAMOS copiado!");
  });

  copyGoBtn.addEventListener("click", () => {
    copyToClipboard(outputCode.textContent, "Código Go puro copiado!");
  });

  // Download do arquivo .go
  downloadGoBtn.addEventListener("click", () => {
    const code = outputCode.textContent;
    const blob = new Blob([code], { type: "text/plain;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "programa_transpilado.go";
    a.click();
    URL.revokeObjectURL(url);
    showToast("Download iniciado: programa_transpilado.go");
  });

  // Limpar Editor
  clearBtn.addEventListener("click", () => {
    editorTextarea.value = "pacote principal\n\nfuncao principal() {\n\t\n}\n";
    updateLineNumbers();
    runTranspilation();
    showToast("Editor limpo.");
  });

  // Modal Cheat Sheet
  cheatSheetBtn.addEventListener("click", () => {
    cheatSheetModal.classList.add("open");
    searchKeyword.focus();
  });

  closeModalBtn.addEventListener("click", () => {
    cheatSheetModal.classList.remove("open");
  });

  cheatSheetModal.addEventListener("click", (e) => {
    if (e.target === cheatSheetModal) {
      cheatSheetModal.classList.remove("open");
    }
  });

  searchKeyword.addEventListener("input", (e) => {
    filterCheatSheet(e.target.value);
  });

  // Tabs do Painel Inferior
  tabBtns.forEach(btn => {
    btn.addEventListener("click", () => {
      tabBtns.forEach(b => b.classList.remove("active"));
      tabContents.forEach(c => c.classList.remove("active"));
      btn.classList.add("active");
      const target = document.getElementById(btn.dataset.tab);
      if (target) target.classList.add("active");
    });
  });
}

function loadSnippet(key) {
  if (SNIPPETS[key]) {
    editorTextarea.value = SNIPPETS[key];
    updateLineNumbers();
    runTranspilation();
    clearTerminal();
    logTerminal("info", `Snippet carregado: "${key}". Clique em "Executar (Simular)" para testar.`);
  }
}

function updateLineNumbers() {
  const lines = editorTextarea.value.split("\n").length;
  let nums = "";
  for (let i = 1; i <= lines; i++) {
    nums += i + "<br>";
  }
  lineNumbers.innerHTML = nums;
}

function runTranspilation() {
  const code = editorTextarea.value;
  const result = transpileVamos(code);

  if (result.success) {
    outputCode.textContent = result.goCode;
    renderTokens(result.tokens);
  } else {
    outputCode.textContent = `// ERRO DE TRANSPILAÇÃO:\n// ${result.error}`;
  }
}

function renderTokens(tokens) {
  if (!tokenTableBody) return;
  tokenTableBody.innerHTML = "";

  tokens.slice(0, 80).forEach(tok => {
    if (tok.type === "WS" || tok.type === "NEWLINE") return;
    const tr = document.createElement("tr");
    
    let typeClass = "";
    if (tok.type === "KEYWORD") typeClass = "tok-kw";
    else if (tok.type === "STRING" || tok.type === "RAW_STRING") typeClass = "tok-str";
    else if (tok.type === "IDENT") typeClass = "tok-ident";
    else if (tok.type === "OP") typeClass = "tok-op";

    tr.innerHTML = `
      <td>${tok.line}:${tok.col}</td>
      <td class="${typeClass}">${tok.type}</td>
      <td><code>${escapeHtml(tok.val)}</code></td>
    `;
    tokenTableBody.appendChild(tr);
  });
}

function simulateExecution() {
  const code = editorTextarea.value;
  clearTerminal();
  logTerminal("info", `[VAMOS CLI] Iniciando transpilação e execução simulada...`);

  const result = transpileVamos(code);
  if (!result.success) {
    logTerminal("error", `[FALHA] Erro léxico/sintático: ${result.error}`);
    return;
  }

  logTerminal("success", `[OK] Transpilação concluída com sucesso (0.001s)`);
  logTerminal("info", `[EXEC] Executando binário gerado nativamente:`);
  logTerminal("info", `--------------------------------------------------`);

  // Simulação inteligente baseada no código
  setTimeout(() => {
    if (code.includes("Olá, desenvolvedor VAMOS-LANG!")) {
      logTerminal("normal", "========================================");
      logTerminal("normal", "  Bem-vindo ao Ecossistema VAMOS-LANG!  ");
      logTerminal("normal", "========================================");
      logTerminal("normal", "Mensagem: Olá, desenvolvedor VAMOS-LANG!");
      logTerminal("normal", "Versão: 1.0 | Ano: 2026");
      logTerminal("normal", `Executado em: ${new Date().toLocaleString("pt-BR")}`);
      logTerminal("normal", "========================================");
    } else if (code.includes("Tarefa estrutura") || code.includes("disparar")) {
      logTerminal("normal", "[RECEBIDO] Tarefa #1 (Download de Dados) concluída!");
      logTerminal("normal", "[RECEBIDO] Tarefa #2 (Processamento de Imagem) concluída!");
      logTerminal("normal", "[RECEBIDO] Tarefa #3 (Envio de Relatório) concluída!");
      logTerminal("success", "Todas as rotinas concorrentes concluídas com êxito!");
    } else if (code.includes("Forma interface")) {
      logTerminal("normal", "Forma: Círculo | Área: 63.62 cm²");
    } else if (code.includes("dividir")) {
      logTerminal("normal", "[OK]   10 / 2 = 5");
      logTerminal("normal", "[OK]   50 / 5 = 10");
      logTerminal("warn", "[ERRO] 100 / 0 -> não é possível dividir por zero");
    } else if (code.includes("servidor_http")) {
      logTerminal("normal", "Configurando rotas HTTP...");
      logTerminal("normal", "Servidor pronto na porta :8080");
      logTerminal("info", "Manipulador GET / registrado.");
    } else {
      logTerminal("normal", "Programa executado com código de saída 0 (sucesso).");
    }
  }, 100);
}

function clearTerminal() {
  terminalOutput.innerHTML = "";
}

function logTerminal(type, text) {
  const line = document.createElement("div");
  line.className = `terminal-line ${type}`;
  line.textContent = text;
  terminalOutput.appendChild(line);
  terminalOutput.scrollTop = terminalOutput.scrollHeight;
}

function populateCheatSheet() {
  refGrid.innerHTML = "";
  
  const allEntries = [
    ...Object.entries(VAMOS_KEYWORDS).map(([v, g]) => ({ v, g, cat: "Palavra-Chave" })),
    ...Object.entries(VAMOS_TYPES_BUILTINS).map(([v, g]) => ({ v, g, cat: "Tipo / Built-in" })),
    ...Object.entries(VAMOS_STDLIB_IMPORTS).map(([v, g]) => ({ v: `"${v}"`, g: `"${g}"`, cat: "Pacote Padrão" })),
    ...Object.entries(VAMOS_METHOD_ALIASES).map(([v, g]) => ({ v: `.${v}()`, g: `.${g}()`, cat: "Método / Interface" }))
  ];

  allEntries.forEach(item => {
    const card = document.createElement("div");
    card.className = "ref-card";
    card.innerHTML = `
      <span class="ref-vamos">${escapeHtml(item.v)}</span>
      <span class="ref-go">↳ ${escapeHtml(item.g)}</span>
    `;
    card.addEventListener("click", () => {
      insertAtCursor(item.v.replace(/[".()]/g, ""));
      cheatSheetModal.classList.remove("open");
      showToast(`Inserido: ${item.v}`);
    });
    refGrid.appendChild(card);
  });
}

function filterCheatSheet(query) {
  const term = query.toLowerCase().trim();
  const cards = refGrid.querySelectorAll(".ref-card");
  cards.forEach(card => {
    const txt = card.textContent.toLowerCase();
    card.style.display = txt.includes(term) ? "flex" : "none";
  });
}

function insertAtCursor(text) {
  const start = editorTextarea.selectionStart;
  const end = editorTextarea.selectionEnd;
  const before = editorTextarea.value.substring(0, start);
  const after = editorTextarea.value.substring(end);
  editorTextarea.value = before + text + after;
  editorTextarea.selectionStart = editorTextarea.selectionEnd = start + text.length;
  editorTextarea.focus();
  updateLineNumbers();
  runTranspilation();
}

function copyToClipboard(text, msg) {
  navigator.clipboard.writeText(text).then(() => {
    showToast(msg);
  }).catch(() => {
    showToast("Erro ao copiar.");
  });
}

function showToast(msg) {
  toast.textContent = msg;
  toast.classList.add("show");
  setTimeout(() => {
    toast.classList.remove("show");
  }, 2500);
}

function escapeHtml(str) {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}
