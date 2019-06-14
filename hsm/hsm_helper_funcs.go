package hsm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

const (
	PercentSign = '%'
)

var hexRegexp, _ = regexp.Compile("[0-9a-fA-F]")
var keySchemeRegexp, _ = regexp.Compile("[UTZ]")

func toASCII(data []byte) []byte {

	return []byte(hex.EncodeToString(data))

}

func parsePrologue(msgBuf *bytes.Buffer, pro *prologue, headerLen int) bool {

	ok := readFixedField(msgBuf, &pro.header, uint(headerLen), String)
	if !ok {
		return false
	} else {
		ok := readFixedField(msgBuf, &pro.commandName, 2, String)
		if !ok {
			return false
		}
	}

	return true

}

func parseEpilogue(msgBuf *bytes.Buffer, epi *epilogue) bool {

	if msgBuf.Len() > 0 {
		//we have some data
		epi.delimiter, _ = msgBuf.ReadByte()
		if epi.delimiter == PercentSign {
			parsedOk := readFixedField(msgBuf, &epi.lmkIdentifier, 2, DecimalInt)
			if !parsedOk {
				return false
			}
			if msgBuf.Len() > 0 {
				//end delimiter present
				epi.endMessageDelimiter, _ = msgBuf.ReadByte()
				if epi.endMessageDelimiter == 0x19 {
					tmp := make([]byte, msgBuf.Len())
					_, _ = msgBuf.Read(tmp)
					epi.messageTrailer = tmp
				} else {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true

}

func readKey(msgBuf *bytes.Buffer, reqStruct interface{}) bool {

	firstByte, _ := msgBuf.ReadByte()
	if keySchemeRegexp.MatchString(string(firstByte)) {
		var tmp []byte
		if firstByte == 'Z' {
			tmp = make([]byte, 16+1)
		} else if firstByte == 'U' {
			tmp = make([]byte, 32+1)
		} else if firstByte == 'T' {
			tmp = make([]byte, 48+1)
		} else {
			return false
		}
		msgBuf.UnreadByte()
		msgBuf.Read(tmp)
		reflect.ValueOf(reqStruct).Elem().SetString(string(tmp))

	} else if hexRegexp.MatchString(string(firstByte)) {
		msgBuf.UnreadByte()
		tmp := make([]byte, 16)
		msgBuf.Read(tmp)
		reflect.ValueOf(reqStruct).Elem().SetString(string(tmp))

	} else {
		return false
	}

	return true

}

func Dump(v interface{}) string {

	strBuilder := bytes.NewBufferString("\n")

	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)
	for i := 0; i < valueOf.NumField(); i++ {
		switch valueOf.Field(i).Kind() {

		default:
			{
				fmt.Println(valueOf.Field(i).String())
			}

		case reflect.Struct:
			{
				strBuilder.WriteString(Dump(valueOf.Field(i)))
				break
			}
		case reflect.String:
			{

				strBuilder.WriteString(fmt.Sprintf("[%-20s] : [%s]\n", typeOf.Field(i).Name, valueOf.Field(i).String()))
				break
			}
		case reflect.Slice:
			{
				strBuilder.WriteString(fmt.Sprintf("[%-20s] : [%s]\n", typeOf.Field(i).Name, hex.EncodeToString(valueOf.Field(i).Bytes())))
				break
			}
		case reflect.Uint:
			{

				strBuilder.WriteString(fmt.Sprintf("[%-20s] : [%d]\n", typeOf.Field(i).Name, valueOf.Field(i).Uint()))
				break
			}
		}
	}

	return string(strBuilder.Bytes())
}

func setFixedField(sf interface{}, fieldSize uint, fieldValue interface{}, dataType int) {

	switch dataType {
	case DecimalInt:
		{
			fmtSpec := fmt.Sprintf("%%0%dd", fieldSize)
			fieldData := fmt.Sprintf(fmtSpec, reflect.ValueOf(fieldValue).Uint())
			//fmt.Println("format spec",fmt_spec,field_data);
			fieldVal := []byte(fieldData)
			reflect.ValueOf(sf).Elem().SetBytes(fieldVal)
			break
		}
	default:
		{
			panic(fmt.Sprintf("set_fixed_field not implemented for this type - %d", dataType))
		}
	}

}

func readFixedField(msgBuf *bytes.Buffer, sf interface{}, size uint, dataType int) bool {

	var tmpDataBuf []byte = make([]byte, size)
	_, err := msgBuf.Read(tmpDataBuf)
	if checkParseError(err) {
		return false
	}

	switch dataType {
	case String:
		{

			reflect.ValueOf(sf).Elem().SetString(string(tmpDataBuf))
			break

		}

	case Binary:
		{
			reflect.ValueOf(sf).Elem().SetBytes(tmpDataBuf)
			break
		}

	case DecimalInt:
		{

			decimalVal, err := strconv.ParseUint(string(tmpDataBuf), 10, 32)
			if checkFormatError(err) {
				return false
			}

			reflect.ValueOf(sf).Elem().SetUint(uint64(decimalVal))
			break
		}

	case HexadecimalInt:
		{

			decimalVal, err := strconv.ParseUint(string(tmpDataBuf), 16, 32)
			if checkFormatError(err) {
				return false
			}
			reflect.ValueOf(sf).Elem().SetUint(decimalVal)
			break
		}
	}

	return true
}

func checkParseError(err error) bool {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "parsing error - %s", err.Error())
		return true
	}

	return false
}

func checkFormatError(err error) bool {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "format error - %s", err.Error())
		return true
	}

	return false
}
