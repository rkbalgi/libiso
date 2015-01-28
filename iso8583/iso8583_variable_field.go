package iso8583

import (
	"bytes"
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

func (field_def *VariableFieldDef) String() string {
	return field_def.name

}
