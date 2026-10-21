# sari

sari runs the commands listed in a TOML file, in order, and prints how long each one took. It runs once and exits; it does not stay in the background.

## What it does

- `sari run` reads `sari.toml`, runs each task from top to bottom, and prints a `$ <name>` line before each task's own output.
- After the tasks it prints a table with three columns: `task`, `exit` and `ms` (task name, exit code, duration in milliseconds).
- By default sari stops at the first task that exits non-zero. `--keep-going` runs the rest anyway. The table is printed in both cases.
- `sari time -- <cmd> [args...]` runs one program, prints its output, then `exit <code>` and `duration <ms>ms`.

## Install

Requires Go 1.24 or newer (from `go.mod`).

```sh
go install github.com/progineousgb/sari@latest
```

This puts the `sari` binary in `$(go env GOPATH)/bin`. Make sure that directory is on your PATH.

### Build from source

Linux and macOS:

```sh
git clone https://github.com/progineousgb/sari.git
cd sari
go build -o sari .
./sari --help
```

Windows (PowerShell or cmd):

```text
git clone https://github.com/progineousgb/sari.git
cd sari
go build -o sari.exe .
.\sari.exe --help
```

Run the tests with:

```sh
go test ./...
```

## Usage

```text
sari run [--keep-going] [-f sari.toml]
sari time -- <command> [args...]
```

- `sari` with no arguments is the same as `sari run`.
- `-f` or `--file` picks a different task file. The default is `sari.toml` in the current directory.
- `help`, `-h` and `--help` print the usage and exit 0.

### Example

`sari.toml`:

```toml
[tasks]
# build steps
hello = echo hello
fail = sh -c "exit 3"
after = echo done
```

`sari run` on Linux:

```text
$ hello
hello
$ fail

task   exit  ms
hello     0  1
fail      3  1
```

sari stopped after `fail` and exited 3. `after` did not run.

`sari run --keep-going`:

```text
$ hello
hello
$ fail
$ after
done

task   exit  ms
hello     0  0
fail      3  1
after     0  1
```

All three tasks ran and sari still exited 3, the first non-zero code.

`sari time -- echo hi`:

```text
hi
exit 0
duration 2ms
```

The millisecond values change from run to run.

## Task file format

The file holds one `[tasks]` section. Each line inside it is `name = command`.

- Blank lines are ignored. A line whose first non-blank character is `#` is a comment.
- The name is shown as written in the `$` line and in the table. It must be non-empty, valid UTF-8, and must not contain `"`.
- If the command starts with `"`, it is a quoted string. Supported escapes are `\\`, `\"`, `\n` and `\t`. Text after the closing quote, or a missing closing quote, is an error.
- Otherwise the command is the rest of the line, taken as is.
- A UTF-8 BOM at the start of the file is allowed.

These are errors: no `[tasks]` section, a second `[tasks]`, any other `[section]`, a line before `[tasks]`, a line without `=`, and an empty command. sari prints `sari: <reason>` (with the line number where there is one) on standard error and exits 2.

## How commands are run

- Tasks in `sari run` go through a shell: `sh -c "<command>"` on Linux and macOS. On Windows sari uses `%ComSpec%` (or `cmd.exe` if unset) with `/d /s /c` followed by the command string unchanged.
- `sari time` does not use a shell. The first word after `--` is the program and the rest are its arguments.
- Task output (stdout and stderr) goes straight to the terminal.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Every task (or the timed command) exited 0 |
| N | The first non-zero exit code from a task, or the timed command's exit code |
| 1 | A command could not be started at all (for example `sari time` with a program that does not exist); sari prints `sari: <error>` |
| 2 | Usage error, unknown command or flag, unreadable task file, invalid task file, or `sari time` with no command |

## Limits

- One `[tasks]` list per file. No other tables, no dependencies between tasks, no parallel runs.
- No scheduling. sari runs when you call it.
- No file watching. It does not rerun when files change.

## Maintained by

Maintained by Pro Gineous: https://progineous.com

## License

MIT, see [LICENSE](LICENSE).
