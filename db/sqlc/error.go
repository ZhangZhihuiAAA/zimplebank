package db

import (
    "errors"

    "github.com/jackc/pgx/v5/pgconn"
)

const (
    FOREIGN_KEY_VIOLATION = "23503"
    UNIQUE_VIOLATION      = "23505"
)

func ErrorCode(err error) string {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return pgErr.Code
    }
    return ""
}
