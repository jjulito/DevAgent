# RAG Pipeline — DevAgent CLI

## Overview

DevAgent CLI includes a **Retrieval-Augmented Generation (RAG)** pipeline that lets you index your project's source code and perform semantic search over it.

## Architecture

```
Source Code → Chunker → Vector Store (Qdrant) → Semantic Search → LLM Context
```

1. **Chunker** (`internal/rag/chunker.go`): Walks the project directory, filters by file extension, and splits files into overlapping chunks (60 lines, 10 line overlap).
2. **Indexer** (`internal/rag/indexer.go`): Orchestrates chunking and storage in the vector engine.
3. **VectorEngine** (`internal/rag/engine.go`): Interface for vector store backends. Current implementation: Qdrant.
4. **Retriever** (`internal/rag/retriever.go`): Performs semantic search and optionally feeds results as context to the LLM.

## Quick Start

### 1. Start Qdrant

```bash
docker compose up -d
```

### 2. Index Your Project

```bash
devagent index .           # index current directory
devagent index ./src       # index specific directory
```

### 3. Search

```bash
devagent search "error handling"
devagent search "authentication" --top 10
```

## Supported File Types

| Extension | Language |
|-----------|----------|
| `.go` | Go |
| `.py` | Python |
| `.js`, `.jsx` | JavaScript |
| `.ts`, `.tsx` | TypeScript |
| `.rs` | Rust |
| `.java` | Java |
| `.c`, `.h`, `.cpp` | C/C++ |
| `.cs` | C# |
| `.rb` | Ruby |
| `.php` | PHP |
| `.sh` | Bash |
| `.md` | Markdown |
| `.yaml`, `.yml` | YAML |
| `.json` | JSON |
| `.sql` | SQL |
| `.html`, `.css` | HTML/CSS |

## Chunking Strategy

- **Max lines per chunk**: 60
- **Overlap**: 10 lines (ensures context continuity)
- **Skip directories**: `.git`, `node_modules`, `vendor`, `__pycache__`, `dist`, `build`, `target`
- **Max file size**: 512KB (larger files are skipped)

## Qdrant Configuration

Default connection in `.env`:

```bash
QDRANT_HOST=localhost
QDRANT_PORT=6334
```

Collection name: `devagent` (created automatically on first index).

## Notes

> The current implementation uses a simple hash-based vector as a placeholder. For production-quality semantic search, integrate a real embedding model (e.g., via OpenAI embeddings API or a local model like `all-MiniLM-L6-v2`).
