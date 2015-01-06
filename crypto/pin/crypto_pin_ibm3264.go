package pin

import (
	"bytes"
	_ "crypto/rand"
	"encoding/hex"
	_ "fmt"
	_ "log"
	_ "strconv"
	"strings"
)

type PinBlock_Ibm3264 struct {
	PinBlocker
}

func (pin_block *PinBlock_Ibm3264) Encrypt(pan string, clear_pin string, key []byte) []byte {

	buf := bytes.NewBufferString(clear_pin)
	for i := buf.Len(); i < 16; i++ {
		buf.WriteString("F")
	}

	pin_block_data, _ := hex.DecodeString(buf.String())
	enc_pin_block := EncryptPinBlock(pin_block_data, key)
	return (enc_pin_block)

}

func (pin_block *PinBlock_Ibm3264) GetPin(pan string, pin_block_data []byte, key []byte) string {

	clear_pin_block := DecryptPinBlock(pin_block_data, key)
	pin_block_str := hex.EncodeToString(clear_pin_block)
	//log.Printf("decrypted pin block =",pin_block_str)
	index_f := strings.Index(pin_block_str, "f")
	return (pin_block_str[:index_f])

}
