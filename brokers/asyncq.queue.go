package brokers

import (
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

type JobPayload struct {
	UserId      string `json:"user_id"`
	TaskType    string `json:"task_type"`
	Title       string `json:"title"`
	TaskDetails string `json:"task_details"`
	Payload     string `json:"payload"`
}

type QueueClient struct {
	client *asynq.Client
}

func NewQueueClient(redisAddress string) *QueueClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddress})
	return &QueueClient{client: client}
}

func (qc *QueueClient) EnqueueJob(payload JobPayload, queueName string) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(payload.TaskType, payloadBytes)

	// Corrected: Removed asynq.WithContext, added WithProcessIn (5s delay for demo)
	info, err := qc.client.Enqueue(task, asynq.Queue(queueName))
	if err != nil {
		return err
	}

	log.Printf("Enqueued task: ID=%s Queue=%s", info.ID, info.Queue)
	return nil
}

func (qc *QueueClient) Close() {
	qc.client.Close()
}
