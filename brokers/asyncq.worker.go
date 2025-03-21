package brokers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

type Worker struct {
	server *asynq.Server
}

func NewWorker(redisAddress string) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddress},
		asynq.Config{Concurrency: 10},
	)
	return &Worker{server: server}
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc("email:send", handleEmailTask)

	log.Println("Worker started... listening for jobs.")
	return w.server.Run(mux)
}

func handleEmailTask(ctx context.Context, task *asynq.Task) error {
	var payload JobPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Processing email task for user: %s, Title: %s", payload.UserId, payload.Title)
	return nil
}

// Usage below

/*
job := brokers.JobPayload{
	UserId:      "123",
	TaskType:    "email:send",
	Title:       "Welcome Email",
	TaskDetails: "Send welcome email to user",
	Payload:     "Welcome to our service!",
}

err = queueService.EnqueueJob(job, "default")
if err != nil {
	log.Fatalf("Failed to enqueue job: %v", err)
}
*/
