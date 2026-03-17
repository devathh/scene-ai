package jwtdomain

import "crypto/rsa"

type KeyLoader interface {
	LoadPublic() (*rsa.PublicKey, error)
	LoadPrivate() (*rsa.PrivateKey, error)
}
