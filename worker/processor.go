package worker

import (
	"context"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/mail"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
    QUEUE_CRITICAL = "critical"
    QUEUE_DEFAULT  = "default"
    QUEUE_LOW      = "low"
)

type TaskProcessor interface {
    Start() error
    ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
    server *asynq.Server
    store  db.Store
    mailer mail.Sender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.Sender) TaskProcessor {
    server := asynq.NewServer(
        redisOpt,
        asynq.Config{
            Queues: map[string]int{
                QUEUE_CRITICAL: 10,
                QUEUE_DEFAULT: 5,
            },
            ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
                log.Error().
                    Err(err).
                    Str("type", task.Type()).
                    Bytes("payload", task.Payload()).
                    Msg("failed to process task")
            }),
            Logger: NewLogger(),
        },
    )

    return &RedisTaskProcessor{
        server: server,
        store:  store,
        mailer: mailer,
    }
}

func (processor *RedisTaskProcessor) Start() error {
    mux := asynq.NewServeMux()
    mux.HandleFunc(TASK_SEND_VERIFICATION_EMAIL, processor.ProcessTaskSendVerificationEmail)
    return processor.server.Start(mux)
}
