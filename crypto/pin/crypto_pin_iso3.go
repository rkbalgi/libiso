package pin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

type PinblockIso3 struct {
	PinBlocker
}

func (pinBlock *PinblockIso3) Encrypt(pan string, clearPin string, key []byte) (res []byte, err error) {

	if len(clearPin) > 12 {
		return nil, ErrInvalidPinLength
	}

	buf := bytes.NewBufferString(fmt.Sprintf("3%X%s", len(clearPin), clearPin))

	//random pads
	fillRandom(buf)

	pinBlockDataA, err := hex.DecodeString(buf.String())
	pan12digits := pan
	if len(pan) != 12 {
		pan12digits = pan[len(pan)-13 : len(pan)-1]
	}

	pinBlockDataB, err := hex.DecodeString("0000" + pan12digits)

	for i, v := range pinBlockDataB {
		pinBlockDataA[i] = pinBlockDataA[i] ^ v
	}

	res, err = EncryptPinBlock(pinBlockDataA, key)
	return

}

func (pinBlock *PinblockIso3) GetPin(pan string, pinBlockData []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pinBlockData, key)

	pan12digits := pan[len(pan)-13 : len(pan)-1]
	pinBlockDataB, _ := hex.DecodeString("0000" + pan12digits)

	for i, v := range pinBlockDataB {
		clearPinBlock[i] = clearPinBlock[i] ^ v
	}

	pinBlockStr := hex.EncodeToString(clearPinBlock)

	nPinDigits, _ := strconv.ParseInt(pinBlockStr[1:2], 16, 16)
	res = pinBlockStr[2:(2 + nPinDigits)]

	return

}
