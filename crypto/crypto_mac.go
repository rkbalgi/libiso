package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"runtime"
	"log"
	"os"
)

var _log *log.Logger=log.New(os.Stdout,"## mac routines ## ",log.LstdFlags)

var __zero_iv []byte = make([]byte, 8)

//generate a X9.9 MAC using a single length key
//data will be zero padded if required

func GenerateMac_X99(in_mac_data []byte, key_data []byte) []byte {

	mac_data := make([]byte, len(in_mac_data))
	copy(mac_data, in_mac_data)

	//add 0 padding
	if len(mac_data) < 8 || len(mac_data)%8 != 0 {
		pads := make([]byte, 8-(len(mac_data)%8))
		mac_data = append(mac_data, pads...)
	}

	result := encrypt_des_cbc(mac_data, key_data)
	return (result[len(result)-8:])

}

//generate a X9.19 MAC (ISO9797 - Algorithm 3) with a double length key (ede2)
//data will be zero padded if required

func GenerateMac_X919(in_mac_data []byte, key_data []byte) []byte {

	_3deskey := make([]byte, len(key_data))
	copy(_3deskey, key_data)

	if len(key_data) == 16 {
		//double length key
		_3deskey = append(_3deskey, key_data[0:8]...)
	}
	triple_des_cipher, err := des.NewTripleDESCipher(_3deskey)
	check_error(err)

	mac_data := make([]byte, len(in_mac_data))
	copy(mac_data, in_mac_data)

	if len(mac_data) < 8 || len(mac_data)%8 != 0 {
		pads := make([]byte, 8-(len(mac_data)%8))
		mac_data = append(mac_data, pads...)
	}

	result := make([]byte, 8)

	if len(mac_data) == 8 {
		//only a single block, so encrypt triple des cbc mode
		//with a 0x00 * 8 iv
		tripledes_cbc_encryptor := cipher.NewCBCEncrypter(triple_des_cipher, __zero_iv)
		tripledes_cbc_encryptor.CryptBlocks(result, (mac_data))
	} else {

		//des-cbc encrypt all but the last block using the left half of
		//the double length key
		data := (mac_data)[0 : len(mac_data)-8]
		des_cipher, err := des.NewCipher(key_data[0:8])
		check_error(err)
		cbc_encryptor := cipher.NewCBCEncrypter(des_cipher, __zero_iv)
		all_but_last := make([]byte, len(data))
		cbc_encryptor.CryptBlocks(all_but_last, data)

		//take the last block of encrypted data
		//and use that as a iv to encrypt the last
		//clear block with triple des
		tripledes_cbc_encryptor := cipher.NewCBCEncrypter(triple_des_cipher, all_but_last[len(all_but_last)-8:len(all_but_last)])
		tripledes_cbc_encryptor.CryptBlocks(result, (mac_data)[len(mac_data)-8:len(mac_data)])

	}

	return (result)

}

func check_error(err error) {
	if err != nil {
		var stack_data []byte
		runtime.Stack(stack_data,false);
		_log.Printf(string(stack_data))
		
		panic(err.Error())
	}
}
