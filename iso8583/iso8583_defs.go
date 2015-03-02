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

var iso8583_msg_def *Iso8583MessageDef
var spec_map map[string]*Iso8583MessageDef



//To send data back and forth between browser and paysim
//application
type WebMsgData struct {
	Spec      string
	Type      string
	DataArray []string
}

func init() {

	//initialize map of all spec's to their definitions
	spec_map = make(map[string]*Iso8583MessageDef)
    //read all specs from the json file 	
	read_spec_defs();
	fmt.Println(spec_map);


	/*iso8583_msg_def = new(Iso8583MessageDef)
	iso8583_msg_def.spec_name = "ISO8583_1 v1 (ASCII)"
	iso8583_msg_def.field_seq = 0
	iso8583_msg_def.fields_def_list = list.New()

	//iso8583_msg_def.add_field(NewFixedFieldDef("Header", ascii_encoding, 4))
	iso8583_msg_def.add_field(iso8583_msg_def, NewFixedFieldDef("Message Type", ascii_encoding, 4))

	bmp := NewBitMap()
	var bmp_field BitmappedField = bmp
	iso8583_msg_def.add_field(iso8583_msg_def, bmp_field)

	//add all subfields of bitmap
	bmp.add_variable_field(2, "PAN", ascii_encoding, ascii_encoding, 2)
	bmp.add_fixed_field(3, "Processing Code", ebcdic_encoding, 6)
	bmp.add_fixed_field(4, "Transaction Amount", ascii_encoding, 12)
	bmp.add_fixed_field(14, "Expiry Date", ascii_encoding, 4)

	bmp.add_variable_field(33, "Test Var Binary", binary_encoding, binary_encoding, 2)
	bmp.add_variable_field(34, "Test Var BCD", bcd_encoding, binary_encoding, 2)

	bmp.add_variable_field(35, "Track II", ebcdic_encoding, ebcdic_encoding, 3)
	bmp.add_fixed_field(38, "Approval Code", ascii_encoding, 6)
	bmp.add_fixed_field(39, "Action Code", ascii_encoding, 3)

	bmp.add_variable_field(55, "ICC Data", ascii_encoding, binary_encoding, 3)
	bmp.add_fixed_field(64, "MAC1", binary_encoding, 8)
	bmp.add_fixed_field(128, "MAC2", binary_encoding, 8)

	//iso8583_msg_def.add_field(NewFixedFieldDef("Trailer", ascii_encoding, 4))

	fmt.Println("initialized -" + iso8583_msg_def.spec_name)
	spec_map[iso8583_msg_def.spec_name] = iso8583_msg_def*/

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
