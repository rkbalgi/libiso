// this package handles various pin block formats
// pin blocks are described in good detail at
// http://www.paymentsystemsblog.com/2010/03/03/pin-block-formats/

package pin

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	_crypt "github.com/rkbalgi/libiso/crypto"
)

var ErrInvalidPinLength = errors.New("libiso: Invalid PIN length (cannot exceed 12)")

// PinBlocker represents a interface for types that can decrypt or encrypt a PIN block
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

// EncryptPinBlock encrypts a new pin block
func EncryptPinBlock(pinBlock []byte, key []byte) (result []byte, err error) {

	if len(key) == 8 {
		result, err = _crypt.EncryptDes(pinBlock, key)
	} else {
		result, err = _crypt.EncryptTripleDes(pinBlock, key)
	}

	return

}

// DecryptPinBlock decrypts a new PIN block
func DecryptPinBlock(pinBlock []byte, key []byte) (result []byte, err error) {

	if len(key) == 8 {
		result, err = _crypt.DecryptDes(pinBlock, key)
	} else {
		result, err = _crypt.DecryptTripleDes(pinBlock, key)
	}

	return

}
