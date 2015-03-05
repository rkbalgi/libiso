package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
)

type BitMap struct {

	//primary secondary and tertiary bitmaps
	_1bmp          *big.Int
	_2bmp          *big.Int
	_3bmp          *big.Int
	sub_field_def  []IsoField
	sub_field_data []*FieldData
	id             int
	iso_msg_def    *Iso8583MessageDef
}

func NewBitMap() *BitMap {

	bitmap := new(BitMap)
	bitmap._1bmp = big.NewInt(0)
	bitmap._2bmp = big.NewInt(0)
	bitmap._3bmp = big.NewInt(0)
	bitmap.sub_field_def = make([]IsoField, 128+1)
	bitmap.sub_field_data = make([]*FieldData, 128+1)

	return (bitmap)
}

type BitmappedField interface {
	Parse(iso_msg *Iso8583Message, buf *bytes.Buffer)
	IsOn(int) bool
	SetOn(int)
	SetOff(int)
	SetId(int)
	GetId() int
	SetSpec(iso_msg_def *Iso8583MessageDef)
	Bytes() []byte
	Def() string
}


func (bmp *BitMap) Def() string{
	return fmt.Sprintf("Bitmap");
}

//add a new  fixed field to the bitmap
func (bmp *BitMap) add_fixed_field(
	p_bmp_pos int,
	p_name string,
	p_data_encoding int,
	p_field_len int) {

	field := NewFixedFieldDef(p_name, p_data_encoding, p_field_len)
	field.SetId(bmp.iso_msg_def.next_field_seq())
	bmp.sub_field_def[p_bmp_pos] = field

}

//add a new  variable field to the bitmap
func (bmp *BitMap) add_variable_field(p_bmp_pos int, p_name string,
	p_len_encoding int,
	p_data_encoding int,
	p_field_len int) {

	field := NewVariableFieldDef(p_name, p_len_encoding, p_data_encoding, p_field_len)
	field.SetId(bmp.iso_msg_def.next_field_seq())
	bmp.sub_field_def[p_bmp_pos] = field

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

	if bmp._2bmp.Uint64() != 0 {
		//second bmp exists
		bmp._1bmp.SetBit(bmp._1bmp, 63, 1)
	} else {
		bmp._1bmp.SetBit(bmp._1bmp, 63, 0)
	}
	if bmp._3bmp.Uint64() != 0 {
		//second bmp exists
		bmp._2bmp.SetBit(bmp._2bmp, 63, 1)
	} else {
		bmp._2bmp.SetBit(bmp._2bmp, 63, 0)
	}
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

func (bmp *BitMap) bit_string() string{
	
	buf:=bytes.NewBufferString("");
	for i:=1;i<129;i++{
		if bmp.IsOn(i){
			buf.Write([]byte("1"));
		}else{
			buf.Write([]byte("0"));
		}
	}
	
	return buf.String();
	
}

func (bmp *BitMap) copy_bits(src_bmp *BitMap) {

	bmp._1bmp = big.NewInt(0).Set(src_bmp._1bmp)
	bmp._2bmp = big.NewInt(0).Set(src_bmp._2bmp)
	bmp._3bmp = big.NewInt(0).Set(src_bmp._3bmp)

}

func (bmp *BitMap) SetId(id int) {
	bmp.id = id
}

func (bmp *BitMap) GetId() int {
	return bmp.id
}

func (bmp *BitMap) SetSpec(iso_msg_def *Iso8583MessageDef) {
	bmp.iso_msg_def = iso_msg_def
}
