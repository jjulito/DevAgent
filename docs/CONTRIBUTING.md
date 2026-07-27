# Contributing to DevAgent CLI

## Getting Started

1. Fork and clone the repository
2. Run the setup script:
   ```bash
   chmod +x scripts/setup.sh
   ./scripts/setup.sh
   ```
3. Create a branch for your changes

## Project Structure

```
cmd/              CLI command definitions (one file per command)
internal/
  config/         Configuration loading and validation
  provider/       LLM provider interface and adapters
  rag/            RAG pipeline (chunker, indexer, retriever)
  mcp/            MCP server and tool handlers
  git/            Git utilities
  review/         Code review logic
  output/         Terminal formatting
docs/             Documentation
```

## Adding a New Command

Create a file in `cmd/`:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // your logic here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Adding a New LLM Provider

1. Create `internal/provider/myprovider.go`
2. Implement the `LLMProvider` interface:
   ```go
   type LLMProvider interface {
       Chat(ctx context.Context, messages []Message) (*Response, error)
       ChatStream(ctx context.Context, messages []Message) (<-chan string, <-chan error)
       Name() string
   }
   ```
3. Register it in `internal/provider/factory.go`
4. Add config fields in `internal/config/config.go`

## Code Style

- Follow standard Go conventions (`go vet`, `go fmt`)
- Keep files focused — one responsibility per file
- Use the shared `parseSSEStream` helper for OpenAI-compatible APIs
- Avoid unnecessary comments — code should be self-documenting

## Testing

```bash
make test       # run all tests
make lint       # run go vet
make build      # verify it compiles
```

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(provider): add Anthropic adapter
fix(review): handle empty git diff gracefully
docs: update provider configuration guide
chore: update dependencies
```
