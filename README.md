# sari

sari runs the commands listed in a TOML file, in order, and prints how long each one took. You run it when you want the list; it does not stay running in the background.

## What you get

`sari run` prints each task name, that task's own output, and then a table. The table columns are `task`, `exit`, and `ms`: the name, the exit code, and the duration in milliseconds.

`sari time` times one program. It prints that program's output, then `exit <code>` and `duration <milliseconds>ms`.

## Install and run

The module in `go.mod` uses Go 1.24. From this directory the build is:

```text
go build -o sari.exe .
```

On this machine `go` on PATH is 1.25.1 and `GOROOT` is the 1.24.1 tree, so that command stops with a compiler version mismatch. The binary used for the output below was built with:

```text
C:\laragon\bin\go\go1.24.1\bin\go.exe build -o sari.exe .
```

In PowerShell the program is `.\sari.exe`, run from this directory:

```text
.\sari.exe run
.\sari.exe time -- python -c "print(123)"
```

`.\sari.exe` with no arguments reads `sari.toml` in the current directory and runs it, same as `.\sari.exe run`.

## Example

The file `sari.toml` in this repo is:

```toml
[tasks]
print = python -c "print(123)"
```

`.\sari.exe run` printed:

```text
$ print
123

task   exit  ms
print     0  73
```

`73` is the millisecond count from that run. The next run will print a different number. The exit code was 0, and `.\sari.exe` itself exited 0.

`.\sari.exe time -- python -c "print(123)"` printed:

```text
123
exit 0
duration 91ms
```

`sari time` started `python` directly. It did not pass the command through `cmd.exe`. The token after `--` is the program, and the rest are its arguments.

## Flags and files

```text
sari run [--keep-going] [-f sari.toml]
sari time -- <command> [args...]
```

`help`, `-h`, and `--help` print those two lines and exit 0.

- No arguments, or `run`, executes the task file from top to bottom. The default file is `sari.toml` in the current directory.
- `-f` and `--file` select a different task file.
- `--keep-going` keeps going after a task exits non-zero. Without it, sari stops at the first non-zero exit and still prints the table for the tasks it did run. sari's own exit code is that first non-zero code, or 0 when every task exits 0.
- `time` requires a command. `sari time --` with nothing after it prints `sari: missing command` and exits 2.

An unknown command or an unknown flag prints the usage on standard error and exits 2.

The task file is a single `[tasks]` list. Blank lines are ignored. A line whose first non-whitespace character is `#` is a comment. Every other line is `name = command`. The name is what the `$` line and the table show, unchanged. A name must be non-empty and must not contain `"`.

If the command starts with `"`, it is a quoted string. The escapes are `\\`, `\"`, `\n`, and `\t`. Trailing text after the closing quote is an error, and an unclosed quote is an error. A command that does not start with `"` is taken as the rest of the line.

A missing `[tasks]`, a second `[tasks]`, any other `[section]`, a line with no `=`, or an empty command makes sari print `sari: ...` on standard error and exit 2. The file may start with a UTF-8 BOM.

From this directory, the tests are:

```text
go test ./...
```

## Limits

sari only reads a `[tasks]` list. It does not read other tables, and it does not schedule or watch the file.

On Windows, sari runs each task line through `cmd.exe`. It uses `ComSpec` when that variable is set, otherwise `cmd.exe`, and it builds the process command line as `/d /s /c` followed by the task string unchanged. sari does not split that string into arguments. `sari time` is separate: the program and arguments after `--` are started directly, with no shell.
