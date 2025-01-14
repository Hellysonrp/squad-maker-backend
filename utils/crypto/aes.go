// https://github.com/gtank/cryptopasta/blob/master/encrypt.go

// cryptopasta - basic cryptography examples
//
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along
// with this software. If not, see // <http://creativecommons.org/publicdomain/zero/1.0/>.

// Provides symmetric authenticated encryption using 256-bit AES-GCM with a random nonce.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"reflect"

	"squad-maker/utils/env"
)

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
// func NewEncryptionKey() *[32]byte {
// 	key := [32]byte{}
// 	_, err := io.ReadFull(rand.Reader, key[:])
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &key
// }

var key [32]byte

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	if reflect.ValueOf(key).IsZero() {
		err = initKey()
		if err != nil {
			return nil, err
		}
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if reflect.ValueOf(key).IsZero() {
		err = initKey()
		if err != nil {
			return nil, err
		}
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func initKey() error {
	b, err := env.GetSecretKey("AES_KEY")
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return errors.New("key is not 32 bytes")
	}
	copy(key[:], b)
	return nil
}
