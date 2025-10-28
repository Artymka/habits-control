package jwt

import (
	"testing"
	"time"

	"github.com/Artymka/habits-control/app/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	testCases := []struct {
		name      string
		userID    int64
		secretKey string
	}{
		{
			name:      "Test 1",
			userID:    123456,
			secretKey: "abcdefgh",
		},
		{
			name:      "Test 2",
			userID:    123456,
			secretKey: "abcdefgh",
		},
		{
			name:      "Test 3",
			userID:    98746621,
			secretKey: "ja;sdf;jasdf5asdf6asd4f6asd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.JWT.SecretKey = tc.secretKey
			token, err := CreateToken(JWTData{UserID: tc.userID}, &cfg)
			assert.Nil(t, err)
			t.Logf("token: %s\n", token)
		})
	}
}

func TestParseTokenAllright(t *testing.T) {
	testCases := []struct {
		name      string
		userID    int64
		secretKey string
		ttl       int
	}{
		{
			name:      "Test 1",
			userID:    123456,
			secretKey: "abcdefgh",
			ttl:       10,
		},
		{
			name:      "Test 2",
			userID:    98746621,
			secretKey: "ja;sdf;jasdf5asdf6asd4f6asd",
			ttl:       10,
		},
		{
			name:      "Test 3",
			userID:    0,
			secretKey: "aa",
			ttl:       10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.JWT.SecretKey = tc.secretKey
			cfg.JWT.TokenSecondsTTL = tc.ttl

			original := JWTData{UserID: tc.userID}
			token, err := CreateToken(original, cfg)
			assert.Nil(t, err)
			t.Logf("created token: %s\n", token)

			created, err := ParseToken(token, cfg)
			assert.Nil(t, err)
			assert.Equal(t, original, created)
		})
	}
}

func TestParseTokenError(t *testing.T) {
	testCases := []struct {
		name      string
		userID    int64
		secretKey string
		ttl       int
		err       error
		wait      int
	}{
		{
			name:      "Token expired",
			userID:    248897,
			secretKey: "abcdefgh",
			ttl:       1,
			err:       ErrTokenExpired,
			wait:      2,
		},
		{
			name:      "Token expired 2",
			userID:    249292,
			secretKey: "kokokokok",
			ttl:       2,
			err:       ErrTokenExpired,
			wait:      3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.JWT.SecretKey = tc.secretKey
			cfg.JWT.TokenSecondsTTL = tc.ttl

			original := JWTData{UserID: tc.userID}
			token, err := CreateToken(original, cfg)
			assert.Nil(t, err)
			t.Logf("created token: %s\n", token)

			if tc.wait != 0 {
				time.Sleep(time.Second * time.Duration(tc.wait))
			}

			_, err = ParseToken(token, cfg)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
