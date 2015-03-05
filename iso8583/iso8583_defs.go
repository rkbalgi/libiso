package iso8583

import (
	"bytes"
	"container/list"
	_ "container/list"
	_ "encoding/binary"
	_ "encoding/hex"
	_ "errors"
	"fmt"
	_ "log"
	_ "os"
	"strconv"
)

const (
	ebcdic_encoding = iota
	ascii_encoding  = iota + 1
	bcd_encoding    = iota + 2
	binary_encoding = iota + 3
)

const (
	V1           = "V1"
	V0           = "V0"
	ISO_MSG_1100 = "1100"
	ISO_MSG_1110 = "1110"
	ISO_MSG_1420 = "1420"
	ISO_MSG_1430 = "1430"
	ISO_MSG_1804 = "1804"
	ISO_MSG_1814 = "1814"

	ISO_RESP_DECLINE  = "100"
	ISO_RESP_APPROVAL = "000"
	ISO_FORMAT_ERROR  = "909"
)

var spec_map map[string]*Iso8583MessageDef

//To send data back and forth between browser and paysim
//application
type WebMsgData struct {
	Spec      string
	Type      string
	DataArray []string
}

func get_encoding_type(encoding int) string {
	switch encoding {
	case ascii_encoding:
		{
			return "ascii"
		}
	case bcd_encoding:
		{
			return "bcd"
		}
	case binary_encoding:
		{
			return "binary"
		}
	case ebcdic_encoding:
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
	to_string([]byte) string
	get_data_encoding() int
	Def() string
}

type Iso8583MessageDef struct {
	spec_name       string
	fields_def_list *list.List
	field_seq       int
}

func (iso_def *Iso8583MessageDef) next_field_seq() int {
	seq := iso_def.field_seq
	iso_def.field_seq = iso_def.field_seq + 1
	return seq
}

func (iso_def *Iso8583MessageDef) add_field(field interface{}) {

	switch field.(type) {
	case IsoField:
		{
			iso_field := field.(IsoField)
			iso_field.SetId(iso_def.next_field_seq())
		}
	case BitmappedField:
		{
			bmp_field := field.(BitmappedField)
			bmp_field.SetSpec(iso_def)
			bmp_field.SetId(iso_def.next_field_seq())
		}
	default:
		{
			fmt.Println("yikes")
		}
	}

	iso_def.fields_def_list.PushBack(field)
}

func str_to_uint64(str_val string) uint64 {

	val, err := strconv.ParseUint(str_val, 10, 64)
	if err != nil {
		panic(err.Error())
	}

	return val

}
