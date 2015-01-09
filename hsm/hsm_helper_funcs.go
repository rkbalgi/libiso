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
	PERCENT_SIGN = '%'
)

var hex_regexp, _ = regexp.Compile("[0-9a-fA-F]")
var key_scheme_regexp, _ = regexp.Compile("[UTZ]")

func to_ascii(data []byte) []byte {

	return ([]byte(hex.EncodeToString(data)))

}

func parse_prologue(msg_buf *bytes.Buffer, pro *prologue, header_len int) bool {

	parse_ok := read_fixed_field(msg_buf, &pro.header, uint(header_len), String)
	if !parse_ok {
		return false
	} else {
		parse_ok := read_fixed_field(msg_buf, &pro.command_name, 2, String)
		if !parse_ok {
			return false
		}
	}

	return true

}

func parse_epilogue(msg_buf *bytes.Buffer, epi *epilogue) bool {

	if msg_buf.Len() > 0 {
		//we have some data
		epi.delimiter, _ = msg_buf.ReadByte()
		if epi.delimiter == PERCENT_SIGN {
			parsed_ok := read_fixed_field(msg_buf, &epi.lmk_identifier, 2, DecimalInt)
			if !parsed_ok {
				return false
			}
			if msg_buf.Len() > 0 {
				//end delimiter present
				epi.end_message_delimiter, _ = msg_buf.ReadByte()
				if epi.end_message_delimiter == 0x19 {
					tmp := make([]byte, msg_buf.Len())
					msg_buf.Read(tmp)
					epi.message_trailer = tmp
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

func read_key(msg_buf *bytes.Buffer, req_struct interface{}) bool {

	first_byte, _ := msg_buf.ReadByte()
	if key_scheme_regexp.MatchString(string(first_byte)) {
		var tmp []byte
		if first_byte == 'Z' {
			tmp = make([]byte, 16+1)
		} else if first_byte == 'U' {
			tmp = make([]byte, 32+1)
		} else if first_byte == 'T' {
			tmp = make([]byte, 48+1)
		} else {
			return false
		}
		msg_buf.UnreadByte()
		msg_buf.Read(tmp)
		reflect.ValueOf(req_struct).Elem().SetString(string(tmp))

	} else if hex_regexp.MatchString(string(first_byte)) {
		msg_buf.UnreadByte()
		tmp := make([]byte, 16)
		msg_buf.Read(tmp)
		reflect.ValueOf(req_struct).Elem().SetString(string(tmp))

	} else {
		return false
	}

	return true

}

func Dump(struct_var interface{}) string {

	str_builder := bytes.NewBufferString("\n")

	value_of := reflect.ValueOf(struct_var)
	type_of := reflect.TypeOf(struct_var)
	for i := 0; i < value_of.NumField(); i++ {
		switch value_of.Field(i).Kind() {
			
			default:{
				fmt.Println(value_of.Field(i).String());
			}

		case reflect.Struct:
			{
				str_builder.WriteString(Dump(value_of.Field(i)))
				break;
			}
		case reflect.String:
			{

				str_builder.WriteString(fmt.Sprintf("[%-20s] : [%s]\n", type_of.Field(i).Name, value_of.Field(i).String()))
				break
			}
		case reflect.Slice:
			{
				str_builder.WriteString(fmt.Sprintf("[%-20s] : [%s]\n", type_of.Field(i).Name, hex.EncodeToString(value_of.Field(i).Bytes())))
				break
			}
		case reflect.Uint:
			{

				str_builder.WriteString(fmt.Sprintf("[%-20s] : [%d]\n", type_of.Field(i).Name, value_of.Field(i).Uint()))
				break
			}
		}
	}

	return (string(str_builder.Bytes()))
}

func set_fixed_field(struct_field interface{}, field_size uint, field_value interface{}, data_type int) {

	switch data_type {
	case DecimalInt:
		{
			fmt_spec := fmt.Sprintf("%%0%dd", field_size)
			field_data := fmt.Sprintf(fmt_spec, reflect.ValueOf(field_value).Uint())
			//fmt.Println("format spec",fmt_spec,field_data);
			field_val := []byte(field_data)
			reflect.ValueOf(struct_field).Elem().SetBytes(field_val)
			break
		}
	default:
		{
			panic(fmt.Sprintf("set_fixed_field not implemented for this type - %d", data_type))
		}
	}

}

func read_fixed_field(msg_buf *bytes.Buffer, struct_field interface{}, size uint, data_type int) bool {

	var tmp_data_buf []byte = make([]byte, size)
	_, err := msg_buf.Read(tmp_data_buf)
	if check_parse_error(err) {
		return false
	}

	switch data_type {
	case String:
		{

			reflect.ValueOf(struct_field).Elem().SetString(string(tmp_data_buf))
			break

		}

	case Binary:
		{
			reflect.ValueOf(struct_field).Elem().SetBytes(tmp_data_buf)
			break
		}

	case DecimalInt:
		{

			decimal_val, err := strconv.ParseUint(string(tmp_data_buf), 10, 32)
			if check_format_error(err) {
				return false
			}

			reflect.ValueOf(struct_field).Elem().SetUint(uint64(decimal_val))
			break
		}

	case HexadecimalInt:
		{

			decimal_val, err := strconv.ParseUint(string(tmp_data_buf), 16, 32)
			if check_format_error(err) {
				return false
			}
			reflect.ValueOf(struct_field).Elem().SetUint(decimal_val)
			break
		}
	}

	return true
}

func check_parse_error(err error) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing error - %s", err.Error())
		return (true)
	}

	return (false)
}

func check_format_error(err error) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "format error - %s", err.Error())
		panic("")
		return (true)
	}

	return (false)
}
