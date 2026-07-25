package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func parseSSEStream(client *http.Client, req *http.Request, textCh chan<- string, errCh chan<- error) {
	resp, err := client.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("request failed: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		errCh <- fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		return
	}

	buf := make([]byte, 4096)
	var remainder string

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			lines := strings.Split(remainder+string(buf[:n]), "\n")
			remainder = lines[len(lines)-1]

			for _, line := range lines[:len(lines)-1] {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "data: ") {
					continue
				}
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					return
				}

				var chunk orStreamChunk
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					continue
				}
				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
					textCh <- chunk.Choices[0].Delta.Content
				}
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				errCh <- readErr
			}
			return
		}
	}
}
