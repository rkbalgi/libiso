// this package handles various pin block formats
// pin blocks are described in good detail at
// http://www.paymentsystemsblog.com/2010/03/03/pin-block-formats/

package pin

import (
	"bytes"
	"encoding/hex"
	_ "fmt"
	_ "log"
	_ "strconv"
	_ "strings"
	//"math"
	"crypto/rand"
	_crypt "github.com/rkbalgi/go/crypto"
	//"github.com/rkbalgi/crypto/pin"
)

type PinBlocker interface {
	Encrypt(pan string, clear_pin string, key []byte) []byte
	GetPin(pan string, pin_block_data []byte, key []byte) string
}

func fill_random(buf *bytes.Buffer) {

	tmp := make([]byte, 1)
	for buf.Len() < 16 {
		rand.Read(tmp)
		buf.WriteString(hex.EncodeToString(tmp))
	}

	buf.Truncate(16)
}

func EncryptPinBlock(pin_block []byte, key []byte) []byte {

	var result []byte
	if len(key) == 8 {
		result = _crypt.EncryptDes(pin_block, key)
	} else {
		result = _crypt.EncryptTripleDes(pin_block, key)
	}

	return result

}

func DecryptPinBlock(pin_block []byte, key []byte) []byte {

	var result []byte
	if len(key) == 8 {
		result = _crypt.DecryptDes(pin_block, key)
	} else {
		result = _crypt.DecryptTripleDes(pin_block, key)
	}

	return result

}
