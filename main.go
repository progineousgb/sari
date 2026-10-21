package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	os.Exit(dispatch(os.Args[1:], os.Stdout, os.Stderr))
}

func dispatch(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runFile("sari.toml", false, stdout, stderr)
	}
	switch args[0] {
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	case "run":
		return dispatchRun(args[1:], stdout, stderr)
	case "time":
		return timeCommand(stripDash(args[1:]), stdout, stderr)
	default:
		if strings.HasPrefix(args[0], "-") {
			return dispatchRun(args, stdout, stderr)
		}
		fmt.Fprintf(stderr, "sari: unknown command %s\n", args[0])
		usage(stderr)
		return 2
	}
}

func dispatchRun(args []string, stdout, stderr io.Writer) int {
	path, keep, err := parseRunFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "sari: %s\n", err)
		usage(stderr)
		return 2
	}
	return runFile(path, keep, stdout, stderr)
}

func stripDash(args []string) []string {
	if len(args) > 0 && args[0] == "--" {
		return args[1:]
	}
	return args
}

func parseRunFlags(args []string) (string, bool, error) {
	path := "sari.toml"
	keep := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--keep-going":
			keep = true
		case "-f", "--file":
			i++
			if i >= len(args) {
				return "", false, fmt.Errorf("missing file after %s", args[i-1])
			}
			path = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", false, fmt.Errorf("unknown flag %s", args[i])
			}
			return "", false, fmt.Errorf("unexpected argument %s", args[i])
		}
	}
	return path, keep, nil
}

func usage(w io.Writer) {
	fmt.Fprint(w, `sari run [--keep-going] [-f sari.toml]
sari time -- <command> [args...]
`)
}
