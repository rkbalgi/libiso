package pin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinblockIso3 struct {
	PinBlocker
}

func (pinBlock *PinblockIso3) Encrypt(pan12digits string, clearPin string, key []byte) (res []byte, err error) {

	if len(clearPin) > 12 {
		panic("pin length > 12")
	}

	buf := bytes.NewBufferString(fmt.Sprintf("3%X%s", len(clearPin), clearPin))

	//random pads
	fillRandom(buf)

	pinBlockDataA, err := hex.DecodeString(buf.String())
	pinBlockDataB, err := hex.DecodeString("0000" + pan12digits)

	for i, v := range pinBlockDataB {
		pinBlockDataA[i] = pinBlockDataA[i] ^ v
	}
	//log.Printf(" xor'ed pin block =", hex.EncodeToString(pin_block_data_a))

	res, err = EncryptPinBlock(pinBlockDataA, key)
	return

}

func (pinBlock *PinblockIso3) GetPin(pan12digits string, pinBlockData []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pinBlockData, key)

	//pan_12digits := pan[len(pan)-13 : len(pan)-1]
	pinBlockDataB, _ := hex.DecodeString("0000" + pan12digits)
	//log.Printf(" pin block (b) =", hex.EncodeToString(pin_block_data_b))

	for i, v := range pinBlockDataB {
		clearPinBlock[i] = clearPinBlock[i] ^ v
	}

	pinBlockStr := hex.EncodeToString(clearPinBlock)
	//log.Printf(" clear pin block (b) =", pin_block_str)

	nPinDigits, _ := strconv.ParseInt(pinBlockStr[1:2], 16, 16)
	res = pinBlockStr[2:(2 + nPinDigits)]

	return

}
