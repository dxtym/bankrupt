package token

import (
	"testing"
	"time"

	"github.com/dxtym/bankrupt/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker([]byte(utils.RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := utils.RandomOwner()
	duration := time.Minute
	createdAt := time.Now()
	expiredAt := createdAt.Add(duration)

	jwtToken, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(jwtToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker([]byte(utils.RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	jwtToken, payload, err := pasetoMaker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(jwtToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
