package main

import (
	"strings"
	"testing"
)

func TestParseOrderQuotesAndUnicode(t *testing.T) {
	src := `
# comment

[tasks]
test = "python -m unittest"
print = python -c "print(123)"
نبض = python -c "print(7)"
`
	tasks, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len %d", len(tasks))
	}
	if tasks[0].Name != "test" || tasks[0].Cmd != "python -m unittest" {
		t.Fatalf("quoted: %#v", tasks[0])
	}
	if tasks[1].Cmd != `python -c "print(123)"` {
		t.Fatalf("unquoted: %#v", tasks[1])
	}
	if tasks[2].Name != "نبض" {
		t.Fatalf("name: %#v", tasks[2])
	}
}

func TestParseMissingSection(t *testing.T) {
	_, err := Parse("print = python -c \"print(1)\"\n")
	if err == nil || !strings.Contains(err.Error(), "[tasks]") {
		t.Fatal(err)
	}
}

func TestParseEmptyTaskList(t *testing.T) {
	tasks, err := Parse("[tasks]\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("%#v", tasks)
	}
}

func TestParseInvalidTaskName(t *testing.T) {
	_, err := Parse("[tasks]\n" + string([]byte{0xff}) + " = echo hi\n")
	if err == nil || !strings.Contains(err.Error(), "invalid task name") {
		t.Fatal(err)
	}
}

func TestParseUnclosedQuote(t *testing.T) {
	_, err := Parse("[tasks]\nbroken = \"python\n")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunFlags(t *testing.T) {
	path, keep, err := parseRunFlags([]string{"--keep-going", "-f", "other.toml"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "other.toml" || !keep {
		t.Fatalf("path %s keep %v", path, keep)
	}
}

func TestParseQuotedValues(t *testing.T) {
	src := "[tasks]\n" +
		"a = \"echo hi\"\n" +
		"b = 'echo \"x\"'\n" +
		"c = \"echo \\u0041\" # note\n" +
		"\"d e\" = echo d\n" +
		"'f' = 'C:\\path'\n"
	tasks, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	want := []Task{{"a", "echo hi"}, {"b", "echo \"x\""}, {"c", "echo A"}, {"d e", "echo d"}, {"f", "C:\\path"}}
	if len(tasks) != len(want) {
		t.Fatalf("%#v", tasks)
	}
	for i := range want {
		if tasks[i] != want[i] {
			t.Fatalf("task %d: got %#v want %#v", i, tasks[i], want[i])
		}
	}
}

func TestParseUnclosedSingleQuote(t *testing.T) {
	_, err := Parse("[tasks]\nx = 'echo\n")
	if err == nil || !strings.Contains(err.Error(), "unclosed quote") {
		t.Fatal(err)
	}
}
