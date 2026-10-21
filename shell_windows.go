//go:build windows

package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func commandFromLine(line string) *exec.Cmd {
	comspec := os.Getenv("ComSpec")
	if comspec == "" {
		comspec = "cmd.exe"
	}
	cmd := exec.Command(comspec)
	exe := comspec
	if strings.ContainsAny(exe, " \t\"") {
		exe = `"` + strings.ReplaceAll(exe, `"`, `\"`) + `"`
	}
	// CmdLine is set explicitly so cmd.exe sees the task string unchanged.
	// os/exec quoting is for CommandLineToArgvW, which cmd /c does not use.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: exe + " /d /s /c " + line,
	}
	return cmd
}
