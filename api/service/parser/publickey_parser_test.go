package parser

import (
	"fmt"
	"github.com/dr-sungate/google-oauth-gateway/api/repository/entity"
	"reflect"
	"testing"
)

func Test_ConvertPublicKey(t *testing.T) {
	publickeyStr := `{
      "kid": "7c309e3a1c1999cb0404ab7125ee40b7cdbcaf7d",
      "e": "AQAB",
      "kty": "RSA",
      "alg": "RS256",
      "n": "3MdFK4pXPvehMipDL_COfqn6o9soHgSaq_V1o8U_5gTZ-j9DxO9PV7BVncXBgHFctnp3JQ1QTDF7txeHeuLOS4KziRw5r4ohaj2WoOTqXh7lqVMR2YDAcBK46asS177NpkQ1CqHIsy3kNfqhXLwTaKfdlwdA_XUfRbKORWbq0kDxV35egx35nHl5qJ6aP6fcpsnnPvHf7KWO0zkdvwuR-IX79HjqUAEg5UERd5FK4y06PRbxuXHjAgVhHu_sk4reNXNp1HRuTYtQ26DFbVaIjsWb8-nQC8-7FkTjlw9FteAwLVGOm9sTLFp73jAf0pWLh7sJ02pBxZKjsxLO1Lvg7w",
      "use": "sig"
    }`
	jwk, err := ConvertJSONWebKey([]byte(publickeyStr))
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	publickey, err := ToRSAPublic(&jwk)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(publickey)
	pempublickey, err := DecodePublicKey(publickey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(fmt.Sprintf("%v", string(pempublickey)))
	encryptkey := &entity.EncryptKey{}
	err = ReadPublicKeyFromByte(pempublickey, encryptkey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(encryptkey.PublicKey)
	if reflect.DeepEqual(encryptkey.PublicKey, publickey) == false {
		t.Fatalf("failed compare publickey ")
	}
	t.Log("success PublicKeyConvertTest")
}

func Test_ConvertPrivatecKey(t *testing.T) {
	privatekeyStr := `{
    "kty":"RSA",
    "kid":"juliet@capulet.lit",
    "use":"enc",
    "n":"t6Q8PWSi1dkJj9hTP8hNYFlvadM7DflW9mWepOJhJ66w7nyoK1gPNqFMSQRyO125Gp-TEkodhWr0iujjHVx7BcV0llS4w5ACGgPrcAd6ZcSR0-Iqom-QFcNP8Sjg086MwoqQU_LYywlAGZ21WSdS_PERyGFiNnj3QQlO8Yns5jCtLCRwLHL0Pb1fEv45AuRIuUfVcPySBWYnDyGxvjYGDSM-AqWS9zIQ2ZilgT-GqUmipg0XOC0Cc20rgLe2ymLHjpHciCKVAbY5-L32-lSeZO-Os6U15_aXrk9Gw8cPUaX1_I8sLGuSiVdt3C_Fn2PZ3Z8i744FPFGGcG1qs2Wz-Q",
    "e":"AQAB",
    "d":"GRtbIQmhOZtyszfgKdg4u_N-R_mZGU_9k7JQ_jn1DnfTuMdSNprTeaSTyWfSNkuaAwnOEbIQVy1IQbWVV25NY3ybc_IhUJtfri7bAXYEReWaCl3hdlPKXy9UvqPYGR0kIXTQRqns-dVJ7jahlI7LyckrpTmrM8dWBo4_PMaenNnPiQgO0xnuToxutRZJfJvG4Ox4ka3GORQd9CsCZ2vsUDmsXOfUENOyMqADC6p1M3h33tsurY15k9qMSpG9OX_IJAXmxzAh_tWiZOwk2K4yxH9tS3Lq1yX8C1EWmeRDkK2ahecG85-oLKQt5VEpWHKmjOi_gJSdSgqcN96X52esAQ",
    "p":"2rnSOV4hKSN8sS4CgcQHFbs08XboFDqKum3sc4h3GRxrTmQdl1ZK9uw-PIHfQP0FkxXVrx-WE-ZEbrqivH_2iCLUS7wAl6XvARt1KkIaUxPPSYB9yk31s0Q8UK96E3_OrADAYtAJs-M3JxCLfNgqh56HDnETTQhH3rCT5T3yJws",
    "q":"1u_RiFDP7LBYh3N4GXLT9OpSKYP0uQZyiaZwBtOCBNJgQxaj10RWjsZu0c6Iedis4S7B_coSKB0Kj9PaPaBzg-IySRvvcQuPamQu66riMhjVtG6TlV8CLCYKrYl52ziqK0E_ym2QnkwsUX7eYTB7LbAHRK9GqocDE5B0f808I4s",
    "dp":"KkMTWqBUefVwZ2_Dbj1pPQqyHSHjj90L5x_MOzqYAJMcLMZtbUtwKqvVDq3tbEo3ZIcohbDtt6SbfmWzggabpQxNxuBpoOOf_a_HgMXK_lhqigI4y_kqS1wY52IwjUn5rgRrJ-yYo1h41KR-vz2pYhEAeYrhttWtxVqLCRViD6c",
    "dq":"AvfS0-gRxvn0bwJoMSnFxYcK1WnuEjQFluMGfwGitQBWtfZ1Er7t1xDkbN9GQTB9yqpDoYaN06H7CFtrkxhJIBQaj6nkF5KKS3TQtQ5qCzkOkmxIe3KRbBymXxkb5qwUpX5ELD5xFc6FeiafWYY63TmmEAu_lRFCOJ3xDea-ots",
    "qi":"lSQi-w9CpyUReMErP1RsBLk7wNtOvs5EQpPqmuMvqW57NBUczScEoPwmUqqabu9V0-Py4dQ57_bapoKRu1R90bvuFnU63SHWEFglZQvJDMeAvmj4sm-Fp0oYu_neotgQ0hzbI5gry7ajdYy9-2lNx_76aBZoOUu9HCJ-UsfSOI8"
    }`
	jwk, err := ConvertJSONWebKey([]byte(privatekeyStr))
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	privatekey, err := ToRSAPrivate(&jwk)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(privatekey)

	pemprikey1 := DecodePrivateKeyPKCS1(privatekey)
	t.Log(fmt.Sprintf("%v", string(pemprikey1)))
	encryptkey := &entity.EncryptKey{}
	err = ReadPrivateKeyFromByte(pemprikey1, encryptkey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	pemprikeyagain1 := DecodePrivateKeyPKCS1(encryptkey.PrivateKey)
	if reflect.DeepEqual(pemprikeyagain1, pemprikey1) == false {
		t.Fatalf("failed compare privatekey PKCS#1")
	}

	pemprikey2, err := DecodePrivateKeyPKCS8(privatekey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(fmt.Sprintf("%v", string(pemprikey2)))
	encryptkey = &entity.EncryptKey{}
	err = ReadPrivateKeyFromByte(pemprikey2, encryptkey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	pemprikeyagain2, err := DecodePrivateKeyPKCS8(encryptkey.PrivateKey)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if reflect.DeepEqual(pemprikeyagain2, pemprikey2) == false {
		t.Fatalf("failed compare privatekey PKCS#8")
	}

	t.Log("success PrivateKeyConvertTest")
}
