package iso8583

import (
	"encoding/hex"
	"github.com/rkbalgi/go/encoding/ebcdic"
	_"fmt"
)

type FieldData struct {
	field_data []byte
	field_def  IsoField
}

func (fld_data *FieldData) SetData(value string) {

	//fmt.Println(*fld_data)
	
	
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
			return ebcdic2ascii(field_data.field_data)
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
