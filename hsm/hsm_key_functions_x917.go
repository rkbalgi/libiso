package hsm

import (
	"github.com/rkbalgi/go/crypto"
)

//encrypt a key under the kek using x917
func encryptKeyKekX917(keyStr string, kek []byte) ([]byte, error) {

	keyData := extractKeyData(keyStr)
	return encryptKeyX917(keyData, kek)
}

func encryptKeyX917(keyData []byte, kek []byte) ([]byte, error) {

	return crypto.EncryptTripleDes(keyData, kek)
}

//decrypt a key from under the kek using x917

func decryptKeyX917(keyData []byte, kek []byte) ([]byte, error) {

	return crypto.DecryptTripleDes(keyData, kek)
}
