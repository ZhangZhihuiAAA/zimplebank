package db

import (
    "context"

    "github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
    EmailID    int64
    SecretCode string
}

type VerifyEmailTxResult struct {
    User              User
    VerificationEmail VerificationEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
    var result VerifyEmailTxResult

    err := store.execTx(ctx, func(q *Queries) error {
        var err error

        result.VerificationEmail, err = q.UpdateVerificationEmail(ctx, UpdateVerificationEmailParams{
            ID:         arg.EmailID,
            SecretCode: arg.SecretCode,
        })
        if err != nil {
            return err
        }

        result.User, err = q.UpdateUser(ctx, UpdateUserParams{
            Username: result.VerificationEmail.Username,
            IsEmailVerified: pgtype.Bool{
                Bool:  true,
                Valid: true,
            },
        })

        return err
    })

    return result, err
}
