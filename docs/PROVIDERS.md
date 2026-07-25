# Providers — DevAgent CLI

## Supported Providers

DevAgent CLI supports 4 LLM providers, all using the same `LLMProvider` interface:

| Provider | API Key Required | Cost | Best For |
|----------|-----------------|------|----------|
| OpenRouter | Yes | Free models available | Default, zero-cost usage |
| Ollama | No | Free (local) | Privacy, offline usage |
| OpenAI | Yes | Paid | GPT-4o, best quality |
| Gemini | Yes | Free tier available | Gemini 2.5 Flash/Pro |

---

## Configuration

### Via `.env` file

```bash
cp .env.example .env
# Edit .env with your keys
```

### Via environment variables

```bash
export DEVAGENT_PROVIDER=openrouter
export OPENROUTER_API_KEY=sk-or-v1-xxxxx
```

### Via CLI flags

```bash
devagent ask "question" -p ollama -m llama3.2
```

**Priority**: CLI flags > environment variables > .env file > defaults

---

## OpenRouter (Default)

Free models available. Get your API key at [openrouter.ai](https://openrouter.ai).

```bash
DEVAGENT_PROVIDER=openrouter
OPENROUTER_API_KEY=sk-or-v1-xxxxxxxxxxxx
DEVAGENT_MODEL=meta-llama/llama-4-scout:free
```

Free models include:
- `meta-llama/llama-4-scout:free`
- `google/gemma-3-27b-it:free`
- `mistralai/mistral-small-3.2-24b-instruct:free`

---

## Ollama (Local)

No API key needed. Install from [ollama.com](https://ollama.com).

```bash
# Install and run a model
ollama pull llama3.2
ollama serve

# Configure DevAgent
DEVAGENT_PROVIDER=ollama
OLLAMA_HOST=http://localhost:11434
DEVAGENT_MODEL=llama3.2
```

---

## OpenAI

Get your API key at [platform.openai.com](https://platform.openai.com).

```bash
DEVAGENT_PROVIDER=openai
OPENAI_API_KEY=sk-xxxxxxxxxxxx
DEVAGENT_MODEL=gpt-4o-mini
```

---

## Gemini

Get your API key at [aistudio.google.com](https://aistudio.google.com).

```bash
DEVAGENT_PROVIDER=gemini
GEMINI_API_KEY=AIzaxxxxxxxxxxxxxxxx
DEVAGENT_MODEL=gemini-2.5-flash
```

---

## Switching Providers On-the-fly

```bash
# Use OpenRouter (default)
devagent ask "explain goroutines"

# Use Ollama for this question
devagent ask "explain goroutines" -p ollama

# Use a specific model
devagent review -p openai -m gpt-4o
```

## Check Configuration

```bash
devagent config
```
