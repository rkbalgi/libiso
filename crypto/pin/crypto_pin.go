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
	Encrypt(pan string, clearPin string, key []byte) ([]byte, error)
	GetPin(pan string, pinBlockData []byte, key []byte) (string, error)
}

func fillRandom(buf *bytes.Buffer) {

	tmp := make([]byte, 1)
	for buf.Len() < 16 {
		_, _ = rand.Read(tmp)
		buf.WriteString(hex.EncodeToString(tmp))
	}

	buf.Truncate(16)
}

func EncryptPinBlock(pinBlock []byte, key []byte) (result []byte, err error) {

	if len(key) == 8 {
		result, err = _crypt.EncryptDes(pinBlock, key)
	} else {
		result, err = _crypt.EncryptTripleDes(pinBlock, key)
	}

	return

}

func DecryptPinBlock(pinBlock []byte, key []byte) (result []byte, err error) {

	if len(key) == 8 {
		result, err = _crypt.DecryptDes(pinBlock, key)
	} else {
		result, err = _crypt.DecryptTripleDes(pinBlock, key)
	}

	return

}
