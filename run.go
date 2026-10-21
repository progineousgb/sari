package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"
)

type result struct {
	name string
	code int
	ms   int64
}

func runFile(path string, keep bool, stdout, stderr io.Writer) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "sari: %s\n", err)
		return 2
	}
	tasks, err := Parse(string(data))
	if err != nil {
		fmt.Fprintf(stderr, "sari: %s\n", err)
		return 2
	}
	var rows []result
	final := 0
	for _, task := range tasks {
		fmt.Fprintf(stdout, "$ %s\n", task.Name)
		start := time.Now()
		cmd := commandFromLine(task.Cmd)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err = cmd.Run()
		code := statusOf(err, stderr)
		rows = append(rows, result{name: task.Name, code: code, ms: time.Since(start).Milliseconds()})
		if code != 0 && final == 0 {
			final = code
		}
		if code != 0 && !keep {
			break
		}
	}
	fmt.Fprintln(stdout)
	writeTable(stdout, rows)
	return final
}

func timeCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "sari: missing command")
		return 2
	}
	start := time.Now()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	code := statusOf(err, stderr)
	fmt.Fprintf(stdout, "exit %d\nduration %dms\n", code, time.Since(start).Milliseconds())
	return code
}

func statusOf(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	fmt.Fprintf(stderr, "sari: %s\n", err)
	return 1
}

func writeTable(w io.Writer, rows []result) {
	width := utf8.RuneCountInString("task")
	for _, row := range rows {
		if n := utf8.RuneCountInString(row.name); n > width {
			width = n
		}
	}
	fmt.Fprintf(w, "%s  exit  ms\n", pad("task", width))
	for _, row := range rows {
		fmt.Fprintf(w, "%s  %4d  %d\n", pad(row.name, width), row.code, row.ms)
	}
}

func pad(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(" ", width-n)
}
