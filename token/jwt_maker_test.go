package token

import (
	"testing"
	"time"

	"github.com/adifahmi/simplebank/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32, false))
	require.NoError(t, err)

	username := util.RandomOwner()
	userID := util.RandomInteger(1, 1000)
	duration := time.Hour
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	vPayload, vErr := maker.VerifyToken(token)
	require.NoError(t, vErr)
	require.NotEmpty(t, vPayload)

	require.Equal(t, username, vPayload.Username)
	require.WithinDuration(t, issuedAt, vPayload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, vPayload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32, false))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.RandomInteger(1, 1000), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.RandomInteger(1, 1000), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32, false))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
