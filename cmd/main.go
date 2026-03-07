package main

import (

	"context"

	"flag"

	"fmt"

	"log"

	"net/http"

	"os"

	"time"



	"github.com/google/generative-ai-go/genai"

	"google.golang.org/api/option"



	"multi-agent-cra/pkg/agent"

	"multi-agent-cra/pkg/domain"

	"multi-agent-cra/pkg/tools"

	"multi-agent-cra/pkg/workflow"

)



func main() {

	role := flag.String("role", "all", "The agent role to run (classifier, auditor, vuln, reporter, or all)")

	mode := flag.String("mode", "batch", "The execution mode (batch or server)")

	flag.Parse()



	ctx := context.Background()



	// Start health check server for container orchestration

	go func() {

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

			w.WriteHeader(http.StatusOK)

			fmt.Fprintln(w, "OK")

		})

		log.Printf("Starting health check server on :8080")

		if err := http.ListenAndServe(":8080", nil); err != nil {

			log.Printf("Health check server failed: %v", err)

		}

	}()



	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {

		log.Fatal("GEMINI_API_KEY is not set")

	}



	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {

		log.Fatal(err)

	}

	defer client.Close()



	if *mode == "server" {

		log.Printf("Running in SERVER mode as role: %s", *role)

		// In a real microservice, we would listen for work here (e.g., via Pub/Sub or gRPC)

		// For this demo, we'll just keep the process alive to satisfy the Deployment.

		select {}

	}



	// Default: Run Batch Assessment

	runBatch(ctx, client)

}



func runBatch(ctx context.Context, client *genai.Client) {

		// --- 1. Initialize Agents ---

		

		aggregatorAgent := agent.New(client, "ResourceAggregator", "Ingestion", "gemini-3.1-flash-lite",

			agent.WithSystemInstruction(`You are a Resource Aggregator. 

			Your task is to ingest relevant project artifacts (code, configs, docs) for a given product and output a structured data repository summary.`),

			agent.WithTools(tools.ScopeTools...), // Reusing existing tools for now, would need specific ones

		)

	

		modelerAgent := agent.New(client, "CRAModeler", "Modeling", "gemini-3.1-flash-lite",

			agent.WithSystemInstruction(`You are a CRA Modeler.

			Your task is to take a structured data repository and apply the Cyber Resilience Act (CRA) compliance framework to generate a compliance model.`),

		)

	

		validatorAgent := agent.New(client, "ComplianceValidator", "Validation", "gemini-3.1-flash-lite",

			agent.WithSystemInstruction(`You are a Compliance Validator.

			Your task is to validate a compliance model against CRA rules and output a compliance report with findings and deviations.`),

			agent.WithTools(tools.RegulatoryCheckerTools...),

		)

	

		reviewerAgent := agent.New(client, "Reviewer", "Approval", "gemini-3.1-flash-lite",

			agent.WithSystemInstruction(`You are a Compliance Reviewer.

			Your task is to review the compliance report and provide an approval status and final report summary.`),

		)

	

		taggerAgent := agent.New(client, "ResourceTagger", "Tagging", "gemini-3.1-flash-lite",

			agent.WithSystemInstruction(`You are a Resource Tagger.

			Your task is to tag resources that have issues identified in the compliance report with information on how to solve them.`),

		)

	

		// --- 2. Initialize Concurrency Coordinator ---

		

		// We use a worker pool of 3 to process products in parallel.

		coordinator := workflow.NewCoordinator(aggregatorAgent, modelerAgent, validatorAgent, reviewerAgent, taggerAgent, 3)

	



	// --- 3. Stream Processing ---



	products := []domain.Product{

		{ID: "p1", Name: "Smart Thermostat", Component: "FreeRTOS", Version: "10.4.1"},

		{ID: "p2", Name: "Enterprise Firewall", Component: "OpenSSL", Version: "1.1.1"},

		{ID: "p3", Name: "Mobile Banking App", Component: "React Native", Version: "0.66"},

		{ID: "p4", Name: "Wireless Keyboard", Component: "Bluetooth Stack", Version: "4.2"},

		{ID: "p5", Name: "Industrial Controller", Component: "BusyBox", Version: "1.30"},

	}



	fmt.Printf("--- Starting Concurrent Assessment for %d Products ---\n", len(products))

	start := time.Now()



	// Create input channel

	inputChan := make(chan domain.Product, len(products))

	for _, p := range products {

		inputChan <- p

	}

	close(inputChan)



	// Consume results from the coordinator

	resultsChan := coordinator.ProcessStream(ctx, inputChan)



	for res := range resultsChan {

		if res.Error != nil {

			fmt.Printf("[❌ ERROR] Product %s: %v\n", res.ProductID, res.Error)

			continue

		}

		fmt.Printf("[✅ DONE] Product: %s | Class: %s | Audit: %s | Vuln: %s\n", 

			res.ProductID, res.Classification, res.AuditStatus, res.VulnReport)

	}



	fmt.Printf("--- Completed in %v ---\n", time.Since(start))

}
