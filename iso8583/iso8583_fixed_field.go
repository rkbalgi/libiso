package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type FixedFieldDef struct {
	name         string
	dataEncoding int
	dataSize     int //in bytes
	id           int
	bPos         int //bit position in bitmap
}

//create a new fixed field definition
func NewFixedFieldDef(pName string, pDataEncoding int, pFieldLen int) *FixedFieldDef {
	field := new(FixedFieldDef)
	field.name = pName
	field.dataEncoding = pDataEncoding
	field.dataSize = pFieldLen

	return field

}

func (fieldDef *FixedFieldDef) Def() string {
	return fmt.Sprintf("Name: %-40s ; Id: %04d ; Type: %-10s ; Length: %04d ;Encoding: [%-10s]",
		fieldDef.name, fieldDef.GetId(), "Fixed", fieldDef.dataSize, getEncodingType(fieldDef.dataEncoding))
}

//parse and return field data by reading appropriate bytes
//from the buffer buf
func (fieldDef *FixedFieldDef) Parse(
	isoMsg *Iso8583Message,
	fData *FieldData,
	buf *bytes.Buffer) *FieldData {

	tmp := make([]byte, fieldDef.dataSize)
	n, err := buf.Read(tmp)

	if n != fieldDef.dataSize || err != nil {
		if n != fieldDef.dataSize {
			isoMsg.bufferUnderflowError(fieldDef.name)
		} else {
			isoMsg.fieldParseError(fieldDef.name, err)
		}
	}

	fData.fieldData = tmp
	isoMsg.log.Printf("parsed: [%s]=[%s]", fieldDef.name, fData.String())

	return fData

}

func (fieldDef *FixedFieldDef) toString(data []byte) string {
	return fmt.Sprintf("[%s] = [%s]", fieldDef.name, hex.EncodeToString(data))
}

func (fieldDef *FixedFieldDef) getDataEncoding() int {
	return fieldDef.dataEncoding
}

//add the field data into buf as per the encoding
func (fieldDef *FixedFieldDef) Assemble(isoMsg *Iso8583Message, buf *bytes.Buffer) {

	//buf.Write(iso_msg);
}

func (fieldDef *FixedFieldDef) IsFixed() bool {
	return true
}

func (fieldDef *FixedFieldDef) DataLength() int {
	return fieldDef.dataSize
}

func (fieldDef *FixedFieldDef) EncodedLength(dataLen int) []byte {
	//not applicable to fixed fields
	panic("illegal operation")
}

func (fieldDef *FixedFieldDef) SetId(id int) {
	fieldDef.id = id
}

func (fieldDef *FixedFieldDef) GetId() int {
	return fieldDef.id
}

func (fieldDef *FixedFieldDef) SetBitPosition(id int) {
	fieldDef.bPos = id
}

func (fieldDef *FixedFieldDef) BitPosition() int {
	return fieldDef.bPos
}

func (fieldDef *FixedFieldDef) String() string {
	return fieldDef.name
}
