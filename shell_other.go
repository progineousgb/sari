//go:build !windows

package main

import "os/exec"

func commandFromLine(line string) *exec.Cmd {
	return exec.Command("sh", "-c", line)
}
