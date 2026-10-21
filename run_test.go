package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// TestHelperProcess is not a real test. The other tests run the test binary
// again with GO_WANT_HELPER_PROCESS=1 and use it as a small portable command,
// so they do not need python or any other tool on PATH.
//
//	print <text>  writes text and a newline to stdout, exits 0
//	exit <code>   exits with code
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "helper: want -- <print|exit> <arg>")
		os.Exit(2)
	}
	switch args[1] {
	case "print":
		fmt.Println(args[2])
		os.Exit(0)
	case "exit":
		code, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "helper:", err)
			os.Exit(2)
		}
		os.Exit(code)
	}
	fmt.Fprintln(os.Stderr, "helper: unknown action", args[1])
	os.Exit(2)
}

// helperArgs returns the argv that runs TestHelperProcess with the given
// action, and sets the env var the helper checks for.
func helperArgs(t *testing.T, action ...string) []string {
	t.Helper()
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	exe, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	return append([]string{exe, "-test.run=TestHelperProcess", "--"}, action...)
}

// helperTask returns a sari.toml task value that runs the helper through
// sh -c or cmd /c. The value is written as a quoted TOML string so the
// backslashes and quotes in a Windows path survive Parse.
func helperTask(t *testing.T, action ...string) string {
	t.Helper()
	args := helperArgs(t, action...)
	var line string
	if runtime.GOOS == "windows" {
		// cmd /s /c strips the first and last quote of the line, so wrap
		// the whole thing once more to keep the quotes around the path.
		line = `"` + `"` + args[0] + `" ` + strings.Join(args[1:], " ") + `"`
	} else {
		// Single quotes keep sh from touching anything in the path.
		line = "'" + strings.ReplaceAll(args[0], "'", `'\''`) + "' " + strings.Join(args[1:], " ")
	}
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(line) + `"`
}

func writeTasks(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "sari.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunSuccess(t *testing.T) {
	dir := t.TempDir()
	body := "[tasks]\nprint = " + helperTask(t, "print", "123") + "\n"
	writeTasks(t, dir, body)
	t.Chdir(dir)
	var out bytes.Buffer
	code := runFile("sari.toml", false, &out, &out)
	text := out.String()
	if code != 0 {
		t.Fatalf("code %d\n%s", code, text)
	}
	if !strings.Contains(text, "123") {
		t.Fatalf("missing output\n%s", text)
	}
	if !strings.Contains(text, "print") || !strings.Contains(text, "exit") {
		t.Fatalf("missing table\n%s", text)
	}
}

func TestRunStopsOnFailure(t *testing.T) {
	dir := t.TempDir()
	body := "[tasks]\n" +
		"first = " + helperTask(t, "print", "1") + "\n" +
		"second = " + helperTask(t, "exit", "3") + "\n" +
		"third = " + helperTask(t, "print", "999") + "\n"
	writeTasks(t, dir, body)
	t.Chdir(dir)
	var out bytes.Buffer
	code := runFile("sari.toml", false, &out, &out)
	text := out.String()
	if code != 3 {
		t.Fatalf("code %d\n%s", code, text)
	}
	if strings.Contains(text, "999") {
		t.Fatalf("later task ran\n%s", text)
	}
	if !strings.Contains(text, "first") || !strings.Contains(text, "second") {
		t.Fatalf("table\n%s", text)
	}
}

func TestRunKeepGoing(t *testing.T) {
	dir := t.TempDir()
	body := "[tasks]\n" +
		"second = " + helperTask(t, "exit", "3") + "\n" +
		"third = " + helperTask(t, "print", "999") + "\n"
	writeTasks(t, dir, body)
	t.Chdir(dir)
	var out bytes.Buffer
	code := runFile("sari.toml", true, &out, &out)
	text := out.String()
	if code != 3 {
		t.Fatalf("code %d\n%s", code, text)
	}
	if !strings.Contains(text, "999") || !strings.Contains(text, "third") {
		t.Fatalf("keep-going\n%s", text)
	}
}

func TestRunPrintsUnicodeName(t *testing.T) {
	dir := t.TempDir()
	body := "[tasks]\nنبض = " + helperTask(t, "print", "7") + "\n"
	writeTasks(t, dir, body)
	t.Chdir(dir)
	var out bytes.Buffer
	code := runFile("sari.toml", false, &out, &out)
	text := out.String()
	if code != 0 || !strings.Contains(text, "نبض") || !strings.Contains(text, "7") {
		t.Fatalf("code %d\n%s", code, text)
	}
}

func TestTimeCommand(t *testing.T) {
	var out bytes.Buffer
	code := timeCommand(helperArgs(t, "print", "123"), &out, &out)
	text := out.String()
	if code != 0 || !strings.Contains(text, "123") || !strings.Contains(text, "exit 0") || !strings.Contains(text, "duration ") {
		t.Fatalf("code %d\n%s", code, text)
	}
}

func TestTimeMissingCommand(t *testing.T) {
	var out bytes.Buffer
	code := timeCommand([]string{"sari-missing-command-xyz"}, &out, &out)
	text := out.String()
	if code == 0 || !strings.Contains(text, "sari:") {
		t.Fatalf("code %d\n%s", code, text)
	}
}

func TestTimeCommandExit(t *testing.T) {
	var out bytes.Buffer
	code := timeCommand(helperArgs(t, "exit", "4"), &out, &out)
	if code != 4 || !strings.Contains(out.String(), "exit 4") {
		t.Fatalf("code %d\n%s", code, out.String())
	}
}
