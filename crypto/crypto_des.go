package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"log"
)

type DesMode int

const (
	Ecb DesMode = iota + 1
	Cbc
)

var ZERO_IV []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0}

//Performs TripleDes encryption of data with the key

func EncryptTripleDesEde2(key []byte, data []byte, paddingType PaddingType) ([]byte, error) {

	var keyData []byte
	keyData = append(keyData, key...)
	keyData = append(keyData, key[:8]...)
	//add padding
	new_data := paddingType.Pad(data)

	log.Printf("(+Pad) %s ", hex.EncodeToString(new_data))

	block, err := des.NewTripleDESCipher(keyData)
	if err != nil {
		log.Printf("failed to encrypt -", err.Error())
		return nil, err
	}

	encryptedData := make([]byte, len(new_data))
	for i := 0; i < len(data); i += 8 {
		block.Encrypt(encryptedData[i:i+8], new_data[i:i+8])
	}

	return encryptedData, err

}

//performs Triple DES decryption

func DecryptTripleDesEde2(key []byte, data []byte, paddingType PaddingType) ([]byte, error) {

	var keyData []byte
	keyData = append(keyData, key...)
	keyData = append(keyData, key[:8]...)

	block, err := des.NewTripleDESCipher(keyData)
	if err != nil {
		log.Printf("failed to encrypt -", err.Error())
		return nil, err
	}

	decryptedData := make([]byte, len(data))
	for i := 0; i < len(data); i += 8 {
		block.Decrypt(decryptedData[i:i+8], data[i:i+8])
	}

	log.Printf("Decrypted Data %s\n", hex.EncodeToString(decryptedData))
	decryptedData = paddingType.RemovePad(decryptedData)

	return decryptedData, err

}

//exported wrappers to des and triple des routines

func EncryptDes(data []byte, key []byte) []byte {
	return encrypt_des(data, key)
}

func DecryptDes(data []byte, key []byte) []byte {
	return decrypt_des(data, key)
}

func EncryptTripleDes(data []byte, key []byte) []byte {
	return encrypt_3des(data, key)
}

func DecryptTripleDes(data []byte, key []byte) []byte {
	return decrypt_3des(data, key)
}

//single des

func decrypt_des(data []byte, key []byte) []byte {

	_des_cipher, err := des.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		_des_cipher.Decrypt(result[i:], data[i:])
	}
	return (result)

}

func encrypt_des(data []byte, key []byte) []byte {

	_des_cipher, err := des.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		_des_cipher.Encrypt(result[i:], data[i:])
	}
	return (result)

}

func EncryptDesCbc(data []byte, key []byte) []byte {
	return encrypt_des_cbc(data, key)
}

func encrypt_des_cbc(data []byte, key []byte) []byte {

	_des_cipher, err := des.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	encrypter := cipher.NewCBCEncrypter(_des_cipher, ZERO_IV)
	result := make([]byte, len(data))
	encrypter.CryptBlocks(result, data)
	return (result)

}

//triple des

func decrypt_3des(data []byte, key []byte) []byte {

	_3deskey := make([]byte, 0)

	if len(key) == 16 {
		//ede2
		_3deskey = append(_3deskey, key...)
		_3deskey = append(_3deskey, key[:8]...)
	} else {
		//ede3
		_3deskey = append(_3deskey, key...)
	}
	_3des_cipher, err := des.NewTripleDESCipher(_3deskey)
	if err != nil {
		panic(err.Error())
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		_3des_cipher.Decrypt(result[i:], data[i:])
	}
	return (result)

}

func encrypt_3des(data []byte, key []byte) []byte {

	_3deskey := make([]byte, 0)

	if len(key) == 16 {
		//ede2
		_3deskey = append(_3deskey, key...)
		_3deskey = append(_3deskey, key[:8]...)
	} else {
		//ede3
		_3deskey = append(_3deskey, key...)
	}

	log.Println(hex.EncodeToString(_3deskey))
	_3des_cipher, err := des.NewTripleDESCipher(_3deskey)
	if err != nil {
		panic(err.Error())
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += 8 {
		_3des_cipher.Encrypt(result[i:], data[i:])
	}
	//_3des_cipher.Encrypt(result, data)

	return (result)

}
