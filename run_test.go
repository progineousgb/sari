package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	dir := t.TempDir()
	body := "[tasks]\nprint = python -c \"print(123)\"\n"
	if err := os.WriteFile(dir+"/sari.toml", []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
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
		"first = python -c \"print(1)\"\n" +
		"second = python -c \"import sys; sys.exit(3)\"\n" +
		"third = python -c \"print(999)\"\n"
	if err := os.WriteFile(dir+"/sari.toml", []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
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
		"second = python -c \"import sys; sys.exit(3)\"\n" +
		"third = python -c \"print(999)\"\n"
	if err := os.WriteFile(dir+"/sari.toml", []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
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
	body := "[tasks]\nنبض = python -c \"print(7)\"\n"
	if err := os.WriteFile(dir+"/sari.toml", []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
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
	code := timeCommand([]string{"python", "-c", "print(123)"}, &out, &out)
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
	code := timeCommand([]string{"python", "-c", "import sys; sys.exit(4)"}, &out, &out)
	if code != 4 || !strings.Contains(out.String(), "exit 4") {
		t.Fatalf("code %d\n%s", code, out.String())
	}
}
