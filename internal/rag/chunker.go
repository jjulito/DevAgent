package rag

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ChunkOptions struct {
	MaxLines  int
	Overlap   int
	RootPath  string
}

var defaultExtensions = map[string]string{
	".go":   "go",
	".py":   "python",
	".js":   "javascript",
	".ts":   "typescript",
	".jsx":  "javascript",
	".tsx":  "typescript",
	".rs":   "rust",
	".java": "java",
	".c":    "c",
	".cpp":  "cpp",
	".h":    "c",
	".cs":   "csharp",
	".rb":   "ruby",
	".php":  "php",
	".sh":   "bash",
	".md":   "markdown",
	".yaml": "yaml",
	".yml":  "yaml",
	".json": "json",
	".toml": "toml",
	".sql":  "sql",
	".html": "html",
	".css":  "css",
}

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"__pycache__":  true,
	".venv":        true,
	"dist":         true,
	"build":        true,
	"target":       true,
	".next":        true,
}

func ChunkDirectory(opts ChunkOptions) ([]Document, error) {
	if opts.MaxLines == 0 {
		opts.MaxLines = 60
	}
	if opts.Overlap == 0 {
		opts.Overlap = 10
	}

	var docs []Document

	err := filepath.Walk(opts.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		lang, ok := defaultExtensions[ext]
		if !ok {
			return nil
		}

		if info.Size() > 512*1024 {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(opts.RootPath, path)
		chunks := splitLines(string(content), opts.MaxLines, opts.Overlap)

		for i, chunk := range chunks {
			id := chunkID(relPath, i)
			docs = append(docs, Document{
				ID:       id,
				Content:  fmt.Sprintf("// File: %s (chunk %d/%d)\n\n%s", relPath, i+1, len(chunks), chunk),
				FilePath: relPath,
				Language: lang,
				Chunk:    i,
			})
		}

		return nil
	})

	return docs, err
}

func splitLines(content string, maxLines, overlap int) []string {
	lines := strings.Split(content, "\n")

	if len(lines) <= maxLines {
		return []string{content}
	}

	var chunks []string
	for i := 0; i < len(lines); i += maxLines - overlap {
		end := i + maxLines
		if end > len(lines) {
			end = len(lines)
		}
		chunk := strings.Join(lines[i:end], "\n")
		chunks = append(chunks, chunk)
		if end == len(lines) {
			break
		}
	}

	return chunks
}

func chunkID(path string, chunk int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", path, chunk)))
	return fmt.Sprintf("%x", h[:8])
}
