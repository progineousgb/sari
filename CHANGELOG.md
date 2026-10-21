# Changelog

## v1.0.3 (2026-10-22)

- Quoted values have their quotes removed: double-quoted basic strings (with `\u` and `\r` escapes), single-quoted literal strings, and a trailing `#` comment after a quoted value.
- Quoted task names are accepted.
- README describes the supported TOML subset instead of claiming full TOML.
- Add a CI workflow.

## v1.0.2 (2026-10-21)

Commit email update. No code changes.

## v1.0.1 (2026-10-21)

Commit metadata cleanup. No code changes.

## v1.0.0 (2026-10-21)

First tagged release. Run project tasks from sari.toml and print how long they took.

- Add sari, a runner for tasks in sari.toml.
- Document sari's task file with a real run and time.
- README: add Linux and macOS steps, drop machine-specific notes
- tests: use a helper process instead of python
