package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type FixedFieldDef struct {
	name          string
	data_encoding int
	data_size     int //in bytes
}

//create a new fixed field definition
func NewFixedFieldDef(p_name string, p_data_encoding int, p_field_len int) *FixedFieldDef {
	field := new(FixedFieldDef)
	field.name = p_name
	field.data_encoding = p_data_encoding
	field.data_size = p_field_len

	return field

}

//parse and return field data by reading appropriate bytes
//from the buffer buf
func (field_def *FixedFieldDef) Parse(iso_msg *Iso8583Message, buf *bytes.Buffer) *FieldData {

	tmp := make([]byte, field_def.data_size)
	n, err := buf.Read(tmp)
	
	if n != field_def.data_size || err != nil {
		if n != field_def.data_size {
			iso_msg.buffer_underflow_error(field_def.name)
		} else {
			iso_msg.field_parse_error(field_def.name, err)
		}
	}

	f_data := new(FieldData)
	f_data.field_def = field_def
	f_data.field_data = tmp
	iso_msg.log.Printf("parsed: [%s]=[%s]", field_def.name, f_data.String())

	return f_data

}

func (field_def *FixedFieldDef) to_string(data []byte) string {
	return fmt.Sprintf("[%s] = [%s]", field_def.name, hex.EncodeToString(data))
}

func (field_def *FixedFieldDef) get_data_encoding() int {
	return field_def.data_encoding
}

//add the field data into buf as per the encoding
func (field_def *FixedFieldDef) Assemble(iso_msg *Iso8583Message, buf *bytes.Buffer) {

	//buf.Write(iso_msg);

}

func (field_def *FixedFieldDef) String() string {
	return field_def.name
}
