package iso8583

import (
	"encoding/hex"
	_ "fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"log"
)

type FieldData struct {
	field_data []byte
	field_def  IsoField
	bmp_def    *BitMap
}

func (fld_data *FieldData) BitmapDef() *BitMap {
	return fld_data.bmp_def
}

func (fld_data *FieldData) Def() IsoField {
	return fld_data.field_def
}

//SetData sets field data as per the encoding
//additional padding will be applied if required
func (fld_data *FieldData) SetData(value string) {

	switch fld_data.field_def.get_data_encoding() {
	case ascii_encoding:
		{
			switch fld_data.field_def.(type) {
			case *FixedFieldDef:
				{
					data := []byte(value)
					fld_data.set_truncate_pad(data)
					break
				}
			default:
				{
					fld_data.field_data = []byte(value)
				}
			}

		}
	case ebcdic_encoding:
		{
			data := ebcdic.Decode(value)
			switch fld_data.field_def.(type) {
			case *FixedFieldDef:
				{
					fld_data.set_truncate_pad(data)
					break
				}
			default:
				{
					fld_data.field_data = data
				}
			}

		}
	case binary_encoding:
		fallthrough
	case bcd_encoding:
		{
			var err error

			data, err := hex.DecodeString(value)
			if err != nil {
				panic(err.Error())
			}
			switch fld_data.field_def.(type) {
			case *FixedFieldDef:
				{
					fld_data.set_truncate_pad(data)
					break
				}
			default:
				{
					fld_data.field_data = data
				}
			}

		}
	default:
		{
			panic("unsupported encoding")
		}

	}

}

func (fld_data *FieldData) set_truncate_pad(data []byte) {

	def_obj := fld_data.field_def.(*FixedFieldDef)
	pad_byte := byte(0x00)
	switch def_obj.get_data_encoding() {
	case ascii_encoding:
		pad_byte = 0x20
	case ebcdic_encoding:
		pad_byte = 0x40
	}

	if len(data) == def_obj.data_size {
		fld_data.field_data = data
	} else if len(data) > def_obj.data_size {
		//truncate
		fld_data.field_data = data[:len(data)]
	} else {

		fld_data.field_data = make([]byte, def_obj.data_size)
		for i, _ := range fld_data.field_data {
			fld_data.field_data[i] = pad_byte
		}
		copy(fld_data.field_data, data)
	}
}

//make a copy of FieldData
func (fld_data *FieldData) copy() *FieldData {

	new_fld_data := new(FieldData)
	new_fld_data.field_data = make([]byte, len(fld_data.field_data))
	copy(new_fld_data.field_data, fld_data.field_data)
	new_fld_data.field_def = fld_data.field_def

	return new_fld_data
}

func (field_data FieldData) String() string {

	if field_data.bmp_def != nil {
		return hex.EncodeToString(field_data.bmp_def.Bytes())
	}

	switch field_data.field_def.get_data_encoding() {
	case ascii_encoding:
		{
			return string(field_data.field_data)
		}
	case ebcdic_encoding:
		{
			encoded := ebcdic.EncodeToString(field_data.field_data)
			log.Println("encoded - ", encoded, "hex ", hex.EncodeToString(field_data.field_data))
			return encoded
		}
	case binary_encoding:
		fallthrough
	case bcd_encoding:
		{
			return hex.EncodeToString(field_data.field_data)
		}
	default:
		{
			panic("unsupported encoding")
		}

	}

}

//return the raw data associated with this field
//it will also include any ll portions for a variable field
func (field_data FieldData) Bytes() []byte {

	if field_data.bmp_def != nil {
		//if it's a bmp field, just return the data
		return field_data.bmp_def.Bytes()
	}

	if field_data.field_def.IsFixed() {
		data_len := field_data.field_def.DataLength()
		if len(field_data.field_data) > data_len {
			log.Printf("Warning: field [%s] length exceeds defined length, will be truncated")
			return field_data.field_data[0:data_len]
		} else if len(field_data.field_data) < data_len {
			//add default padding
			new_fld_data := make([]byte, data_len)
			copy(new_fld_data, field_data.field_data)
			return new_fld_data
		}
		return field_data.field_data[0:data_len]

	} else {
		//variable fields should have length indicators
		data_len := len(field_data.field_data)
		ll := field_data.field_def.EncodedLength(data_len)
		ll_data := append(ll, field_data.field_data...)
		return ll_data

	}

}
