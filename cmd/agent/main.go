package main

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/service/agent"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel  := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer ctxCancel()

	config := agent.CreateConfig()
	agent, err := agent.NewAgent(config)

	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	if err := agent.Run(ctx); err != nil {
		log.Fatalf("Failed to run an agent: %v", err)
	}
}