package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// OllamaRequest represents the request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// GenerateQA sends a prompt to the Ollama API and returns the response
func GenerateQA(model, url, prompt string) (string, error) {
	body, _ := json.Marshal(OllamaRequest{Model: model, Prompt: prompt})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var fullResponse string
	dec := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := dec.Decode(&chunk); err != nil {
			break
		}
		if r, ok := chunk["response"].(string); ok {
			fullResponse += r
		}
		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}
	return fullResponse, nil
}

// ParseFlashcards parses the Ollama response and returns a list of questions and answers
func ParseFlashcards(response string) ([]map[string]string, error) {
	var qas []map[string]string
	lines := splitLines(response)
	var q, a string
	for _, line := range lines {
		l := trimSpace(line)
		if len(l) > 1 && (l[:2] == "Q:" || l[:2] == "A:") {
			if l[:2] == "Q:" {
				q = trimSpace(l[2:])
			} else if l[:2] == "A:" {
				a = trimSpace(l[2:])
				if q != "" && a != "" {
					qas = append(qas, map[string]string{"question": q, "answer": a})
					q, a = "", ""
				}
			}
		}
	}
	return qas, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == '*' || s[0] == ' ' || s[0] == ':') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == '*' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}
