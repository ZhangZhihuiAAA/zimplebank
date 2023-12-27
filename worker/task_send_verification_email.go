package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TASK_SEND_VERIFICATION_EMAIL = "task:send_verification_email"

type PayloadSendVerificationEmail struct {
    Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerificationEmail(
    ctx context.Context,
    payload *PayloadSendVerificationEmail,
    opts ...asynq.Option,
) error {
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal task payload: %w", err)
    }

    task := asynq.NewTask(TASK_SEND_VERIFICATION_EMAIL, jsonPayload, opts...)
    taskInfo, err := distributor.client.EnqueueContext(ctx, task)
    if err != nil {
        return fmt.Errorf("failed to enqueue task: %w", err)
    }

    log.Info().
        Str("type", task.Type()).
        Bytes("payload", task.Payload()).
        Str("queue", taskInfo.Queue).
        Int("max_retry", taskInfo.MaxRetry).
        Msg("enqueued task")
    return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error {
    var payload PayloadSendVerificationEmail
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
    }

    user, err := processor.store.GetUser(ctx, payload.Username)
    if err != nil {
        // if err == db.ErrNoRows {
        //     return fmt.Errorf("user not found: %w", asynq.SkipRetry)
        // }
        return fmt.Errorf("failed to get user: %w", err)
    }

    vEmail, err := processor.store.CreateVerificationEmail(ctx, db.CreateVerificationEmailParams{
        Username:   user.Username,
        Email:      user.Email,
        SecretCode: util.RandomString(32),
    })
    if err != nil {
        return fmt.Errorf("failed to create verification email: %w", err)
    }

    subject := "Welcome to Zimple Bank"
    vURL := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", vEmail.ID, vEmail.SecretCode)
    content := fmt.Sprintf(`Hello %s,<br/>
    Thank you for registering with us!<br/>
    Please <a href="%s">click here</a> to verify your email address.<br>
    `, user.FullName, vURL)
    to := []string{user.Email}
    err = processor.mailer.Send(subject, content, to, nil, nil, nil)
    if err != nil {
        return fmt.Errorf("failed to send verification email: %w", err)
    }

    log.Info().
        Str("type", task.Type()).
        Bytes("payload", task.Payload()).
        Str("email", user.Email).
        Msg("processed task")
    return nil
}
