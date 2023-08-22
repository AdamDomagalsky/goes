package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
	payload   *Payload
}

func NewJWTMaker(secretKey string) (Maker, error) {

	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey, &Payload{}}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	id := payload.ID.String()
	if id == "" {
		return "", nil, ErrInvalidToken
	}
	jwtPayload := jwt.RegisteredClaims{
		ID:        id,
		Subject:   payload.Username,
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}
	return signedToken, payload, nil
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		if errors.Is(err, jwt.ErrTokenUnverifiable) {
			return nil, ErrInvalidToken
		}

		return nil, ErrInvalidToken
	}
	jwtPayload, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(jwtPayload.ID)

	if err != nil {
		return nil, ErrInvalidToken

	}
	return &Payload{
		ID:        id,
		Username:  jwtPayload.Subject,
		IssuedAt:  jwtPayload.IssuedAt.Time,
		ExpiresAt: jwtPayload.ExpiresAt.Time,
	}, nil
}
