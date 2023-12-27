package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
    Querier
    TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
    CreateUserTx(context.Context, CreateUserTxParams) (CreateUserTxResult, error)
    VerifyEmailTx(context.Context, VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// SQLStore provides all functions to execute SQL quries and transactions
type SQLStore struct {
    connPool *pgxpool.Pool
    *Queries
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
    return &SQLStore{
        connPool: connPool,
        Queries:  New(connPool),
    }
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
    tx, err := store.connPool.Begin(ctx)
    if err != nil {
        return err
    }

    q := store.Queries.WithTx(tx)
    err = fn(q)
    if err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit(ctx)
}
