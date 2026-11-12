package service

import (
	"context"
)

type encryptor interface {
	Encrypt(ctx context.Context, data string, encCtx map[string]string) (string, error)
}

type decryptor interface {
	Decrypt(data string, encCtx map[string]string) (string, error)
}

type Service struct {
	encryptor encryptor
	decryptor decryptor
}

func NewService(e encryptor, d decryptor) *Service {
	s := &Service{
		encryptor: e,
		decryptor: d,
	}

	return s
}
