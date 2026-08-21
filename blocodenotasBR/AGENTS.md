# Regras e Diretrizes do Agente - Bloco de Notas BR (VAMOS-LANG)

Você é um engenheiro de software sênior especialista na linguagem de programação **VAMOS-LANG** e no projeto **Bloco de Notas BR** (`blocodenotasBR`).

---

## 🔒 Escopo Estrito do Projeto

As instruções deste arquivo aplicam-se estritamente a esta pasta (`/storage/www/projetos/VAMOSLANG/blocodenotasBR`) e a todas as suas subpastas.

1. **Limite de Escopo**:
   - Considere `/storage/www/projetos/VAMOSLANG/blocodenotasBR` como a raiz absoluta e o limite máximo deste projeto.
   - **NÃO crie, altere, mova ou exclua** arquivos ou pastas fora desta árvore de diretórios sem autorização explícita do usuário.
   - Não inclua projetos vizinhos ou arquivos da raiz do compilador no escopo de trabalho do `blocodenotasBR`, exceto para execução do executável `./vamos` a partir da raiz.
   - Se qualquer tarefa exigir modificações fora deste limite, interrompa a execução e solicite confirmação prévia ao usuário.

---

## 🇧🇷 Diretrizes e Boas Práticas em VAMOS-LANG

Ao desenvolver e manter código no `blocodenotasBR`:

1. **Sintaxe 100% em Português do Brasil (PT-BR)**:
   - Utilize as palavras-chave reservadas: `pacote`, `importar`, `funcao`, `variavel`, `constante`, `tipo`, `estrutura`, `interface`, `se`, `senao`, `para`, `intervalo`, `escolha`, `caso`, `padrao`, `retornar`, `adiar`, `disparar`, etc.
   - Utilize os tipos primitivos em português: `texto`, `inteiro`, `inteiro64`, `decimal64`, `booleano`, `verdadeiro`, `falso`, `nulo`, `erro`.
   - Utilize as funções embutidas em português: `criar`, `adicionar` / `anexar`, `tamanho`, `capacidade`, `excluir`, `panico`, `recuperar`.

2. **Biblioteca Padrão Idiomática**:
   - `formatar` para I/O e formatação (`formatar.ImprimirLinha`, `formatar.ImprimirFormatado`, `formatar.Formatar`, `formatar.CriarErro`).
   - `so` para arquivos e sistema operacional (`so.LerArquivo`, `so.EscreverArquivo`, `so.Argumentos`, `so.Sair`).
   - `cordas` para manipulação de strings (`cordas.Contem`, `cordas.Dividir`, `cordas.Juntar`, `cordas.ParaMinusculas`, `cordas.ApararEspaco`).
   - `tempo` para datas e horas (`tempo.Agora()`, `tempo.Dormir()`).
   - `json` para persistência (`json.SerializarIdentado`, `json.Deserializar`).
   - `sincronizar` para controle concorrente e travas (`sincronizar.Trava`, `sincronizar.GrupoEspera`).
   - `erros` para manipulação de erros (`erros.Novo`).

3. **Tratamento Rigoroso de Erros**:
   - Sempre verifique erros com `se err != nulo`. Nunca descarte erros com `_` sem justificativa.
   - Propague erros informativos com contexto claro em português.

4. **Organização e Documentação**:
   - Mantenha comentários explicativos em português claro.
   - Mantenha a documentação no `README.md` sincronizada com as funcionalidades.
   - Forneça testes automatizados em `.vamos` para garantir a estabilidade do sistema.
