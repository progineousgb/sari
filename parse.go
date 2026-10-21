package main

import (
	"fmt"
	"strconv"
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
		name, err := parseKey(strings.TrimSpace(line[:eq]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		command, err := parseCommand(strings.TrimSpace(line[eq+1:]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		if name == "" || !utf8.ValidString(name) {
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

// parseKey accepts a bare key, a "basic" quoted key, or a 'literal' key.
func parseKey(s string) (string, error) {
	if s == "" || (s[0] != '"' && s[0] != '\'') {
		if strings.ContainsAny(s, "\"'") {
			return "", fmt.Errorf("invalid task name")
		}
		return s, nil
	}
	v, rest, err := parseQuoted(s)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(rest) != "" {
		return "", fmt.Errorf("invalid task name")
	}
	return v, nil
}

// parseCommand reads a value. A value in double quotes ("basic") or single
// quotes ('literal') has its quotes removed; a "#" comment may follow it.
// An unquoted value is taken as written.
func parseCommand(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("empty command")
	}
	if s[0] != '"' && s[0] != '\'' {
		return s, nil
	}
	v, rest, err := parseQuoted(s)
	if err != nil {
		return "", err
	}
	rest = strings.TrimSpace(rest)
	if rest != "" && !strings.HasPrefix(rest, "#") {
		return "", fmt.Errorf("trailing characters after quoted command")
	}
	return v, nil
}

// parseQuoted reads one quoted string at the start of s and returns its
// value and the text after the closing quote.
func parseQuoted(s string) (string, string, error) {
	quote := s[0]
	var b strings.Builder
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c == quote {
			return b.String(), s[i+1:], nil
		}
		if c != '\\' || quote == '\'' {
			b.WriteByte(c)
			continue
		}
		i++
		if i >= len(s) {
			break
		}
		switch s[i] {
		case '\\', '"':
			b.WriteByte(s[i])
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case 'u', 'U':
			n := 4
			if s[i] == 'U' {
				n = 8
			}
			if i+n >= len(s) {
				return "", "", fmt.Errorf("bad escape \\%c", s[i])
			}
			code, err := strconv.ParseUint(s[i+1:i+1+n], 16, 32)
			if err != nil || !utf8.ValidRune(rune(code)) {
				return "", "", fmt.Errorf("bad escape \\%c", s[i])
			}
			b.WriteRune(rune(code))
			i += n
		default:
			return "", "", fmt.Errorf("bad escape \\%c", s[i])
		}
	}
	return "", "", fmt.Errorf("unclosed quote")
}
