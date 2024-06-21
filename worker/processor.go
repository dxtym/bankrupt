package worker

import (
	"context"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

// for mock testing
type TaskProcessor interface {
	Run() error
	ProcessorTaskSendEmail(ctx context.Context, task *asynq.Task) error
}

// task processor
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		opts,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (rtp RedisTaskProcessor) Run() error {
	mux := asynq.NewServeMux()

	// register task handlers
	mux.HandleFunc(TaskSendEmail, rtp.ProcessorTaskSendEmail)
	return rtp.server.Start(mux)
}
