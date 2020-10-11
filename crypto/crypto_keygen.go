package crypto

import (
	"crypto/rand"

	"errors"
)

//596b0abc958120f62a7edc6896c99144

var ErrInvalidKeyLength = errors.New("invalid key length")

// Generates a DES key, keyLen should be specified in multiple of 8
func GenerateDesKey(keyLen int) ([]byte, error) {

	if keyLen%8 != 0 {
		return nil, ErrInvalidKeyLength
	}

	key := make([]byte, keyLen)
	_, err := rand.Read(key) // hex.DecodeString("596b0abc958120f62a7edc6896c991442a7edc6896c99144")
	if err != nil {
		return nil, err
	}
	toOddParity(key)

	return key, nil

}
