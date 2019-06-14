/**
**/

package mac

import (
	"encoding/hex"
	"github.com/rkbalgi/go/crypto"
)

//generate a X9.9 MAC using a single length key
//data will be zero padded if required

func GenerateMacX99(inMacData []byte, keyData []byte) []byte {

	macData := make([]byte, len(inMacData))
	copy(macData, inMacData)

	println(hex.EncodeToString(inMacData), hex.EncodeToString(macData))
	println(len(macData))

	//add 0 padding
	if len(macData) < 8 || len(macData)%8 != 0 {
		pads := make([]byte, 8-(len(macData)%8))
		println("pads ", len(pads))
		macData = append(macData, pads...)
	}

	result := crypto.EncryptDesCbc(macData, keyData)
	return result[len(result)-8:]

}
