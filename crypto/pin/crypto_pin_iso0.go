package pin

import (
	"bytes"
	_ "crypto/rand"
	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinBlock_Iso0 struct {
	PinBlocker
}

func (pin_block *PinBlock_Iso0) Encrypt(pan_12digits string, clear_pin string, key []byte) ([]byte, error) {

	if len(clear_pin) > 12 {
		panic("pin length > 12")
	}

	buf := bytes.NewBufferString(fmt.Sprintf("0%X%s", len(clear_pin), clear_pin))
	for buf.Len() < 16 {
		buf.WriteString("F")
	}

	pinBlockDataA, _ := hex.DecodeString(buf.String())
	//log.Printf(" pin block (a) =", buf.String())

	//pan_12digits := pan[len(pan)-13 : len(pan)-1]
	pinBlockDataB, _ := hex.DecodeString("0000" + pan_12digits)
	//log.Printf(" pin block (b) =", hex.EncodeToString(pin_block_data_b))

	for i, v := range pinBlockDataB {
		pinBlockDataA[i] = pinBlockDataA[i] ^ v
	}
	//log.Printf(" xor'ed pin block =", hex.EncodeToString(pin_block_data_a))

	encPinBlock, err := EncryptPinBlock(pinBlockDataA, key)
	return encPinBlock, err

}

func (pin_block *PinBlock_Iso0) GetPin(pan_12digits string, pin_block_data []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pin_block_data, key)

	//pan_12digits := pan[len(pan)-13 : len(pan)-1]
	pinBlockDataB, _ := hex.DecodeString("0000" + pan_12digits)
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
