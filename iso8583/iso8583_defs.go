package iso8583

import (
	"bytes"
	"container/list"
	_ "container/list"
	_ "encoding/binary"
	_ "encoding/hex"
	_ "errors"
	"log"
	_ "log"
	_ "os"
)

const (
	ebcdicEncoding = iota
	asciiEncoding  = iota + 1
	bcdEncoding    = iota + 2
	binaryEncoding = iota + 3
)

const (
	V1         = "V1"
	V0         = "V0"
	IsoMsg1100 = "1100"
	IsoMsg1110 = "1110"
	IsoMsg1420 = "1420"
	IsoMsg1430 = "1430"
	IsoMsg1804 = "1804"
	IsoMsg1814 = "1814"

	IsoRespDecline  = "100"
	IsoRespPickup   = "200"
	IsoRespApproval = "000"
	IsoFormatError  = "909"
	IsoRespDrop     = "999"
)

var specMap map[string]*MessageDef

//To send data back and forth between browser and paysim
//application
type WebMsgData struct {
	Spec      string
	Type      string
	DataArray []string
}

func getEncodingType(encoding int) string {
	switch encoding {
	case asciiEncoding:
		{
			return "ascii"
		}
	case bcdEncoding:
		{
			return "bcd"
		}
	case binaryEncoding:
		{
			return "binary"
		}
	case ebcdicEncoding:
		{
			return "ebcdic"
		}
	default:
		{
			return "unknown"
		}
	}
}

type IsoField interface {
	Parse(*Iso8583Message, *FieldData, *bytes.Buffer) *FieldData
	Assemble(*Iso8583Message, *bytes.Buffer)
	String() string
	IsFixed() bool
	SetId(int)
	GetId() int
	DataLength() int
	EncodedLength(int) []byte
	toString([]byte) string
	getDataEncoding() int
	Def() string

	BitPosition() int
	SetBitPosition(int)
}

type MessageDef struct {
	specName      string
	fieldsDefList *list.List
	fieldSeq      int
}

func (isoDef *MessageDef) nextFieldSeq() int {
	seq := isoDef.fieldSeq
	isoDef.fieldSeq = isoDef.fieldSeq + 1
	return seq
}

func (isoDef *MessageDef) addField(field interface{}) {

	switch field.(type) {
	case IsoField:
		{
			isoField := field.(IsoField)
			isoField.SetId(isoDef.nextFieldSeq())
		}
	case BitmappedField:
		{
			bmpField := field.(BitmappedField)
			bmpField.SetSpec(isoDef)
			bmpField.SetId(isoDef.nextFieldSeq())
		}
	default:
		{
			log.Println("yikes")
		}
	}

	isoDef.fieldsDefList.PushBack(field)
}

func (isoDef *MessageDef) Name() string {
	return isoDef.specName
}
