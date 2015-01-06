package pin


import (
	"encoding/hex"
	"testing"
)

func Test_Ibm3264Format_Encrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pin_block PinBlocker=new(PinBlock_Ibm3264);
	
	enc_pin_block := pin_block.Encrypt("","1274", key)

	if hex.EncodeToString(enc_pin_block) != "89fa441aa25ff0cc" {
		t.Errorf("%s!=%s","89fa441aa25ff0cc",hex.EncodeToString(enc_pin_block))
		t.Fail()
	}

}

func Test_Ibm3264Format_Decrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	enc_pin_block,_:=hex.DecodeString("89fa441aa25ff0cc")
	var pin_block PinBlocker=new(PinBlock_Ibm3264);
	
	pin := pin_block.GetPin("",enc_pin_block,key)

	if pin != "1274" {
		t.Errorf("%s!=%s","1274",pin)
		t.Fail()
	}

}


func Test_Iso0Format_Encrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pin_block PinBlocker=new(PinBlock_Iso0);
	
	enc_pin_block := pin_block.Encrypt("4111111111111111","1234", key)

	if hex.EncodeToString(enc_pin_block) != "6042012526a9c2e0" {
		t.Errorf("%s!=%s","6042012526a9c2e0",hex.EncodeToString(enc_pin_block))
		t.Fail()
	}

}

func Test_Iso0Format_Decrypt(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	enc_pin_block,_:=hex.DecodeString("6042012526a9c2e0")
	var pin_block PinBlocker=new(PinBlock_Iso0);
	
	pin := pin_block.GetPin("4111111111111111",enc_pin_block,key)

	if pin != "1234" {
		t.Errorf("%s!=%s","1234",pin)
		t.Fail()
	}

}
func Test_Iso1Format(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pin_block PinBlocker=new(PinBlock_Iso1);
	
	enc_pin_block := pin_block.Encrypt("","2278", key)
	
	pin := pin_block.GetPin("",enc_pin_block,key)

	if pin != "2278" {
		t.Errorf("%s!=%s","2278",pin)
		t.Fail()
	}


}

func Test_Iso3Format(t *testing.T) {

	key, _ := hex.DecodeString("e0f4543f3e2a2c5ffc7e5e5a222e3e4d")
	var pin_block PinBlocker=new(PinBlock_Iso3);
	
	enc_pin_block := pin_block.Encrypt("4111111111111111","22781", key)
	//t.Log(hex.EncodeToString(enc_pin_block))
	pin := pin_block.GetPin("4111111111111111",enc_pin_block,key)

	if pin != "22781" {
		t.Errorf("%s!=%s","2278",pin)
		t.Fail()
	}


}



