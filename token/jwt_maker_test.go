package token

import (
	"testing"
	"time"

	"github.com/dxtym/bankrupt/utils"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	jwtMaker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := utils.RandomOwner()
	duration := time.Minute
	createdAt := time.Now()
	expiredAt := createdAt.Add(duration)

	jwtToken, payload, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(jwtToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	jwtToken, payload, err := jwtMaker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(jwtToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jwtMaker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	payload, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
