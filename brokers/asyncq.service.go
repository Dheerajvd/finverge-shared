package brokers

import (
	"log"
)

type QueueService struct {
	queueClient *QueueClient
	worker      *Worker
}

func NewQueueService(redisAddress string) *QueueService {
	return &QueueService{
		queueClient: NewQueueClient(redisAddress),
		worker:      NewWorker(redisAddress),
	}
}

func (qs *QueueService) StartWorker() {
	go func() {
		if err := qs.worker.Start(); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()
}

func (qs *QueueService) EnqueueJob(payload JobPayload, queueName string) error {
	return qs.queueClient.EnqueueJob(payload, queueName)
}

func (qs *QueueService) Close() {
	qs.queueClient.Close()
}
