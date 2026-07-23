# Architecture — DevAgent CLI

## Overview

DevAgent CLI is a modular terminal tool for developers that follows the principles of **Clean Architecture** in Go. It allows you to analyze code, perform automated code reviews, conduct semantic search (RAG) on projects, and expose tools via MCP — all using your own API key (BYOK).

## Layer Diagram

```
┌──────────────────────────────────────────────┐
│                   cmd/                        │  ← CLI Layer
│  Flag parsing, logic invocation               │     Only I/O, no logic
├──────────────────────────────────────────────┤
│               internal/                       │
│  ┌────────────┐ ┌────────────┐ ┌───────────┐ │
│  │  provider/  │ │    rag/    │ │   mcp/    │ │  ← Domain Layer
│  │ LLMProvider │ │VectorEngine│ │ MCPServer │ │     Interfaces + logic
│  └─────┬──────┘ └─────┬──────┘ └─────┬─────┘ │
│        │              │              │        │
│  ┌─────▼──────┐ ┌─────▼──────┐      │        │
│  │ openrouter │ │  qdrant    │      │        │  ← Adapters
│  │ openai     │ │            │      │        │     Implementations
│  │ ollama     │ └────────────┘      │        │
│  │ gemini     │                     │        │
│  └────────────┘                     │        │
├──────────────────────────────────────────────┤
│              config/                          │  ← Configuration
│  .env, flags, Viper                           │     Centralized
├──────────────────────────────────────────────┤
│              output/                          │  ← Presentation
│  Colores, terminal formatting                 │     No business logic
└──────────────────────────────────────────────┘
```

## Principles

### Dependency Rule
Internal layers (`provider`, `rag`) define **interfaces**. External layers implement them. Never the other way around.

### One file, one responsibility
Each file has a clear purpose. One adapter per file, one command per file.

### Factory Pattern
Providers are instantiated through `provider.New(cfg)` — the consuming code does not know the concrete implementation.

### Cascade Configuration
`.env` → Environment variables → CLI flags. Each level overwrites the previous one (Viper).

## Technical Stack

| Component | Technology | Purpose |
|-----------|-----------|-----------|
| CLI Framework | Cobra + Viper | Subcommands, flags, configuration |
| LLM | OpenRouter, OpenAI, Gemini, Ollama | Interchangeable AI providers |
| Vector Store | Qdrant | Embeddings for RAG |
| MCP | go-sdk oficial | Expose tools to agents |
| Containerization | Docker Compose | Complete stack in one command |
