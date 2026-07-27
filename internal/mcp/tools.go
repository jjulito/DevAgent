package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ReadFileArgs struct {
	Path string `json:"path" jsonschema:"description=Absolute or relative path to the file to read"`
}

type ListDirArgs struct {
	Path string `json:"path" jsonschema:"description=Directory path to list contents of"`
}

type SearchFilesArgs struct {
	Pattern string `json:"pattern" jsonschema:"description=Glob pattern or filename to search for"`
	Path    string `json:"path,omitempty" jsonschema:"description=Root directory to search in (defaults to current dir)"`
}

type GrepArgs struct {
	Query string `json:"query" jsonschema:"description=Text pattern to search for in files"`
	Path  string `json:"path,omitempty" jsonschema:"description=Root directory to search in"`
}

func HandleReadFile(_ context.Context, args ReadFileArgs) (string, error) {
	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", args.Path, err)
	}

	if len(data) > 100*1024 {
		return string(data[:100*1024]) + "\n\n... [truncated at 100KB]", nil
	}

	return string(data), nil
}

func HandleListDir(_ context.Context, args ListDirArgs) (string, error) {
	entries, err := os.ReadDir(args.Path)
	if err != nil {
		return "", fmt.Errorf("failed to list %s: %w", args.Path, err)
	}

	var lines []string
	for _, entry := range entries {
		prefix := "📄"
		if entry.IsDir() {
			prefix = "📁"
		}
		info, _ := entry.Info()
		size := ""
		if info != nil && !entry.IsDir() {
			size = fmt.Sprintf(" (%d bytes)", info.Size())
		}
		lines = append(lines, fmt.Sprintf("%s %s%s", prefix, entry.Name(), size))
	}

	return strings.Join(lines, "\n"), nil
}

func HandleSearchFiles(_ context.Context, args SearchFilesArgs) (string, error) {
	root := args.Path
	if root == "" {
		root = "."
	}

	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		matched, _ := filepath.Match(args.Pattern, info.Name())
		if matched {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "No files found matching: " + args.Pattern, nil
	}

	return strings.Join(matches, "\n"), nil
}

func HandleGrep(_ context.Context, args GrepArgs) (string, error) {
	root := args.Path
	if root == "" {
		root = "."
	}

	var results []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Size() > 512*1024 {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.Contains(line, args.Query) {
				results = append(results, fmt.Sprintf("%s:%d: %s", path, i+1, strings.TrimSpace(line)))
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "No matches found for: " + args.Query, nil
	}

	if len(results) > 50 {
		results = results[:50]
		results = append(results, "... [truncated, showing first 50 results]")
	}

	return strings.Join(results, "\n"), nil
}
