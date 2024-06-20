package worker

import (
	"context"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Run() error
	ProcessorTaskSendEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store
}

func NewRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		opts,
		asynq.Config{},
	)

	return &RedisTaskProcessor{
		server: server,
		store: store,
	}
}

func (rtp RedisTaskProcessor) Run() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendEmail, rtp.ProcessorTaskSendEmail)
	return rtp.server.Start(mux)
}
