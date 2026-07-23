# Commands — DevAgent CLI

## General Usage

```bash
devagent [command] [arguments] [flags]
```

### Global Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--provider` | `-p` | LLM Provider | `openrouter` |
| `--model` | `-m` | Specific model | according to provider |
| `--verbose` | `-v` | Detailed output | `false` |
| `--config` | `-c` | Configuration file path | `.env` |

---

## `devagent ask`

Sends a question to the configured LLM and shows the response in streaming.

```bash
devagent ask "what is a goroutine?"
devagent ask "explain this error" -p ollama -m llama3.2
```

**Arguments:**
- `question` (required): The question to send to the model.

---

## `devagent review` *(Phase 2)*

Generates an automatic code review of the current Git diff.

```bash
devagent review              # review of unstaged diff
devagent review --staged     # review of staged changes
```

---

## `devagent index` *(Phase 3)*

Indexes the project source code in the vector store for semantic search.

```bash
devagent index .             # index current directory
devagent index ./src         # index sub-directory
```

---

## `devagent search` *(Phase 3)*

Performs semantic search in the indexed code.

```bash
devagent search "error handling"
devagent search "authentication" --top 10
```

---

## `devagent config` *(Phase 2)*

Shows or edits the active configuration.

```bash
devagent config --show       # show current config
```

---

## `devagent serve` *(Phase 4)*

Runs an MCP server that exposes project tools to AI agents.

```bash
devagent serve               # default port (3000)
devagent serve --port 8080   # custom port
```
