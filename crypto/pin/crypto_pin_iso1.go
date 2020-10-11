package pin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinblockIso1 struct {
	PinBlocker
}

func (pinBlock *PinblockIso1) Encrypt(pan string, clearPin string, key []byte) (res []byte, err error) {

	if len(clearPin) > 12 {
		return nil, ErrInvalidPinLength
	}
	buf := bytes.NewBufferString(fmt.Sprintf("1%X%s", len(clearPin), clearPin))
	fillRandom(buf)

	pinBlockData, err := hex.DecodeString(buf.String())
	res, err = EncryptPinBlock(pinBlockData, key)
	return

}

func (pinBlock *PinblockIso1) GetPin(pan string, pinBlockData []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pinBlockData, key)
	pinBlockStr := hex.EncodeToString(clearPinBlock)

	nPinDigits, _ := strconv.ParseInt(pinBlockStr[1:2], 16, 16)
	res = pinBlockStr[2:(2 + nPinDigits)]

	return

}
