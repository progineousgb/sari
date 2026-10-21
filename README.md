sari runs the commands listed in `sari.toml`, in order, then prints how long each one took. One command, no daemon.

Build:

```
go build -o sari.exe .
```

Run the tasks in `sari.toml` in the current directory:

```
sari run
```

`sari` with no arguments does the same thing. On Windows each task string is passed to `cmd.exe`. sari stops at the first non-zero exit unless `--keep-going` is set. The last lines are a table of task name, exit code, and milliseconds.

Time one command without a shell. Arguments after `--` are the program and its arguments:

```
sari time -- python -c "print(123)"
```

Example `sari.toml`:

```
[tasks]
print = python -c "print(123)"
```

Task names are printed as written, including non-ASCII characters.

```
go test ./...
```
