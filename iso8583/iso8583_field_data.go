package iso8583

import (
	"encoding/hex"
	_"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"log"
)

type FieldData struct {
	field_data []byte
	field_def  IsoField
}

func (fld_data *FieldData) SetData(value string) {

	switch fld_data.field_def.get_data_encoding() {
	case ascii_encoding:
		{
			fld_data.field_data = []byte(value)

		}
	case ebcdic_encoding:
		{
			fld_data.field_data = ebcdic.Decode(value)
		}
	case binary_encoding:
		fallthrough
	case bcd_encoding:
		{
			var err error

			fld_data.field_data, err = hex.DecodeString(value)
			if err != nil {
				panic(err.Error())
			}
		}
	default:
		{
			panic("unsupported encoding")
		}

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

	switch field_data.field_def.get_data_encoding() {
	case ascii_encoding:
		{
			return string(field_data.field_data)
		}
	case ebcdic_encoding:
		{
			return ebcdic.EncodeToString(field_data.field_data)
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
