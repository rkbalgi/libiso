package iso8583

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Test_Bitmap(t *testing.T) {

	bmp := NewBitMap()

	bmp.SetOn(1)
	bmp.SetOn(64)
	bmp.SetOn(52)
	bmp.SetOn(40)
	bmp.SetOn(128)

	//bmp.SetBit(bmp, 0, 0)
	//bmp.SetBit(bmp, 4, 1)

	//t.Log(bmp.Uint64())

	t.Log(hex.EncodeToString(bmp.Bytes()))
}

// 1101 B
func Test_Iso8583Message(t *testing.T) {

	data, _ := hex.DecodeString("31313030F0040000E000000000000000000000013135333731313131313131313131313134F0F0F4F0F0F0303030303030303030313232313231320010e1e2e3a1a2a3a4d1d2d3a2a3a4d1d2d30010f1f2f3a1a2a3a4d1a2a3F0F2F8F3F7F1F1F1F1F1F1F1F1F1F1F1F1F47EF1F2F1F2F5F6F5F6F5F5F5F4e201f245ed4abb00")

	buf := bytes.NewBuffer(data)
	t.Log(hex.EncodeToString(buf.Bytes()))
	_, err := Handle(buf)
	if err != nil {
		t.Fail()
	}

}

