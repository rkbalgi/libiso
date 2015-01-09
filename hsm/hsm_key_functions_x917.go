package hsm

import (
	"github.com/rkbalgi/go/crypto"
)

//encrypt a key under the kek using x917
func encrypt_key_kek_x917(key_str string, kek []byte, key_type string) []byte {

	key_data := extract_key_data(key_str)
	return (encrypt_key_x917(key_data, kek, key_type))
}

func encrypt_key_x917(key_data []byte, kek []byte, key_type string) []byte {

	return (crypto.EncryptTripleDes(key_data, kek))
}

//decrypt a key from under the kek using x917
func decrypt_key_kek_x917(key_str string, kek []byte, key_type string) []byte {

	key_data := extract_key_data(key_str)
	return (decrypt_key_x917(key_data, kek, key_type))
}

func decrypt_key_x917(key_data []byte, kek []byte, key_type string) []byte {

	return (crypto.DecryptTripleDes(key_data, kek))
}
