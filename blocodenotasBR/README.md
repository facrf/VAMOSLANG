# 📝 Bloco de Notas BR (`blocodenotasBR`)

> **Gerenciador de notas moderno, robusto e de alta performance desenvolvido na linguagem VAMOS-LANG (Go em Português do Brasil).**

[![VAMOS-LANG](https://img.shields.io/badge/Linguagem-VAMOS--LANG-00ADD8?style=flat&logo=go)](../README.md)
[![Status](https://img.shields.io/badge/Status-Concluído-brightgreen.svg)]()
[![Licença](https://img.shields.io/badge/Licença-MIT-green.svg)](../LICENSE)

---

## 📖 Visão Geral

O **Bloco de Notas BR** é uma aplicação completa de gerenciamento de anotações pessoais, estudos e trabalho desenvolvida com código-fonte 100% em sintaxe **VAMOS-LANG**.

O projeto suporta tanto uso direto via linha de comando (**CLI**) quanto um modo de terminal iterativo (**Console REPL**), com persistência automática de dados em formato JSON, sistema de busca textual, marcação de favoritos, cálculo de métricas e exportação de relatórios em texto.

---

## ✨ Funcionalidades

- ✍️ **Criação Rápida de Notas**: Adicione notas com título, conteúdo e categorização personalizada (ex: *Trabalho*, *Estudos*, *Pessoal*).
- 📋 **Listagem Formatada em Tabela**: Visualização tabular com IDs, indicador de favoritos (`★`), categoria e datas de criação.
- 🔍 **Busca Textual Integrada**: Pesquise anotações por palavras-chave presentes no título, conteúdo ou categoria.
- 🏷️ **Categorização e Edição**: Altere conteúdo e categoria de notas existentes com atualização automática de timestamps.
- ⭐ **Sistema de Favoritos**: Marque e desmarque anotações prioritárias com um único comando.
- 📊 **Estatísticas e Métricas**: Resumo com total de notas, contagem estimada de palavras e distribuição por categorias.
- 💾 **Persistência Confiável em JSON**: Gravação e leitura automática em arquivo JSON estruturado (`notas_db.json`).
- 📤 **Exportação em Texto**: Exporte todas as anotações organizadas para um arquivo `.txt` formatado.
- 💻 **Modo Interativo (REPL)**: Menu interativo amigável para uso contínuo no terminal.
- 🧪 **Suíte de Testes Automatizados**: Testes unitários e de integração escritos na própria linguagem VAMOS-LANG.

---

## 📂 Estrutura do Projeto

```
blocodenotasBR/
├── AGENTS.md            # Regras de isolamento e escopo para LLMs e agentes de IA
├── .gitignore           # Ignora binários compilados, temporários e logs
├── README.md            # Documentação completa e guia de uso
├── main.vamos           # Código-fonte da aplicação Bloco de Notas BR
├── testes.vamos         # Suíte de testes automatizados do sistema
└── notas_db.json        # Banco de dados local gerado automaticamente
```

---

## 🚀 Como Executar e Compilar

Os comandos a seguir devem ser executados a partir da raiz do repositório (`/storage/www/projetos/VAMOSLANG`):

### 1. Executando Diretamente com o Transpilador (`rodar` / `run`)

```bash
# Exibir tela de ajuda com todos os comandos
./vamos rodar blocodenotasBR/main.vamos ajuda

# Criar uma nova nota
./vamos rodar blocodenotasBR/main.vamos criar "Reunião de Arquitetura" "Definir módulos do sistema" Trabalho

# Listar todas as notas
./vamos rodar blocodenotasBR/main.vamos listar

# Iniciar o modo interativo (REPL)
./vamos rodar blocodenotasBR/main.vamos interativo
```

### 2. Compilando para Binário Executável Nativo (`compilar` / `build`)

```bash
# Compila o projeto gerando o binário nativo 'blocodenotas'
./vamos compilar blocodenotasBR/main.vamos -o blocodenotasBR/blocodenotas

# Executar o binário diretamente
./blocodenotasBR/blocodenotas listar
./blocodenotasBR/blocodenotas estatisticas
```

### 3. Executando a Suíte de Testes Automatizados (`testes.vamos`)

```bash
./vamos rodar blocodenotasBR/testes.vamos
```

Saída esperada:
```
==================================================
  🧪 SUÍTE DE TESTES: BLOCO DE NOTAS BR (VAMOS)   
==================================================
[✓ TESTE 1] Inicialização de Bloco de Notas: OK
[✓ TESTE 2] Criação de notas com persistência: OK
[✓ TESTE 3] Validação de título obrigatório: OK
[✓ TESTE 4] Consulta por ID existente/inexistente: OK
[✓ TESTE 5] Alternância de status favorito: OK
[✓ TESTE 6] Busca textual em conteúdo: OK
[✓ TESTE 7] Recarga de dados persistidos em JSON: OK
[✓ TESTE 8] Exclusão de nota: OK
--------------------------------------------------
Resultado: 8 de 8 testes passaram com êxito!
Status: TODOS OS TESTES PASSARAM COM SUCESSO! ✓
==================================================
```

### 4. Formatando e Verificando o Código com Linter

```bash
# Formatar os arquivos .vamos do projeto
./vamos fmt blocodenotasBR/

# Executar análise estática de boas práticas
./vamos lint blocodenotasBR/
```

---

## 🛠️ Referência de Comandos da Linha de Comando (CLI)

| Comando | Argumentos | Descrição | Exemplo |
| :--- | :--- | :--- | :--- |
| `criar` | `<titulo> <conteudo> [categoria]` | Cria uma nova anotação | `./vamos rodar blocodenotasBR/main.vamos criar "Estudar VAMOS" "Revisar sintaxe" Estudos` |
| `listar` | `[filtro]` | Lista notas em tabela formatada | `./vamos rodar blocodenotasBR/main.vamos listar` |
| `ver` | `<id>` | Exibe detalhes completos da nota | `./vamos rodar blocodenotasBR/main.vamos ver 1` |
| `editar` | `<id> <novo_conteudo>` | Atualiza o conteúdo da nota | `./vamos rodar blocodenotasBR/main.vamos editar 1 "Novo texto"` |
| `categoria`| `<id> <nova_categoria>` | Altera a categoria da nota | `./vamos rodar blocodenotasBR/main.vamos categoria 1 Projetos` |
| `favoritar`| `<id>` | Alterna o status de favorita (`★`) | `./vamos rodar blocodenotasBR/main.vamos favoritar 1` |
| `buscar` | `<termo>` | Pesquisa em títulos, conteúdos e tags | `./vamos rodar blocodenotasBR/main.vamos buscar VAMOS` |
| `excluir` | `<id>` | Remove permanentemente uma nota | `./vamos rodar blocodenotasBR/main.vamos excluir 1` |
| `estatisticas`| — | Exibe métricas de notas, categorias e palavras | `./vamos rodar blocodenotasBR/main.vamos estatisticas` |
| `exportar` | `[arquivo.txt]` | Exporta relatório formatado | `./vamos rodar blocodenotasBR/main.vamos exportar relatorio.txt` |
| `interativo`| — | Inicia console interativo com prompt | `./vamos rodar blocodenotasBR/main.vamos interativo` |
| `ajuda` | — | Exibe as instruções de uso | `./vamos rodar blocodenotasBR/main.vamos ajuda` |

---

## 💻 Demonstração do Modo Interativo

Ao executar sem argumentos ou com o comando `interativo`:

```
====================================================================
       📝 BLOCO DE NOTAS BR - GERENCIADOR EM VAMOS-LANG 🇧🇷          
           Versão 1.0.0 | Alta Performance com Sintaxe PT-BR          
====================================================================
Modo Interativo Ativado! Digite 'ajuda' para comandos ou 'sair' para encerrar.
--------------------------------------------------------------------

bloco-notas-br> criar
Título da Nota: Aprender Concorrência
Categoria (ex: Pessoal, Trabalho, Estudos) [Geral]: Estudos
Conteúdo da Nota: Praticar o uso de 'disparar', 'canal' e 'sincronizar.GrupoEspera'.
✓ Nota #3 criada com sucesso com o título 'Aprender Concorrência'!

bloco-notas-br> listar
┌─────┬───┬──────────────────────────────┬──────────────┬─────────────────────┐
│ ID  │ ★ │ TÍTULO                       │ CATEGORIA    │ DATA DE CRIAÇÃO     │
├─────┼───┼──────────────────────────────┼──────────────┼─────────────────────┤
│ 1   │ ★ │ Minha Primeira Nota          │ Estudos      │ 21/08/2026 16:55:54 │
│ 2   │   │ Compilador VAMOS             │ Tecnologia   │ 21/08/2026 16:55:56 │
│ 3   │   │ Aprender Concorrência        │ Estudos      │ 21/08/2026 17:00:10 │
└─────┴───┴──────────────────────────────┴──────────────┴─────────────────────┘
Total listado: 3 nota(s)

bloco-notas-br> stats
📊 Estatísticas do Bloco de Notas:
  • Total de Notas:     3
  • Notas Favoritas:    1
  • Total de Palavras:  30
  • Categorias:
      - Estudos        : 2 nota(s)
      - Tecnologia     : 1 nota(s)

bloco-notas-br> sair
Encerrando Bloco de Notas BR. Até logo!
```

---

## 🇧🇷 Boas Práticas e Idiomas VAMOS-LANG Utilizados

O código-fonte do `Bloco de Notas BR` segue rigorosamente as convenções de engenharia de software em **VAMOS-LANG**:

1. **Tipagem e Estruturas Expressivas**:
   ```vamos
   tipo Nota estrutura {
       Id           inteiro  `json:"id"`
       Titulo       texto    `json:"titulo"`
       Conteudo     texto    `json:"conteudo"`
       Categoria    texto    `json:"categoria"`
       CriadoEm     texto    `json:"criado_em"`
       AtualizadoEm texto    `json:"atualizado_em"`
       Favorita     booleano `json:"favorita"`
   }
   ```

2. **Manipulação de Coleções e Fatias**:
   - Uso das funções nativas `criar`, `adicionar` e `tamanho`.
   - Iteração limpa com `para _, nota := intervalo b.Notas`.

3. **Biblioteca Padrão e Formatação**:
   - `formatar.ImprimirLinha`, `formatar.ImprimirFormatado`, `formatar.CriarErro`.
   - `so.LerArquivo` e `so.EscreverArquivo` para operações seguras de I/O em disco.
   - `json.SerializarIdentado` e `json.Deserializar` para serialização de dados.
   - `tempo.Agora().Formato(...)` para carimbos temporais precisos.
   - `cordas.Contem`, `cordas.ParaMinusculas`, `cordas.ApararEspaco` e `cordas.Dividir`.

4. **Tratamento Seguro de Erros**:
   - Verificação explícita de `se err != nulo` em todas as chamadas propensas a falha.
   - Sem descartes silenciosos de erros.

---

## 📄 Licença

Este projeto está sob a licença [MIT](../LICENSE).
