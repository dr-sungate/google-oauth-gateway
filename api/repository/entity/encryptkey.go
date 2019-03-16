package entity

import (
	"crypto/rsa"
)

type EncryptKey struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []*rsa.PublicKey
}
