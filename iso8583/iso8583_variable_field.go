package iso8583

import (
	"bytes"
	"encoding/binary"
	_ "encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"strconv"
)

type VariableFieldDef struct {
	name            string
	data_encoding   int
	length_encoding int
	len_ind_size    int //in bytes
}

//create a new fixed field definition
func NewVariableFieldDef(p_name string,
	p_len_encoding int,
	p_data_encoding int,
	p_len_ind_size int) *VariableFieldDef {

	field := new(VariableFieldDef)
	field.name = p_name
	field.data_encoding = p_data_encoding
	field.length_encoding = p_len_encoding
	field.len_ind_size = p_len_ind_size

	return field

}

func (field_def *VariableFieldDef) to_string(data []byte) string {
	return fmt.Sprintf("[%s] = [%s]", field_def.name, hex.EncodeToString(data))
}

func (field_def *VariableFieldDef) get_data_encoding() int {
	return field_def.data_encoding
}

func (field_def *VariableFieldDef) Parse(iso_msg *Iso8583Message, buf *bytes.Buffer) *FieldData {

	tmp := make([]byte, field_def.len_ind_size)
	n, err := buf.Read(tmp)
	if n != field_def.len_ind_size || err != nil {

		if n != field_def.len_ind_size {
			iso_msg.buffer_underflow_error(field_def.name)
		} else {
			iso_msg.field_parse_error(field_def.name, err)
		}

		iso_msg.field_parse_error(field_def.name, err)
	}

	var data_len uint64 = 0
	switch field_def.length_encoding {
	case ascii_encoding:
		{
			data_len, _ = strconv.ParseUint(string(tmp), 10, 64)
		}
	case ebcdic_encoding:
		{

			data_len, _ = strconv.ParseUint(ebcdic.EncodeToString(tmp), 10, 64)
		}
	case binary_encoding:
		{
			data_len, _ = strconv.ParseUint(hex.EncodeToString(tmp), 16, 64)
		}
	case bcd_encoding:
		{
			data_len, _ = strconv.ParseUint(hex.EncodeToString(tmp), 10, 64)
		}
	default:
		{
			panic("unsupported encoding")
		}
	}

	f_data := new(FieldData)
	f_data.field_def = field_def
	b_field_data := make([]byte, data_len)
	n, err = buf.Read(b_field_data)
	if uint64(n) != data_len || err != nil {

		if uint64(n) != data_len {
			iso_msg.buffer_underflow_error(field_def.name)
		} else {
			iso_msg.field_parse_error(field_def.name, err)
		}

	}
	f_data.field_data = b_field_data
	iso_msg.log.Printf("parsed: [%s]=[%s]", field_def.name, f_data.String())

	return f_data

}

//add the field data into buf as per the encoding
func (field_def *VariableFieldDef) Assemble(iso_msg *Iso8583Message, buf *bytes.Buffer) {

}

func (field_def *VariableFieldDef) IsFixed() bool {
	return false
}

func (field_def *VariableFieldDef) DataLength() int {
	//not applicable
	return -1
}

//return the length part of the variable field
//as a []byte slice
func (field_def *VariableFieldDef) EncodedLength(data_len int) []byte {
	//not applicable

	if field_def.len_ind_size > 4 &&
		(field_def.length_encoding == bcd_encoding || field_def.length_encoding == binary_encoding) {
		panic("[llvar] invalid length indicator size for bcd/binary - >4")
	}

	var ll []byte

	switch field_def.length_encoding {
	case binary_encoding:
		{
			switch field_def.len_ind_size {
			case 1:
				{
					ll = []byte{byte(data_len)}
				}
			case 2:
				{
					ll = make([]byte, 2)
					binary.BigEndian.PutUint16(ll, uint16(data_len))
				}
			case 4:
				{
					ll = make([]byte, 4)
					binary.BigEndian.PutUint32(ll, uint32(data_len))
				}
			default:
				{
					panic(fmt.Sprintf("[llvar] invalid length indicator size for binary field - %d", field_def.len_ind_size))
				}
			}

		}

	case bcd_encoding:
		{

			switch field_def.len_ind_size {
			case 1:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%0d", data_len))
				}
			case 2:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%04d", data_len))
				}
			case 4:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%08d", data_len))
				}
			default:
				{
					panic(fmt.Sprintf("[llvar] invalid length indicator size for binary field - %d", field_def.len_ind_size))
				}
			}

		}

	case ascii_encoding:
		{

			len_str := encoded_length_as_string(field_def.len_ind_size, data_len)
			ll = []byte(len_str)

		}

	case ebcdic_encoding:
		{

			len_str := encoded_length_as_string(field_def.len_ind_size, data_len)
			ll = ebcdic.Decode(len_str)

		}

	}

	return ll
}

func encoded_length_as_string(len_ind_size int, data_len int) string {

	var tmp string

	switch len_ind_size {
	case 1:
		{
			if data_len > 9 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", data_len))
			}
			tmp = fmt.Sprintf("%d", data_len)
		}
	case 2:
		{
			if data_len > 99 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", data_len))
			}
			tmp = fmt.Sprintf("%02d", data_len)
		}
	case 3:
		{
			if data_len > 999 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", data_len))
			}
			tmp = fmt.Sprintf("%03d", data_len)
		}
	case 4:
		{
			if data_len > 9999 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", data_len))
			}
			tmp = fmt.Sprintf("%04d", data_len)
		}
	default:
		{
			panic(fmt.Sprintf("[llvar] invalid length indicator size for  field - %d", len_ind_size))
		}
	}

	return tmp
}

func (field_def *VariableFieldDef) String() string {
	return field_def.name

}
