package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

const (
	SecretKeySize = 32
)

type Encrypt struct {
	gcm cipher.AEAD
}

func NewEncrypt(cfg config.Encrypt) (*Encrypt, error) {
	if len(cfg.Key) != SecretKeySize {
		return nil, fmt.Errorf("secret key must be 32 length")
	}
	c, err := aes.NewCipher([]byte(cfg.Key))
	if err != nil {
		return nil, fmt.Errorf("error init cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error init gcm: %w", err)
	}
	return &Encrypt{
		gcm: gcm,
	}, nil
}

func (e *Encrypt) Encrypt(data string) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		logger.Error.Println("error read nonce:", err.Error())
		return "", status.Error(codes.Internal, "")
	}
	res := e.gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(res), nil
}

func (e *Encrypt) Decrypt(data string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "invalid string")
	}

	nonceSize := e.gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", status.Error(codes.InvalidArgument, "string too short")
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error.Println("error decrypt message:", err.Error())
		return "", status.Error(codes.Internal, "internal")
	}
	return string(plaintext), nil
}
