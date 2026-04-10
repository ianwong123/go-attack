//go:build !windows

package internal

import "os/exec"

func hideWindow(cmd *exec.Cmd) {}
