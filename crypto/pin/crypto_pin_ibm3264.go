package pin

import (
	"bytes"
	"encoding/hex"
	"strings"
)

type PinblockIbm3264 struct {
	PinBlocker
}

func (pinBlock *PinblockIbm3264) Encrypt(pan string, clearPin string, key []byte) ([]byte, error) {

	var buf = bytes.NewBufferString(clearPin)
	for i := buf.Len(); i < 16; i++ {
		buf.WriteString("F")
	}

	pinBlockData, _ := hex.DecodeString(buf.String())
	encPinBlock, err := EncryptPinBlock(pinBlockData, key)
	return encPinBlock, err

}

func (pinBlock *PinblockIbm3264) GetPin(pan string, pinBlockData []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pinBlockData, key)
	pinBlockStr := hex.EncodeToString(clearPinBlock)
	//log.Printf("decrypted pin block =",pin_block_str)
	indexF := strings.Index(pinBlockStr, "f")
	res = pinBlockStr[:indexF]
	return

}
