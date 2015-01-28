/**
**/


package mac

import (
	"github.com/rkbalgi/go/crypto"
)

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

	result := crypto.EncryptDesCbc(mac_data, key_data)
	return (result[len(result)-8:])

}
