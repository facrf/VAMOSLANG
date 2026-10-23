const vscode = require('vscode');
const { spawn } = require('child_process');

let lspProcess = null;
let msgId = 1;
const pendingRequests = new Map();
let docVersions = new Map();

function sendLSPRequest(method, params) {
  return new Promise((resolve, reject) => {
    if (!lspProcess || lspProcess.killed) {
      return reject(new Error('LSP não está em execução'));
    }

    const id = msgId++;
    const payload = JSON.stringify({
      jsonrpc: '2.0',
      id: id,
      method: method,
      params: params
    });

    const header = `Content-Length: ${Buffer.byteLength(payload, 'utf8')}\r\n\r\n`;
    pendingRequests.set(id, { resolve, reject });
    lspProcess.stdin.write(header + payload);
  });
}

function sendLSPNotification(method, params) {
  if (!lspProcess || lspProcess.killed) return;
  const payload = JSON.stringify({
    jsonrpc: '2.0',
    method: method,
    params: params
  });
  const header = `Content-Length: ${Buffer.byteLength(payload, 'utf8')}\r\n\r\n`;
  lspProcess.stdin.write(header + payload);
}

function startLSP(diagnosticCollection) {
  try {
    const config = vscode.workspace.getConfiguration('vamos');
    const binPath = config.get('lspPath') || 'vamos';

    lspProcess = spawn(binPath, ['lsp']);

    let buffer = Buffer.alloc(0);

    lspProcess.stdout.on('data', (chunk) => {
      buffer = Buffer.concat([buffer, chunk]);

      while (true) {
        const headerEnd = buffer.indexOf('\r\n\r\n');
        if (headerEnd === -1) break;

        const headerStr = buffer.slice(0, headerEnd).toString('utf8');
        const match = headerStr.match(/Content-Length:\s*(\d+)/i);
        if (!match) {
          buffer = buffer.slice(headerEnd + 4);
          continue;
        }

        const contentLength = parseInt(match[1], 10);
        const totalMsgLength = headerEnd + 4 + contentLength;
        if (buffer.length < totalMsgLength) {
          break; // Aguarda mais dados
        }

        const bodyBuf = buffer.slice(headerEnd + 4, totalMsgLength);
        buffer = buffer.slice(totalMsgLength);

        try {
          const msg = JSON.parse(bodyBuf.toString('utf8'));

          // 1. Diagnósticos recebidos
          if (msg.method === 'textDocument/publishDiagnostics' && msg.params) {
            const uri = vscode.Uri.parse(msg.params.uri);
            const diags = (msg.params.diagnostics || []).map((d) => {
              const range = new vscode.Range(
                d.range.start.line,
                d.range.start.character,
                d.range.end.line,
                d.range.end.character
              );
              const severity =
                d.severity === 1
                  ? vscode.DiagnosticSeverity.Error
                  : d.severity === 2
                  ? vscode.DiagnosticSeverity.Warning
                  : vscode.DiagnosticSeverity.Information;
              return new vscode.Diagnostic(range, d.message, severity);
            });
            diagnosticCollection.set(uri, diags);
          }

          // 2. Resposta de requisição síncrona
          if (msg.id && pendingRequests.has(msg.id)) {
            const { resolve } = pendingRequests.get(msg.id);
            pendingRequests.delete(msg.id);
            resolve(msg.result);
          }
        } catch (e) {
          console.error('Erro ao processar mensagem LSP:', e);
        }
      }
    });

    lspProcess.on('error', (err) => {
      console.warn('VAMOS LSP não pôde ser iniciado:', err.message);
    });

    lspProcess.on('exit', () => {
      lspProcess = null;
    });

    // Envia Initialize
    sendLSPRequest('initialize', {
      capabilities: {},
      processId: process.pid,
      rootUri: vscode.workspace.workspaceFolders?.[0]?.uri.toString() || null
    }).then(() => {
      sendLSPNotification('initialized', {});
    }).catch(() => {});
  } catch (err) {
    console.error('Falha ao inicializar LSP do VAMOS-LANG:', err);
  }
}

function activate(context) {
  const diagnosticCollection = vscode.languages.createDiagnosticCollection('vamos');
  context.subscriptions.push(diagnosticCollection);

  // Inicia Language Server
  startLSP(diagnosticCollection);

  // Monitora abertura e alteração de arquivos .vamos
  context.subscriptions.push(
    vscode.workspace.onDidOpenTextDocument((doc) => {
      if (doc.languageId === 'vamos') {
        sendLSPNotification('textDocument/didOpen', {
          textDocument: {
            uri: doc.uri.toString(),
            languageId: 'vamos',
            version: doc.version,
            text: doc.getText()
          }
        });
      }
    })
  );

  context.subscriptions.push(
    vscode.workspace.onDidChangeTextDocument((event) => {
      if (event.document.languageId === 'vamos') {
        sendLSPNotification('textDocument/didChange', {
          textDocument: {
            uri: event.document.uri.toString(),
            version: event.document.version
          },
          contentChanges: [{ text: event.document.getText() }]
        });
      }
    })
  );

  context.subscriptions.push(
    vscode.workspace.onDidSaveTextDocument((doc) => {
      if (doc.languageId === 'vamos') {
        sendLSPNotification('textDocument/didSave', {
          textDocument: { uri: doc.uri.toString() }
        });
      }
    })
  );

  // 1. Provedor de Hover Inteligente
  const hoverProvider = vscode.languages.registerHoverProvider('vamos', {
    async provideHover(document, position) {
      try {
        const res = await sendLSPRequest('textDocument/hover', {
          textDocument: { uri: document.uri.toString() },
          position: { line: position.line, character: position.character }
        });
        if (res && res.contents && res.contents.value) {
          return new vscode.Hover(new vscode.MarkdownString(res.contents.value));
        }
      } catch (e) {}
      return null;
    }
  });
  context.subscriptions.push(hoverProvider);

  // 2. Provedor de Autocompletar (Completions)
  const completionProvider = vscode.languages.registerCompletionItemProvider(
    'vamos',
    {
      async provideCompletionItems(document, position) {
        try {
          const res = await sendLSPRequest('textDocument/completion', {
            textDocument: { uri: document.uri.toString() },
            position: { line: position.line, character: position.character }
          });
          if (Array.isArray(res)) {
            return res.map((item) => {
              const ci = new vscode.CompletionItem(item.label, item.kind || vscode.CompletionItemKind.Keyword);
              ci.detail = item.detail;
              ci.documentation = new vscode.MarkdownString(item.documentation || '');
              if (item.insertText) {
                ci.insertText = item.insertText;
              }
              return ci;
            });
          }
        } catch (e) {}
        return [];
      }
    },
    '.',
    ' '
  );
  context.subscriptions.push(completionProvider);

  // 3. Provedor de Formatação (Format Document)
  const formattingProvider = vscode.languages.registerDocumentFormattingEditProvider(
    'vamos',
    {
      async provideDocumentFormattingEdits(document) {
        try {
          const edits = await sendLSPRequest('textDocument/formatting', {
            textDocument: { uri: document.uri.toString() }
          });
          if (Array.isArray(edits) && edits.length > 0) {
            return edits.map((e) => {
              const range = new vscode.Range(
                e.range.start.line,
                e.range.start.character,
                e.range.end.line,
                e.range.end.character
              );
              return vscode.TextEdit.replace(range, e.newText);
            });
          }
        } catch (e) {}
        return [];
      }
    }
  );
  context.subscriptions.push(formattingProvider);

  // 4. Comandos de Terminal
  context.subscriptions.push(
    vscode.commands.registerCommand('vamos.runCurrentFile', () => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) return;
      const term = vscode.window.activeTerminal || vscode.window.createTerminal('VAMOS');
      term.show();
      term.sendText(`vamos run "${editor.document.fileName}"`);
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('vamos.buildCurrentFile', () => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) return;
      const term = vscode.window.activeTerminal || vscode.window.createTerminal('VAMOS');
      term.show();
      term.sendText(`vamos build "${editor.document.fileName}"`);
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('vamos.testWorkspace', () => {
      const term = vscode.window.activeTerminal || vscode.window.createTerminal('VAMOS');
      term.show();
      term.sendText(`vamos test`);
    })
  );
}

function deactivate() {
  if (lspProcess) {
    lspProcess.kill();
    lspProcess = null;
  }
}

module.exports = {
  activate,
  deactivate
};
