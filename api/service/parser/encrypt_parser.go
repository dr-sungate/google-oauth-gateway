package parser

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dr-sungate/google-oauth-gateway/api/repository/entity"
	"io/ioutil"
)

const (
	RSAPRIVATEKEY_MESSAGE = "RSA PRIVATE KEY"
	PRIVATEKEY_MESSAGE    = "PRIVATE KEY"
	PUBLICKEY_MESSAGE     = "PUBLIC KEY"
)

func ReadPrivateKey(filepath string, encryptkey *entity.EncryptKey) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return ReadPrivateKeyFromByte(bytes, encryptkey)
}

func ReadPublicKey(filepath string, encryptkey *entity.EncryptKey) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return ReadPublicKeyFromByte(bytes, encryptkey)
}

func ReadPrivateKeyFromByte(bytedata []byte, encryptkey *entity.EncryptKey) error {
	block, _ := pem.Decode(bytedata)
	if block == nil {
		return errors.New("failed to decode private key data")
	}
	var key *rsa.PrivateKey
	var err error
	switch block.Type {
	case RSAPRIVATEKEY_MESSAGE:
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
	case PRIVATEKEY_MESSAGE:
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		var ok bool
		key, ok = keyInterface.(*rsa.PrivateKey)
		if !ok {
			return errors.New("not RSA private key")
		}
	default:
		return fmt.Errorf("invalid private key type : %s", block.Type)
	}
	key.Precompute()
	encryptkey.PrivateKey = key
	return nil
}

func ReadPublicKeyFromByte(bytedata []byte, encryptkey *entity.EncryptKey) error {
	block, _ := pem.Decode(bytedata)
	if block == nil || block.Type != PUBLICKEY_MESSAGE {
		return errors.New("failed to decode PEM block containing public key")
	}
	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	key, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("not RSA public key")
	}
	encryptkey.PublicKey = key
	return nil
}

func DecodePrivateKeyPKCS1(pubkey *rsa.PrivateKey) []byte {
	prikey_bytes := x509.MarshalPKCS1PrivateKey(pubkey)
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  RSAPRIVATEKEY_MESSAGE,
			Bytes: prikey_bytes,
		},
	)
	return pemdata
}

func DecodePrivateKeyPKCS8(pubkey *rsa.PrivateKey) ([]byte, error) {
	prikey_bytes, err := x509.MarshalPKCS8PrivateKey(pubkey)
	if err != nil {
		return nil, err
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  PRIVATEKEY_MESSAGE,
			Bytes: prikey_bytes,
		},
	)
	return pemdata, nil
}

func DecodePublicKey(pubkey *rsa.PublicKey) ([]byte, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  PUBLICKEY_MESSAGE,
			Bytes: pubkey_bytes,
		},
	)
	return pemdata, nil
}
