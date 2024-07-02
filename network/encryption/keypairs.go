package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

type KeypairBytes struct {
	PublicKey  []byte
	PrivateKey []byte
}

func MakeKeypairBytes() (KeypairBytes, error) {
	const bitSize = 1024

	var privateKey *rsa.PrivateKey
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return KeypairBytes{}, err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return KeypairBytes{}, err
	}

	log.Printf("Generated %d-bit RSA keypair\n", bitSize)

	return KeypairBytes{
		PublicKey:  publicKeyBytes,
		PrivateKey: privateKeyBytes,
	}, nil
}
