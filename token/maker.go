package token

import (
	"errors"
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
    // CreateToken creates a new token for a specific username and duration
    CreateToken(username string, duration time.Duration) (string, error)

    // VerifyToken checks if the token is valid or not
    VerifyToken(token string) (*Payload, error)
}

var (
    ErrInvalidToken = errors.New("token is invalid")
    ErrExpiredToken = errors.New("token has expired")
)