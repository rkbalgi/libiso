package pin

import (
	"bytes"

	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinBlock_Iso1 struct {
	PinBlocker
}

func (pin_block *PinBlock_Iso1) Encrypt(pan string, clear_pin string, key []byte) []byte {

	buf := bytes.NewBufferString(fmt.Sprintf("1%X%s", len(clear_pin), clear_pin))
	fill_random(buf)

	//log.Println("block =", buf.String())

	pin_block_data, _ := hex.DecodeString(buf.String())
	enc_pin_block := EncryptPinBlock(pin_block_data, key)
	return (enc_pin_block)

}

func (pin_block *PinBlock_Iso1) GetPin(pan string, pin_block_data []byte, key []byte) string {

	clear_pin_block := DecryptPinBlock(pin_block_data, key)
	pin_block_str := hex.EncodeToString(clear_pin_block)

	n_pin_digits, _ := strconv.ParseInt(pin_block_str[1:2], 16, 16)
	clear_pin := pin_block_str[2:(2 + n_pin_digits)]

	return (clear_pin)

}
