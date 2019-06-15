package pin

import (
	"encoding/hex"
	"testing"
)

func Test_Ibm3264Format_Encrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pinBlock PinBlocker = new(PinblockIbm3264)

	if encPinBlock, err := pinBlock.Encrypt("", "1274", key); err != nil {
		t.Error(err)
	} else {
		if hex.EncodeToString(encPinBlock) != "89fa441aa25ff0cc" {
			t.Errorf("%s!=%s", "89fa441aa25ff0cc", hex.EncodeToString(encPinBlock))
			t.Fail()
		}
	}

}

func Test_Ibm3264Format_Decrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	encPinBlock, _ := hex.DecodeString("89fa441aa25ff0cc")
	var pinBlock PinBlocker = new(PinblockIbm3264)

	pin, _ := pinBlock.GetPin("", encPinBlock, key)

	if pin != "1274" {
		t.Errorf("%s!=%s", "1274", pin)
		t.Fail()
	}

}

func Test_Iso0Format_Encrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pinBlock PinBlocker = new(PinBlock_Iso0)

	encPinBlock, _ := pinBlock.Encrypt("4111111111111111", "1234", key)

	if hex.EncodeToString(encPinBlock) != "6042012526a9c2e0" {
		t.Errorf("%s!=%s", "6042012526a9c2e0", hex.EncodeToString(encPinBlock))
		t.Fail()
	}

}

func Test_Iso0Format_Decrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	encPinBlock, _ := hex.DecodeString("6042012526a9c2e0")
	var pinBlock PinBlocker = new(PinBlock_Iso0)

	pin, _ := pinBlock.GetPin("4111111111111111", encPinBlock, key)

	if pin != "1234" {
		t.Errorf("%s!=%s", "1234", pin)
		t.Fail()
	}

}
func Test_Iso1Format(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pinBlock PinBlocker = new(PinblockIso1)

	encPinBlock, _ := pinBlock.Encrypt("", "2278", key)

	pin, _ := pinBlock.GetPin("", encPinBlock, key)

	if pin != "2278" {
		t.Errorf("%s!=%s", "2278", pin)
		t.Fail()
	}

}

func Test_Iso3Format(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pinBlock PinBlocker = new(PinblockIso3)

	encPinBlock, _ := pinBlock.Encrypt("4111111111111111", "22781", key)
	//t.Log(hex.EncodeToString(enc_pin_block))
	pin, _ := pinBlock.GetPin("4111111111111111", encPinBlock, key)

	if pin != "22781" {
		t.Errorf("%s!=%s", "2278", pin)
		t.Fail()
	}

}
