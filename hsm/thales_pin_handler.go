package hsm

import (
	"fmt"
	"github.com/rkbalgi/go/crypto/pin"
	"log"
)

var thales_pinfmt_map map[int]pin.PinBlocker

type pin_handler struct {
	pan            string
	in_pinblk_fmt  int
	out_pinblk_fmt int
	in_key         []byte
	out_key        []byte

	clear_pin string
}

func init() {
	thales_pinfmt_map = make(map[int]pin.PinBlocker, 5)
	thales_pinfmt_map[1] = new(pin.PinBlock_Iso0)
	thales_pinfmt_map[5] = new(pin.PinBlock_Iso1)
	thales_pinfmt_map[47] = new(pin.PinBlock_Iso3)
	thales_pinfmt_map[03] = new(pin.PinBlock_Ibm3264)

	log.Printf("[%d] pin block formats registered.", len(thales_pinfmt_map))
}

func new_pin_handler(pan string, in_pinblk_fmt int, out_pinblk_fmt int, in_key []byte, out_key []byte) *pin_handler {
	ph := new(pin_handler)
	ph.pan = pan
	ph.in_pinblk_fmt = in_pinblk_fmt
	ph.out_pinblk_fmt = out_pinblk_fmt
	ph.in_key = in_key
	ph.out_key = out_key

	return ph

}

func (ph *pin_handler) decrypt_and_extract_pin(in_pinblk []byte) {

	pin_blocker := thales_pinfmt_map[ph.in_pinblk_fmt]

	if pin_blocker == nil {
		panic(fmt.Sprintf("unsupported pin block format - %d", ph.in_pinblk_fmt))
	} else {
		ph.clear_pin = pin_blocker.GetPin(ph.pan, in_pinblk, ph.in_key)
	}
}

func (ph *pin_handler) get_clear_pin() string {
	return ph.clear_pin
}

func (ph *pin_handler) create_pin_block() []byte {

	pin_blocker := thales_pinfmt_map[ph.out_pinblk_fmt]
	dest_pin_block := pin_blocker.Encrypt(ph.pan, ph.clear_pin, ph.out_key)
	return (dest_pin_block)

}

func (ph *pin_handler) translate(in_pin_block []byte) []byte {

	ph.decrypt_and_extract_pin(in_pin_block)
	dest_pin_block := ph.create_pin_block()
	return (dest_pin_block)

}
