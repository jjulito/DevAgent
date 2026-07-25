package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type DiffOptions struct {
	Staged bool
	Path   string
}

type DiffResult struct {
	Raw       string
	FileCount int
	Files     []string
}

func GetDiff(opts DiffOptions) (*DiffResult, error) {
	args := []string{"diff", "--no-color"}

	if opts.Staged {
		args = append(args, "--cached")
	}

	if opts.Path != "" {
		args = append(args, "--", opts.Path)
	}

	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w\n%s", err, string(out))
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, fmt.Errorf("no changes found (try --staged for staged changes)")
	}

	files := parseChangedFiles(raw)

	return &DiffResult{
		Raw:       raw,
		FileCount: len(files),
		Files:     files,
	}, nil
}

func GetLog(n int) (string, error) {
	args := []string{"log", fmt.Sprintf("-n%d", n), "--oneline", "--no-color"}

	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w\n%s", err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}

func parseChangedFiles(diff string) []string {
	var files []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				file := strings.TrimPrefix(parts[3], "b/")
				files = append(files, file)
			}
		}
	}
	return files
}
