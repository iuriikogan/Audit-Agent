package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"multi-agent-cra/pkg/agent"
	"multi-agent-cra/pkg/domain"
)

// Coordinator acts as the Concurrency Agent, orchestrating the flow of data

// between specialized agents using Go channels and goroutines.

type Coordinator struct {

	aggregator agent.Agent

	modeler    agent.Agent

	validator  agent.Agent

	reviewer   agent.Agent

	tagger     agent.Agent

	concurrency int

}



// NewCoordinator initializes the workflow manager with specific agents.

func NewCoordinator(aggregator, modeler, validator, reviewer, tagger agent.Agent, workers int) *Coordinator {

	if workers <= 0 {

		workers = 1

	}

	return &Coordinator{

		aggregator:  aggregator,

		modeler:     modeler,

		validator:   validator,

		reviewer:    reviewer,

		tagger:      tagger,

		concurrency: workers,

	}

}





// ProcessStream takes a stream of products and returns a stream of results.

// It manages a worker pool to process items concurrently.

func (c *Coordinator) ProcessStream(ctx context.Context, input <-chan domain.Product) <-chan domain.AssessmentResult {

	results := make(chan domain.AssessmentResult)



	go func() {

		defer close(results)

		

		var wg sync.WaitGroup

		

		// Launch worker pool

		for i := 0; i < c.concurrency; i++ {

			wg.Add(1)

			go func(workerID int) {

				defer wg.Done()

				c.workerLoop(ctx, workerID, input, results)

			}(i)

		}

		

		wg.Wait()

	}()



	return results

}



// workerLoop consumes products and runs the agent pipeline for each.

func (c *Coordinator) workerLoop(ctx context.Context, id int, input <-chan domain.Product, output chan<- domain.AssessmentResult) {

	for p := range input {

		// Respect context cancellation

		select {

		case <-ctx.Done():

			return

		default:

		}



		res := c.analyzeProduct(ctx, p)

		output <- res

	}

}



// analyzeProduct executes the sequential logic for a single item:

// Aggregator -> Modeler -> Validator -> Reviewer -> Tagger

func (c *Coordinator) analyzeProduct(ctx context.Context, p domain.Product) domain.AssessmentResult {

	res := domain.AssessmentResult{ProductID: p.ID}



	// 1. Resource Aggregator

	stepCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)

	defer cancel()



	dataRepo, err := c.aggregator.Chat(stepCtx, fmt.Sprintf("Ingest resources for product: %s (Component: %s, Version: %s)", p.Name, p.Component, p.Version))

	if err != nil {

		res.Error = fmt.Errorf("aggregation failed: %w", err)

		return res

	}



	// 2. CRA Modeler

	stepCtx, cancel = context.WithTimeout(ctx, 2*time.Minute)

	defer cancel()



	complianceModel, err := c.modeler.Chat(stepCtx, fmt.Sprintf("Model CRA compliance for data: %s", dataRepo))

	if err != nil {

		res.Error = fmt.Errorf("modeling failed: %w", err)

		return res

	}



	// 3. Compliance Validator

	stepCtx, cancel = context.WithTimeout(ctx, 2*time.Minute)

	defer cancel()



	complianceReport, err := c.validator.Chat(stepCtx, fmt.Sprintf("Validate model against CRA rules: %s", complianceModel))

	if err != nil {

		res.Error = fmt.Errorf("validation failed: %w", err)

		return res

	}

	res.VulnReport = complianceReport // Mapping report to existing field for now



	// 4. Reviewer

	stepCtx, cancel = context.WithTimeout(ctx, 2*time.Minute)

	defer cancel()



	approval, err := c.reviewer.Chat(stepCtx, fmt.Sprintf("Review compliance report: %s", complianceReport))

	if err != nil {

		res.Error = fmt.Errorf("review failed: %w", err)

		return res

	}

	res.AuditStatus = approval



	// 5. Resource Tagger (only if issues found or requested)

	stepCtx, cancel = context.WithTimeout(ctx, 2*time.Minute)

	defer cancel()

	

	tags, err := c.tagger.Chat(stepCtx, fmt.Sprintf("Tag resources based on report: %s", complianceReport))

	if err != nil {

		// Non-blocking error, log but continue

		fmt.Printf("Tagging warning for %s: %v\n", p.ID, err)

	} else {

		res.Classification = tags // Storing tags in classification field for now

	}



	return res

}
