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
	_ "math/big"
)

var const_zero = make([]byte, 16)
var const_rb = make([]byte, 16)
var zero_iv = make([]byte, 16)

const (
	AES_BLOCK_SIZE = 16
)

func init() {
	const_rb[len(const_rb)-1] = 0x87

}

func left_shift(in_data []byte) []byte {

	data := make([]byte, len(in_data))
	copy(data, in_data)

	var carry bool = false
	for i := len(data) - 1; i >= 0; i-- {

		data[i] = data[i] << 1

		if carry {
			data[i] = data[i] | 0x01
		}
		if in_data[i]&0x80 == 0x80 {
			carry = true
		} else {
			carry = false
		}

	}
	return data

}

func add_padding(data []byte) []byte {
	n_pads := 0

	if len(data) < 16 {
		n_pads = 16 - len(data)
	} else if len(data) > 16 {
		n_pads = len(data) % AES_BLOCK_SIZE
	}

	pads := make([]byte, n_pads)
	pads[0] = 0x80

	data = append(data, pads...)

	return data

}

func sub_keys(key []byte) ([]byte, []byte) {

	l := aes_encrypt(key, const_zero)
	k1, k2 := make([]byte, 16), make([]byte, 16)

	//is msb of l=0 then K1= l<<1 else K1= const_rb ^ (l<<1)
	if l[0]&0x80 == 0x80 {
		copy(k1, l)
		k1 = left_shift(k1)
		k1 = xor(k1, const_rb)
	} else {
		copy(k1, l)
		k1 = left_shift(k1)
	}
	//is msb of k1=0 then K2= k1<<1 else K2= const_rb ^ (k1<<1)
	if k1[0]&0x80 == 0x80 {
		copy(k2, k1)
		k2 = left_shift(k2)
		k2 = xor(k2, const_rb)
	} else {
		copy(k2, k1)
		left_shift(k2)
	}
	return k1, k2

}

//compute AES 128 CMAC using key and message M (in_data)
func AesCmac128(key []byte, in_data []byte) []byte {

	k1, k2 := sub_keys(key)

	var data []byte
	flag := false
	if len(in_data) < AES_BLOCK_SIZE {
		data = add_padding(in_data)
	} else if len(in_data)%AES_BLOCK_SIZE != 0 {
		data = add_padding(in_data)
	} else {
		//flag is true
		flag = true
		data = in_data
	}
	var last_block []byte

	if flag {
		//a complete block was detected
		last_block = xor(data[len(data)-16:], k1)
	} else {
		//a block requiring padding was detected
		last_block = xor(data[len(data)-16:], k2)
	}

	aes_block, _ := aes.NewCipher(key)
	aes_encryptor := cipher.NewCBCEncrypter(aes_block, zero_iv)

	var enc_block []byte = make([]byte, 16)

	if len(data)/AES_BLOCK_SIZE > 1 {
		enc_block = make([]byte, len(data)-16)
		aes_encryptor.CryptBlocks(enc_block, data[:len(data)-16])
	} else {
		enc_block = make([]byte, 16)
	}

	y := xor(enc_block[len(enc_block)-16:], last_block)
	t := make([]byte, 16)
	aes_block.Encrypt(t, y)
	return (t)

}

func xor(a []byte, b []byte) []byte {

	result := make([]byte, len(a))
	for i, v := range a {
		result[i] = v ^ b[i]
	}
	return result
}

func aes_encrypt(key []byte, data []byte) []byte {

	aes_block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	enc_data := make([]byte, len(data))

	aes_block.Encrypt(enc_data, data)
	return enc_data

}
