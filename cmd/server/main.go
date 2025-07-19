package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"temporal-crm-ingestor/config"

	"go.temporal.io/sdk/client"
)

const webhookPath = "/webhook"

func main() {
	// Load config
	err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Setup Temporal client
	tClient, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer tClient.Close()

	http.HandleFunc(webhookPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Simple token check
		if r.Header.Get("X-Auth-Token") != config.Cfg.ExpectedToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse into map
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		log.Printf("✅ Received webhook payload: %+v\n", payload)

		// Start Temporal workflow
		workflowOptions := client.StartWorkflowOptions{
			TaskQueue: "lead-task-queue",
		}

		we, err := tClient.ExecuteWorkflow(r.Context(), workflowOptions, "CreateLeadWorkflow", payload)
		if err != nil {
			log.Printf("❌ Failed to start workflow: %v", err)
			http.Error(w, "workflow start error", http.StatusInternalServerError)
			return
		}

		log.Printf("🚀 Workflow started. ID: %s RunID: %s\n", we.GetID(), we.GetRunID())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Webhook received and workflow started")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🌐 Listening on http://localhost:%s%s", port, webhookPath)
	http.ListenAndServe(":"+port, nil)
}
