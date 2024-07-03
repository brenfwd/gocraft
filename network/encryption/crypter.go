package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type Crypter struct {
	decStream *cipher.Stream
	encStream *cipher.Stream
}

func NewCrypter(key []byte) (*Crypter, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("key passed to NewCrypter should be 16 bytes but received %d bytes", len(key))
	}

	// Note: Here we use the shared key for both the key and the IV.
	//       Not 100% sure this is correct/secure, but it *does* work...
	decBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// decStream := cipher.NewCFBDecrypter(decBlock, key[:decBlock.BlockSize()])
	decStream := NewCFB8Decrypter(decBlock, key[:decBlock.BlockSize()])

	encBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// encStream := cipher.NewCFBEncrypter(encBlock, key[:encBlock.BlockSize()])
	encStream := NewCFB8Encrypter(encBlock, key[:encBlock.BlockSize()])

	c := Crypter{
		decStream: &decStream,
		encStream: &encStream,
	}
	return &c, nil
}

// Encrypts data **in-place**.
func (c *Crypter) Encrypt(data *[]byte) {
	(*c.encStream).XORKeyStream(*data, *data)
}

// Decrypts data **in-place**.
func (c *Crypter) Decrypt(data *[]byte) {
	(*c.decStream).XORKeyStream(*data, *data)
}
