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

// Decrypts data **in-place** using the server private key using `rsa.DecryptPKCS1v15`.
// Requires data to have been encrypted using the server public key.
func (kp *KeypairBytes) DecryptWithPrivateKey(data *[]byte) error {
	block, _ := x509.ParsePKCS1PrivateKey(kp.PrivateKey)
	decrypted, err := rsa.DecryptPKCS1v15(nil, block, *data)
	if err != nil {
		return err
	}

	*data = decrypted
	return nil
}
