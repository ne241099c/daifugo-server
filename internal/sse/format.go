package sse

import (
	"encoding/json"
	"fmt"
	"io"
)

// writeSSE は SSE フォーマットで1イベントを書き込む
func writeSSE(w io.Writer, ev Event) error {
	// retry（任意）
	if ev.Retry != nil {
		if _, err := fmt.Fprintf(w, "retry: %d\n", *ev.Retry); err != nil {
			return err
		}
	}
	// id
	if _, err := fmt.Fprintf(w, "id: %d\n", ev.ID); err != nil {
		return err
	}
	// event（任意）
	if ev.Type != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", ev.Type); err != nil {
			return err
		}
	}

	// data（複数行対応）
	b, err := json.Marshal(ev.Data)
	if err != nil {
		return err
	}
	for _, line := range splitLines(string(b)) {
		if _, err := fmt.Fprintf(w, "data: %s\n", line); err != nil {
			return err
		}
	}

	// イベント終端（空行）
	_, err = fmt.Fprint(w, "\n")
	return err
}

func splitLines(s string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
