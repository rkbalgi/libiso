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
	id            int
	b_pos int //bit position in bitmap
}

//create a new fixed field definition
func NewFixedFieldDef(p_name string, p_data_encoding int, p_field_len int) *FixedFieldDef {
	field := new(FixedFieldDef)
	field.name = p_name
	field.data_encoding = p_data_encoding
	field.data_size = p_field_len

	return field

}

func (field *FixedFieldDef) Def() string {
	return fmt.Sprintf("Name: %-40s ; Id: %04d ; Type: %-10s ; Length: %04d ;Encoding: [%-10s]",
		field.name, field.GetId(), "Fixed", field.data_size, get_encoding_type(field.data_encoding))
}

//parse and return field data by reading appropriate bytes
//from the buffer buf
func (field_def *FixedFieldDef) Parse(
	iso_msg *Iso8583Message,
	f_data *FieldData,
	buf *bytes.Buffer) *FieldData {

	tmp := make([]byte, field_def.data_size)
	n, err := buf.Read(tmp)

	if n != field_def.data_size || err != nil {
		if n != field_def.data_size {
			iso_msg.buffer_underflow_error(field_def.name)
		} else {
			iso_msg.field_parse_error(field_def.name, err)
		}
	}

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

func (field_def *FixedFieldDef) IsFixed() bool {
	return true
}

func (field_def *FixedFieldDef) DataLength() int {
	return field_def.data_size
}

func (field_def *FixedFieldDef) EncodedLength(data_len int) []byte {
	//not applicable to fixed fields
	panic("illegal operation")
}

func (f_def *FixedFieldDef) SetId(id int) {
	f_def.id = id
}

func (f_def *FixedFieldDef) GetId() int {
	return f_def.id
}

func (f_def *FixedFieldDef) SetBitPosition(id int) {
	f_def.b_pos=id;
}

func (f_def *FixedFieldDef) BitPosition() int {
	return f_def.b_pos
}

func (field_def *FixedFieldDef) String() string {
	return field_def.name
}
