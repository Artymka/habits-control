package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Artymka/habits-control/app/internal/config"
)

type JWTData struct {
	UserID int64
}

var (
	ErrTokenExpired = jwt.ErrTokenExpired
	ErrInvalidToken = errors.New("jwt token is invalid")
	ErrInvalidField = errors.New("jwt token has invalid field")
)

func CreateToken(data JWTData, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": data.UserID,
		"exp":    time.Now().Add(time.Second * time.Duration(cfg.JWT.TokenSecondsTTL)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string, cfg *config.Config) (JWTData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})

	if err != nil {
		return JWTData{}, err
	}
	if !token.Valid {
		return JWTData{}, ErrInvalidToken
	}

	res := JWTData{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// converting time
		timestamp, ok := claims["exp"].(float64)
		if !ok {
			return JWTData{}, fmt.Errorf("%w: %s", ErrInvalidField, "exp")
		}
		expires_at := time.Unix(int64(timestamp), 0)
		if time.Until(expires_at) < 0 {
			return JWTData{}, ErrTokenExpired
		}

		floatUserID, ok := claims["userID"].(float64)
		if !ok {
			return JWTData{}, fmt.Errorf("%w: %s", ErrInvalidField, "userID")
		}
		res.UserID = int64(floatUserID)
		return res, nil

	} else {
		return JWTData{}, ErrInvalidToken
	}
}
