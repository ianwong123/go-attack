package protocol

import "time"

// beacon is sent from agent to server on every checkin
type Beacon struct {
	AgentID  string    `json:"agent_id"`
	Hostname string    `json:"hostname"`
	Username string    `json:"username"`
	OS       string    `json:"os"`
	Arch     string    `json:"arch"`
	IPs      []string  `json:"ips"`
	PID      int       `json:"pid"`
	CheckIn  time.Time `json:"check_in"`
	Result   *Result   `json:"result,omitempty"`
}

// task sent from server to agent
type Task struct {
	TaskID  string `json:"task_id"`
	Command string `json:"command"`
	Args    string `json:"args,omitempty"`
}

// result embedded in next beacon
type Result struct {
	TaskID  string `json:"task_id"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
}
