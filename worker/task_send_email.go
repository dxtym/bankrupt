package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendEmail = "task:send_email"

type PayloadSendEmail struct {
	Username string `json:"username"`
}

// task distributor
func (rtd RedisTaskDistributor) DistributorTaskSendEmail(
	ctx context.Context,
	payload PayloadSendEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendEmail, jsonPayload, opts...) // create a new task
	info, err := rtd.client.EnqueueContext(ctx, task) // enqueue the task
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("task enqueued")
	return nil
}

// task processor
func (rtp RedisTaskProcessor) ProcessorTaskSendEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := rtp.store.GetUser(ctx, payload.Username) // get user from db
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// todo: send email to user
	log.Info().Str("username", user.Username).Str("email", user.Email).Msg("task processed")
	return nil
}
