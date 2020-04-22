package iso8583

import (
	"bytes"
	"encoding/binary"
	_ "encoding/binary"
	"encoding/hex"
	"fmt"
	"libiso/encoding/ebcdic"
	"strconv"
)

type VariableFieldDef struct {
	name           string
	dataEncoding   int
	lengthEncoding int
	lenIndSize     int //in bytes
	id             int
	bPos           int
}

//create a new fixed field definition
func NewVariableFieldDef(pName string,
	pLenEncoding int,
	pDataEncoding int,
	pLenIndSize int) *VariableFieldDef {

	field := new(VariableFieldDef)
	field.name = pName
	field.dataEncoding = pDataEncoding
	field.lengthEncoding = pLenEncoding
	field.lenIndSize = pLenIndSize

	return field

}

func (field *VariableFieldDef) Def() string {
	return fmt.Sprintf("Name: %-40s ; Id: %04d ; Type: %-10s ; Length Indicator Length: %04d ;Encoding (Length Indicator): [%-10s] ; Encoding (Data): [%-10s]",
		field.name, field.GetId(), "Variable", field.lenIndSize, getEncodingType(field.lengthEncoding), getEncodingType(field.dataEncoding))
}

func (field *VariableFieldDef) toString(data []byte) string {
	return fmt.Sprintf("[%s] = [%s]", field.name, hex.EncodeToString(data))
}

func (field *VariableFieldDef) getDataEncoding() int {
	return field.dataEncoding
}

func (field *VariableFieldDef) Parse(
	isoMsg *Iso8583Message,
	fData *FieldData,
	buf *bytes.Buffer) *FieldData {

	tmp := make([]byte, field.lenIndSize)
	n, err := buf.Read(tmp)
	if n != field.lenIndSize || err != nil {

		if n != field.lenIndSize {
			isoMsg.bufferUnderflowError(field.name)
		} else {
			isoMsg.fieldParseError(field.name, err)
		}

		isoMsg.fieldParseError(field.name, err)
	}

	var dataLen uint64 = 0
	switch field.lengthEncoding {
	case asciiEncoding:
		{
			dataLen, _ = strconv.ParseUint(string(tmp), 10, 64)
		}
	case ebcdicEncoding:
		{

			dataLen, _ = strconv.ParseUint(ebcdic.EncodeToString(tmp), 10, 64)
		}
	case binaryEncoding:
		{
			dataLen, _ = strconv.ParseUint(hex.EncodeToString(tmp), 16, 64)
		}
	case bcdEncoding:
		{
			dataLen, _ = strconv.ParseUint(hex.EncodeToString(tmp), 10, 64)
		}
	default:
		{
			panic("unsupported encoding")
		}
	}

	bFieldData := make([]byte, dataLen)
	n, err = buf.Read(bFieldData)
	if uint64(n) != dataLen || err != nil {

		if uint64(n) != dataLen {
			isoMsg.bufferUnderflowError(field.name)
		} else {
			isoMsg.fieldParseError(field.name, err)
		}

	}
	fData.fieldData = bFieldData
	isoMsg.log.Printf("parsed: [%s]=[%s]", field.name, fData.String())

	return fData

}

//add the field data into buf as per the encoding
func (field *VariableFieldDef) Assemble(isoMsg *Iso8583Message, buf *bytes.Buffer) {

}

func (field *VariableFieldDef) IsFixed() bool {
	return false
}

func (field *VariableFieldDef) DataLength() int {
	//not applicable
	return -1
}

//return the length part of the variable field
//as a []byte slice
func (field *VariableFieldDef) EncodedLength(dataLen int) []byte {
	//not applicable

	if field.lenIndSize > 4 &&
		(field.lengthEncoding == bcdEncoding || field.lengthEncoding == binaryEncoding) {
		panic("[llvar] invalid length indicator size for bcd/binary - >4")
	}

	var ll []byte

	switch field.lengthEncoding {
	case binaryEncoding:
		{
			switch field.lenIndSize {
			case 1:
				{
					ll = []byte{byte(dataLen)}
				}
			case 2:
				{
					ll = make([]byte, 2)
					binary.BigEndian.PutUint16(ll, uint16(dataLen))
				}
			case 4:
				{
					ll = make([]byte, 4)
					binary.BigEndian.PutUint32(ll, uint32(dataLen))
				}
			default:
				{
					panic(fmt.Sprintf("[llvar] invalid length indicator size for binary field - %d", field.lenIndSize))
				}
			}

		}

	case bcdEncoding:
		{

			switch field.lenIndSize {
			case 1:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%0d", dataLen))
				}
			case 2:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%04d", dataLen))
				}
			case 4:
				{
					ll, _ = hex.DecodeString(fmt.Sprintf("%08d", dataLen))
				}
			default:
				{
					panic(fmt.Sprintf("[llvar] invalid length indicator size for binary field - %d", field.lenIndSize))
				}
			}

		}

	case asciiEncoding:
		{

			lenStr := encodedLengthAsString(field.lenIndSize, dataLen)
			ll = []byte(lenStr)

		}

	case ebcdicEncoding:
		{

			lenStr := encodedLengthAsString(field.lenIndSize, dataLen)
			ll = ebcdic.Decode(lenStr)

		}

	}

	return ll
}

func encodedLengthAsString(lenIndSize int, dataLen int) string {

	var tmp string

	switch lenIndSize {
	case 1:
		{
			if dataLen > 9 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", dataLen))
			}
			tmp = fmt.Sprintf("%d", dataLen)
		}
	case 2:
		{
			if dataLen > 99 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", dataLen))
			}
			tmp = fmt.Sprintf("%02d", dataLen)
		}
	case 3:
		{
			if dataLen > 999 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", dataLen))
			}
			tmp = fmt.Sprintf("%03d", dataLen)
		}
	case 4:
		{
			if dataLen > 9999 {
				panic(fmt.Sprintf("[llvar] data length > %d\n", dataLen))
			}
			tmp = fmt.Sprintf("%04d", dataLen)
		}
	default:
		{
			panic(fmt.Sprintf("[llvar] invalid length indicator size for  field - %d", lenIndSize))
		}
	}

	return tmp
}

func (field *VariableFieldDef) SetId(id int) {
	field.id = id
}

func (field *VariableFieldDef) GetId() int {
	return field.id
}

func (field *VariableFieldDef) String() string {
	return field.name

}

func (field *VariableFieldDef) SetBitPosition(id int) {
	field.bPos = id
}

func (field *VariableFieldDef) BitPosition() int {
	return field.bPos
}
