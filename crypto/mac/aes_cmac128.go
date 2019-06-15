package mac

//Implementation of RFC 4493 for cipher based macing (CMAC) using AES

//described here - https://tools.ietf.org/html/rfc4493
//and also  here http://csrc.nist.gov/publications/nistpubs/800-38B/SP_800-38B.pdf

import (
	_ "bytes"
	"crypto/aes"
	"crypto/cipher"
	_ "encoding/hex"
	_ "fmt"
	"log"
	_ "math/big"
)

var constZero = make([]byte, 16)
var constRb = make([]byte, 16)
var zeroIv = make([]byte, 16)

const (
	AesBlockSize = 16
)

func init() {
	constRb[len(constRb)-1] = 0x87

}

func leftShift(inData []byte) []byte {

	data := make([]byte, len(inData))
	copy(data, inData)

	var carry bool = false
	for i := len(data) - 1; i >= 0; i-- {

		data[i] = data[i] << 1

		if carry {
			data[i] = data[i] | 0x01
		}
		if inData[i]&0x80 == 0x80 {
			carry = true
		} else {
			carry = false
		}

	}
	return data

}

func addPadding(data []byte) []byte {
	nPads := 0

	if len(data) < 16 {
		nPads = 16 - len(data)
	} else if len(data) > 16 {
		nPads = len(data) % AesBlockSize
	}

	pads := make([]byte, nPads)
	pads[0] = 0x80

	data = append(data, pads...)

	return data

}

func subKeys(key []byte) ([]byte, []byte) {

	var err error
	l, err := aesEncrypt(key, constZero)
	if err != nil {
		log.Print(err)
	}
	k1, k2 := make([]byte, 16), make([]byte, 16)

	//is msb of l=0 then K1= l<<1 else K1= const_rb ^ (l<<1)
	if l[0]&0x80 == 0x80 {
		copy(k1, l)
		k1 = leftShift(k1)
		k1 = xor(k1, constRb)
	} else {
		copy(k1, l)
		k1 = leftShift(k1)
	}
	//is msb of k1=0 then K2= k1<<1 else K2= const_rb ^ (k1<<1)
	if k1[0]&0x80 == 0x80 {
		copy(k2, k1)
		k2 = leftShift(k2)
		k2 = xor(k2, constRb)
	} else {
		copy(k2, k1)
		leftShift(k2)
	}
	return k1, k2

}

//compute AES 128 CMAC using key and message M (in_data)
func AesCmac128(key []byte, inData []byte) []byte {

	k1, k2 := subKeys(key)

	var data []byte
	flag := false
	if len(inData) < AesBlockSize {
		data = addPadding(inData)
	} else if len(inData)%AesBlockSize != 0 {
		data = addPadding(inData)
	} else {
		//flag is true
		flag = true
		data = inData
	}
	var lastBlock []byte

	if flag {
		//a complete block was detected
		lastBlock = xor(data[len(data)-16:], k1)
	} else {
		//a block requiring padding was detected
		lastBlock = xor(data[len(data)-16:], k2)
	}

	aesBlock, _ := aes.NewCipher(key)
	aesEncryptor := cipher.NewCBCEncrypter(aesBlock, zeroIv)

	var encBlock []byte = make([]byte, 16)

	if len(data)/AesBlockSize > 1 {
		encBlock = make([]byte, len(data)-16)
		aesEncryptor.CryptBlocks(encBlock, data[:len(data)-16])
	} else {
		encBlock = make([]byte, 16)
	}

	y := xor(encBlock[len(encBlock)-16:], lastBlock)
	t := make([]byte, 16)
	aesBlock.Encrypt(t, y)
	return (t)

}

func xor(a []byte, b []byte) []byte {

	result := make([]byte, len(a))
	for i, v := range a {
		result[i] = v ^ b[i]
	}
	return result
}

func aesEncrypt(key []byte, data []byte) ([]byte, error) {

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encData := make([]byte, len(data))

	aesBlock.Encrypt(encData, data)
	return encData, nil

}
