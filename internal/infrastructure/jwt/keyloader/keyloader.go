package jwtkeyloader

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type KeyLoader struct {
	publicKeyPath  string
	privateKeyPath string
}

func New(pubPath, privPath string) *KeyLoader {
	return &KeyLoader{
		publicKeyPath:  pubPath,
		privateKeyPath: privPath,
	}
}

func (kl *KeyLoader) LoadPrivate() (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(kl.privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key pem: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}

func (kl *KeyLoader) LoadPublic() (*rsa.PublicKey, error) {
	bytes, err := os.ReadFile(kl.publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key pem: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return key, nil
}
