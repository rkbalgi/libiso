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

var ZeroIv []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0}

//Performs TripleDes encryption of data with the key

func EncryptTripleDesEde2(key []byte, data []byte, paddingType PaddingType) ([]byte, error) {

	var keyData []byte
	keyData = append(keyData, key...)
	keyData = append(keyData, key[:8]...)
	//add padding
	newData := paddingType.Pad(data)

	log.Printf("(+Pad) %s ", hex.EncodeToString(newData))

	block, err := des.NewTripleDESCipher(keyData)
	if err != nil {
		log.Printf("failed to encrypt - %s", err.Error())
		return nil, err
	}

	encryptedData := make([]byte, len(newData))
	for i := 0; i < len(data); i += 8 {
		block.Encrypt(encryptedData[i:i+8], newData[i:i+8])
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
		log.Printf("failed to encrypt - %s", err.Error())
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

func EncryptDes(data []byte, key []byte) ([]byte, error) {
	return encryptDes(data, key)
}

func DecryptDes(data []byte, key []byte) ([]byte, error) {
	return decryptDes(data, key)
}

func EncryptTripleDes(data []byte, key []byte) ([]byte, error) {
	return encrypt3des(data, key)
}

func DecryptTripleDes(data []byte, key []byte) ([]byte, error) {
	return decrypt3des(data, key)
}

//single des

func decryptDes(data []byte, key []byte) ([]byte, error) {

	desCipher, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		desCipher.Decrypt(result[i:], data[i:])
	}
	return result, nil

}

func encryptDes(data []byte, key []byte) ([]byte, error) {

	desCipher, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		desCipher.Encrypt(result[i:], data[i:])
	}
	return result, nil

}

func EncryptDesCbc(data []byte, key []byte) ([]byte, error) {
	return encryptDesCbc(data, key)
}

func encryptDesCbc(data []byte, key []byte) ([]byte, error) {

	desCipher, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(desCipher, ZeroIv)
	result := make([]byte, len(data))
	encrypter.CryptBlocks(result, data)
	return result, nil

}

//triple des

func decrypt3des(data []byte, key []byte) ([]byte, error) {

	tripleDESKey := make([]byte, 0)

	if len(key) == 16 {
		//ede2
		tripleDESKey = append(tripleDESKey, key...)
		tripleDESKey = append(tripleDESKey, key[:8]...)
	} else {
		//ede3
		tripleDESKey = append(tripleDESKey, key...)
	}
	tripleDESCipher, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return nil, err
	}

	result := make([]byte, len(data))
	//log.Printf("%s %s", hex.EncodeToString(result), hex.EncodeToString(data))
	for i := 0; i < len(data); i += 8 {
		tripleDESCipher.Decrypt(result[i:], data[i:])
	}
	return result, err

}

func encrypt3des(data []byte, key []byte) ([]byte, error) {

	tripleDESKey := make([]byte, 0)

	if len(key) == 16 {
		//ede2
		tripleDESKey = append(tripleDESKey, key...)
		tripleDESKey = append(tripleDESKey, key[:8]...)
	} else {
		//ede3
		tripleDESKey = append(tripleDESKey, key...)
	}

	log.Println(hex.EncodeToString(tripleDESKey))
	tripleDESCipher, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return nil, err
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += 8 {
		tripleDESCipher.Encrypt(result[i:], data[i:])
	}
	//tripleDESCipher.Encrypt(result, data)

	return result, nil

}
