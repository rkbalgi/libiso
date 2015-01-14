package iso8583

import (
	"encoding/hex"
	"testing"
)


func Test_Bitmap(t *testing.T) {

	bmp := NewBitmap()

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
