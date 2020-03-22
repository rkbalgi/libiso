package hsm

import (
	"fmt"
	"go/crypto/pin"
	"log"
)

var thalesPinfmtMap map[int]pin.PinBlocker

type pinHandler struct {
	pan          string
	inPinblkFmt  int
	outPinblkFmt int
	inKey        []byte
	outKey       []byte
	clearPin     string
}

func init() {
	thalesPinfmtMap = make(map[int]pin.PinBlocker, 5)
	thalesPinfmtMap[1] = new(pin.PinBlock_Iso0)
	thalesPinfmtMap[5] = new(pin.PinblockIso1)
	thalesPinfmtMap[47] = new(pin.PinblockIso3)
	thalesPinfmtMap[03] = new(pin.PinblockIbm3264)

	log.Printf("[%d] pin block formats registered.", len(thalesPinfmtMap))
}

func newPinHandler(pan string, inPinblkFmt int,
	outPinblkFmt int, inKey []byte, outKey []byte) *pinHandler {
	ph := new(pinHandler)
	ph.pan = pan
	ph.inPinblkFmt = inPinblkFmt
	ph.outPinblkFmt = outPinblkFmt
	ph.inKey = inKey
	ph.outKey = outKey

	return ph

}

func (ph *pinHandler) decryptAndExtractPin(inPinblk []byte) error {

	pinBlocker := thalesPinfmtMap[ph.inPinblkFmt]
	var err error
	if pinBlocker == nil {
		return fmt.Errorf("unsupported pin block format - %d", ph.inPinblkFmt)
	} else {
		ph.clearPin, err = pinBlocker.GetPin(ph.pan, inPinblk, ph.inKey)
		if err != nil {
			return err
		}

	}
	return err
}

func (ph *pinHandler) getClearPin() string {
	return ph.clearPin
}

func (ph *pinHandler) createPinBlock() ([]byte, error) {

	pinBlocker := thalesPinfmtMap[ph.outPinblkFmt]
	destPinBlock, err := pinBlocker.Encrypt(ph.pan, ph.clearPin, ph.outKey)
	if err != nil {
		return nil, err
	}
	return destPinBlock, nil

}

func (ph *pinHandler) translate(inPinBlock []byte) ([]byte, error) {

	if err := ph.decryptAndExtractPin(inPinBlock); err != nil {
		return nil, err
	}
	destPinBlock, err := ph.createPinBlock()
	return destPinBlock, err

}
