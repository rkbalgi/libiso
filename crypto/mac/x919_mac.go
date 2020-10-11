package mac

import (
	"crypto/cipher"
	"crypto/des"
)

var zeroIV []byte = make([]byte, 8)

// GenerateMacX919 generate a X9.19 MAC (ISO9797 - Algorithm 3) with a double length key (ede2)  data will be zero padded if required
func GenerateMacX919(inMacData []byte, keyData []byte) ([]byte, error) {

	tripleDesKey := make([]byte, len(keyData))
	copy(tripleDesKey, keyData)

	if len(keyData) == 16 {
		//double length key
		tripleDesKey = append(tripleDesKey, keyData[0:8]...)
	}
	if tripleDesCipher, err := des.NewTripleDESCipher(tripleDesKey); err != nil {
		return nil, err
	} else {

		macData := make([]byte, len(inMacData))
		copy(macData, inMacData)

		if len(macData) < 8 || len(macData)%8 != 0 {
			pads := make([]byte, 8-(len(macData)%8))
			macData = append(macData, pads...)
		}

		result := make([]byte, 8)

		if len(macData) == 8 {
			//only a single block, so encrypt triple des cbc mode
			//with a 0x00 * 8 iv
			tripledesCbcEncryptor := cipher.NewCBCEncrypter(tripleDesCipher, zeroIV)
			tripledesCbcEncryptor.CryptBlocks(result, macData)
		} else {

			var err error
			//des-cbc encrypt all but the last block using the left half of
			//the double length key
			data := (macData)[0 : len(macData)-8]
			desCipher, err := des.NewCipher(keyData[0:8])
			if err != nil {
				return nil, err
			}
			cbcEncryptor := cipher.NewCBCEncrypter(desCipher, zeroIV)
			allButLast := make([]byte, len(data))
			cbcEncryptor.CryptBlocks(allButLast, data)

			//take the last block of encrypted data
			//and use that as a iv to encrypt the last
			//clear block with triple des
			tripledesCbcEncryptor := cipher.NewCBCEncrypter(tripleDesCipher, allButLast[len(allButLast)-8:])
			tripledesCbcEncryptor.CryptBlocks(result, (macData)[len(macData)-8:len(macData)])

		}
		return result, nil
	}

}
