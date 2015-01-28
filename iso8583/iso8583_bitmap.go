package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
)

type BitMap struct {

	//primary secondary and tertiary bitmaps
	_1bmp *big.Int
	_2bmp *big.Int
	_3bmp *big.Int
}

func NewBitMap() *BitMap {

	bitmap := new(BitMap)
	bitmap._1bmp = big.NewInt(0)
	bitmap._2bmp = big.NewInt(0)
	bitmap._3bmp = big.NewInt(0)

	return (bitmap)
}

//check if the bit at position 'pos' is
//set as 1
func (bmp *BitMap) IsOn(pos int) bool {
	pos = pos - 1
	i_bmp, i_pos := bmp.get_bmp_and_pos(pos)
	if i_bmp.Bit(i_pos) == 1 {
		return true
	} else {
		return false
	}

}

func (bmp *BitMap) get_bmp_and_pos(pos int) (*big.Int, int) {

	var target_bmp *big.Int
	var i_pos int

	switch {
	case pos >= 0 && pos < 64:
		{
			target_bmp = bmp._1bmp
			i_pos = 63 - pos

		}
	case pos >= 64 && pos < 128:
		{

			target_bmp = bmp._2bmp
			i_pos = 127 - pos

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

func (bmp *BitMap) SetOff(pos int) {
	pos = pos - 1
	target_bmp, i_pos := bmp.get_bmp_and_pos(pos)
	target_bmp.SetBit(target_bmp, i_pos, 0)
}

func (bmp *BitMap) SetOn(pos int) {
	pos = pos - 1
	target_bmp, i_pos := bmp.get_bmp_and_pos(pos)
	target_bmp.SetBit(target_bmp, i_pos, 1)

}

func to_octet(in_data []byte) []byte {

	//fmt.Println("to_octet -", hex.EncodeToString(in_data))
	if len(in_data) != 8 {
		n_pads := 8 - len(in_data)
		with_pads := make([]byte, n_pads)
		with_pads = append(with_pads, in_data...)

		return with_pads
	}

	return in_data

}

func (bmp *BitMap) Parse(iso_msg *Iso8583Message, buf *bytes.Buffer) {

	if buf.Len() >= 8 {
		tmp := make([]byte, 8)
		_, err := buf.Read(tmp)
		iso_msg.handle_error(err)
		bmp._1bmp.SetBytes(tmp)

		if tmp[0]&0x80 == 0x80 {
			if buf.Len() >= 8 {
				tmp := make([]byte, 8)
				_, err := buf.Read(tmp)
				iso_msg.handle_error(err)
				bmp._2bmp.SetBytes(tmp)

				if tmp[0]&0x80 == 0x80 {
					if buf.Len() >= 8 {
						tmp := make([]byte, 8)
						_, err := buf.Read(tmp)

						iso_msg.handle_error(err)
						bmp._3bmp.SetBytes(tmp)

					} else {
						iso_msg.buffer_underflow_error("bitmap")
					}
				}

			} else {
				iso_msg.buffer_underflow_error("bitmap")
			}
		}

	}

	iso_msg.log.Printf("parsed bitmap: %s\n", hex.EncodeToString(bmp.Bytes()))

}

func (bmp *BitMap) Bytes() []byte {

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

func (bmp *BitMap) Copy() *BitMap {

	new_bmp := NewBitMap()

	new_bmp._1bmp = big.NewInt(0).Set(bmp._1bmp)
	new_bmp._2bmp = big.NewInt(0).Set(bmp._2bmp)
	new_bmp._3bmp = big.NewInt(0).Set(bmp._3bmp)
	return new_bmp

}
