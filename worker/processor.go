package worker

import (
	"context"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
			ErrorHandler: asynq.ErrorHandlerFunc(func (ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("task processing failed")
			}),
			Logger: NewLogger(),
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
