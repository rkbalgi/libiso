package iso8583

import (
	"fmt"
	"math/big"
	"encoding/hex"
)

const (
	V1 = "V1"
	V0 = "V0"
)

type Bitmap struct {

	//primary secondary and tertiary bitmaps
	_1bmp *big.Int
	_2bmp *big.Int
	_3bmp *big.Int
}

func NewBitmap() *Bitmap {

	bitmap := new(Bitmap)
	bitmap._1bmp = big.NewInt(0)
	bitmap._2bmp = big.NewInt(0)
	bitmap._3bmp = big.NewInt(0)

	return (bitmap)
}

func (bmp *Bitmap) get_bmp_and_pos(pos int) (*big.Int, int) {

	var target_bmp *big.Int
	var i_pos int
	
	fmt.Println(pos);

	switch {
	case pos >= 0 && pos < 64:
		{
			target_bmp = bmp._1bmp
			i_pos = 63 - pos
			fmt.Println("i_pos= ",i_pos);

		}
	case pos >= 64 && pos < 128:
		{
			 
			target_bmp = bmp._2bmp
			i_pos = 127 - pos
			fmt.Println("i_pos= ",i_pos);

		}
	case pos >= 128 && pos < 192:
		{
			target_bmp = bmp._3bmp
			i_pos = 191 - pos
		}
	default:
		{
			panic(fmt.Sprint("invalid position in bitmap", pos))
		}

	}

	return target_bmp, i_pos

}

func (bmp *Bitmap) SetOff(pos int) {
	pos = pos - 1
	target_bmp, i_pos := bmp.get_bmp_and_pos(pos)
	target_bmp.SetBit(target_bmp, i_pos, 0)
}

func (bmp *Bitmap) SetOn(pos int) {
	pos = pos - 1
	target_bmp, i_pos := bmp.get_bmp_and_pos(pos)
	target_bmp.SetBit(target_bmp, i_pos, 1)

}

func to_octet(in_data []byte) []byte {

    fmt.Println("to_octet -",hex.EncodeToString(in_data)); 
	if len(in_data) != 8 {
		n_pads := 8 - len(in_data)
		with_pads := make([]byte, n_pads)
		with_pads = append(with_pads, in_data...)

		return with_pads
	}
	

	return in_data

}

func (bmp *Bitmap) Bytes() []byte {

	var bmp_data []byte

	bmp_data = to_octet(bmp._1bmp.Bytes())
	if bmp._1bmp.Bit(63) == 1 {
		tmp := to_octet(bmp._2bmp.Bytes())
		bmp_data = append(bmp_data, tmp...)
	}
	if bmp._2bmp.Bit(63) == 1 {
		tmp := to_octet(bmp._3bmp.Bytes())
		bmp_data = append(bmp_data, tmp...)
	}

	return bmp_data

}
