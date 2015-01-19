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

	data, _ := hex.DecodeString("31313030D0040000000000000000000000000001313533373131313131313131313131313430303030303030303031323231323132e201f245ed4abb00")

	buf := bytes.NewBuffer(data)
	iso_msg, err := Handle(buf)
	t.Log(iso_msg.Dump());

	if err != nil {
		t.Fail()
	}
	//bmp.SetOn(1)
	//bmp.SetOn(64)
	//bmp.SetOn(52)
	//bmp.SetOn(40)
	//bmp.SetOn(128)

	//bmp.SetBit(bmp, 0, 0)
	//bmp.SetBit(bmp, 4, 1)

	//t.Log(bmp.Uint64())

	//t.Log(hex.EncodeToString(bmp.Bytes()))
}
