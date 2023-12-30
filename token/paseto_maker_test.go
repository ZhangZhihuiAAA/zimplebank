package token

import (
	"testing"
	"time"

	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
    maker, err := NewPasetoMaker(util.RandomString(32))
    require.NoError(t, err)

    username := util.RandomOwner()
    role := util.DEPOSITOR_ROLE
    duration := time.Minute

    issuedAt := time.Now()
    expiresAt := issuedAt.Add(duration)

    token, payload, err := maker.CreateToken(username, role, duration)
    require.NoError(t, err)
    require.NotEmpty(t, token)
    require.NotEmpty(t, payload)

    payload, err = maker.VerifyToken(token)
    require.NoError(t, err)
    require.NotEmpty(t, payload)

    require.NotZero(t, payload.ID)
    require.Equal(t, username, payload.Username)
    require.Equal(t, role, payload.Role)
    require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
    require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
    maker, err := NewPasetoMaker(util.RandomString(32))
    require.NoError(t, err)

    token, payload, err := maker.CreateToken(util.RandomOwner(), util.DEPOSITOR_ROLE, -time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, token)
    require.NotEmpty(t, payload)

    payload, err = maker.VerifyToken(token)
    require.Error(t, err)
    require.EqualError(t, err, ErrExpiredToken.Error())
    require.Nil(t, payload)
}