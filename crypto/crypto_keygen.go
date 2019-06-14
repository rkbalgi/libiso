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
	//if keyLen == 8 {
	//	key = tmp[0:8]
	//} else if keyLen == 16 {
	//	key = tmp[:16]
	//} else {
	//	key = tmp[:24]
	//}

	toOddParity(key)

	return key, nil

}

/*
func GenerateDesKey(key_len int) []byte {

	if(key_len%8!=0){
		panic("invalid keylen for generation")
	}

	key := make([]byte, key_len)
	n, err := rand.Read(key)
	if n != key_len || err != nil {
		panic("key gen failure" + err.Error())
	}

	return key

}
*/
