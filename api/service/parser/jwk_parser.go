package parser

import (
	"crypto/rsa"
	"errors"
	jose "gopkg.in/square/go-jose.v2"
)

func ConvertJSONWebKey(input []byte) (jose.JSONWebKey, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON(input)
	return jwk, err
}

func ToRSAPublic(key *jose.JSONWebKey) (*rsa.PublicKey, error) {
	res, ok := key.Key.(*rsa.PublicKey)
	if !ok {
		return res, errors.New("Could not convert key to RSA Private Key.")
	}
	return res, nil
}

func ToRSAPrivate(key *jose.JSONWebKey) (*rsa.PrivateKey, error) {
	res, ok := key.Key.(*rsa.PrivateKey)
	if !ok {
		return res, errors.New("Could not convert key to RSA Private Key.")
	}
	return res, nil
}
