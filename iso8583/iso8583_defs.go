package iso8583

import (
	"bytes"
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

func init() {

	//initialize map of all spec's to their definitions
	spec_map = make(map[string]*Iso8583MessageDef)

	iso8583_msg_def = new(Iso8583MessageDef)
	iso8583_msg_def.spec_name = "ISO8583_1 v1 (ASCII)"
	iso8583_msg_def.fields = make([]IsoField, 128+1)
	//lets use a 1 based slice accessor

	//add all defined fields
	iso8583_msg_def.fields[2] = NewVariableFieldDef("PAN", ascii_encoding, ascii_encoding, 2)
	iso8583_msg_def.fields[3] = NewFixedFieldDef("Processing Code", ebcdic_encoding, 6)
	iso8583_msg_def.fields[4] = NewFixedFieldDef("Transaction Amount", ascii_encoding, 12)
	iso8583_msg_def.fields[14] = NewFixedFieldDef("Expiry Date", ascii_encoding, 4)

	iso8583_msg_def.fields[33] = NewVariableFieldDef("Test Var Binary", binary_encoding, binary_encoding, 2)
	iso8583_msg_def.fields[34] = NewVariableFieldDef("Test Var BCD", bcd_encoding, binary_encoding, 2)

	iso8583_msg_def.fields[35] = NewVariableFieldDef("Track II", ebcdic_encoding, ebcdic_encoding, 3)
	iso8583_msg_def.fields[38] = NewFixedFieldDef("Approval Code", ascii_encoding, 6)
	iso8583_msg_def.fields[39] = NewFixedFieldDef("Action Code", ascii_encoding, 3)

	iso8583_msg_def.fields[55] = NewVariableFieldDef("ICC Data", ascii_encoding, binary_encoding, 3)
	iso8583_msg_def.fields[64] = NewFixedFieldDef("MAC1", binary_encoding, 8)
	iso8583_msg_def.fields[128] = NewFixedFieldDef("MAC2", binary_encoding, 8)

	fmt.Println("initialized -" + iso8583_msg_def.spec_name)
	spec_map[iso8583_msg_def.spec_name] = iso8583_msg_def

}

type IsoField interface {
	Parse(*Iso8583Message, *bytes.Buffer) *FieldData
	Assemble(*Iso8583Message, *bytes.Buffer)
	String() string
	IsFixed() bool
	DataLength() int
	EncodedLength(int) []byte
	to_string([]byte) string
	get_data_encoding() int
}

type Iso8583MessageDef struct {
	spec_name string
	fields    []IsoField
}

func str_to_uint64(str_val string) uint64 {

	val, err := strconv.ParseUint(str_val, 10, 64)
	if err != nil {
		panic(err.Error())
	}

	return val

}
