package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Task struct {
	Name string
	Cmd  string
}

func Parse(src string) ([]Task, error) {
	src = strings.TrimPrefix(src, "\uFEFF")
	var tasks []Task
	seen := false
	for i, line := range strings.Split(src, "\n") {
		line = strings.TrimRight(line, "\r")
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		if strings.HasPrefix(trim, "[") {
			if trim != "[tasks]" {
				return nil, fmt.Errorf("line %d: unknown section %s", i+1, trim)
			}
			if seen {
				return nil, fmt.Errorf("line %d: duplicate [tasks]", i+1)
			}
			seen = true
			continue
		}
		if !seen {
			return nil, fmt.Errorf("line %d: expected [tasks]", i+1)
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			return nil, fmt.Errorf("line %d: expected name = command", i+1)
		}
		name := strings.TrimSpace(line[:eq])
		command, err := parseCommand(strings.TrimSpace(line[eq+1:]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		if name == "" || strings.Contains(name, "\"") || !utf8.ValidString(name) {
			return nil, fmt.Errorf("line %d: invalid task name", i+1)
		}
		if strings.TrimSpace(command) == "" {
			return nil, fmt.Errorf("line %d: empty command", i+1)
		}
		tasks = append(tasks, Task{Name: name, Cmd: command})
	}
	if !seen {
		return nil, fmt.Errorf("missing [tasks]")
	}
	return tasks, nil
}

func parseCommand(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("empty command")
	}
	if s[0] != '"' {
		return s, nil
	}
	var b strings.Builder
	escaped := false
	closed := false
	for i := 1; i < len(s); i++ {
		c := s[i]
		if escaped {
			switch c {
			case '\\', '"':
				b.WriteByte(c)
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			default:
				return "", fmt.Errorf("bad escape \\%c", c)
			}
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '"' {
			rest := strings.TrimSpace(s[i+1:])
			if rest != "" {
				return "", fmt.Errorf("trailing characters after quoted command")
			}
			closed = true
			break
		}
		b.WriteByte(c)
	}
	if escaped || !closed {
		return "", fmt.Errorf("unclosed quote")
	}
	return b.String(), nil
}
