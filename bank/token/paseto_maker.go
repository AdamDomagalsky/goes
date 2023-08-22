package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
	"golang.org/x/crypto/chacha20poly1305"
)

type PASETOMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPASETOMaker(symmetricKey string) (Maker, error) {

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters long", chacha20poly1305.KeySize)
	}

	return &PASETOMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (maker *PASETOMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	// footer := "optional footer"

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil

}

func (maker *PASETOMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
