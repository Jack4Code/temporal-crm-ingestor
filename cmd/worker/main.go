package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-crm-ingestor/config"
	"temporal-crm-ingestor/internal/workflows"
)

func main() {
	// Load config
	err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to Temporal
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, "lead-task-queue", worker.Options{})

	// Register workflows and activities
	w.RegisterWorkflow(workflows.CreateLeadWorkflow)
	w.RegisterActivity(workflows.CreateLeadActivity)

	log.Println("👷 Worker started for task queue: lead-task-queue")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("worker failed to start: %v", err)
	}
}
