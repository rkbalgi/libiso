package iso8583

//fixed field attributes = key field{true|false};length;encoding;padding_type
//variable field attributes = key field{true|false};len indicator size in bytes;len ind encoding;data encoding_;nibbles={true|false}  (no padding type)

//encoding can be ascii, binary, bcd or ebcdic

//nibbles=true implies that the size is indicated in number of half bytes (or nibbles) -
// i.e. is field value is "00000000" then if nibbles is true length indicator is 8 else 4

//a key field uniquely identifies a message/transaction
//more than a single field can be a key - a key is formed by concatenation
//padding_type has not be implemented, leave as default

// _ is used when the attribute is not applicable
//only Bitmapped field will have children, and the BitPosition indicates their position
//in the bitmap

import (
	_ "bufio"
	"bytes"
	"container/list"
	"encoding/json"

	"github.com/rkbalgi/go/paysim/demo"
	pylog "github.com/rkbalgi/go/paysim/log"
	_ "io"
	"log"
	"os"
	_ "strconv"
	_ "strings"
)

type FieldAttributes struct {
	DataEncoding           string
	FieldLength            int
	Key                    bool
	FieldIndicatorLength   int
	FieldIndicatorEncoding string
	Padding                string
}

//JsonFieldDef represents the definition of a field in the specs.json file
type JsonFieldDef struct {
	BitPosition int
	Name        string
	Type        string
	Attrs       FieldAttributes
	Children    []JsonFieldDef
}

//JsonSpecDef represents the definition of a field in the specs.json file
type JsonSpecDef struct {
	SpecName string
	Fields   []JsonFieldDef
}

type JsonSpecDefs struct {
	Specs []JsonSpecDef
}

func display_specs() {

	for k, v := range spec_map {
		spec := v
		buf := bytes.NewBufferString("")
		for l := spec.fields_def_list.Front(); l != nil; l = l.Next() {

			switch obj := l.Value.(type) {
			case IsoField:
				{
					buf.WriteString(obj.Def() + "\n")
					break
				}
			case BitmappedField:
				{
					bmp := obj.(*BitMap)
					buf.WriteString(obj.Def() + "\n")

					for _, f_def := range bmp.sub_field_def {
						if f_def != nil {
							buf.WriteString(f_def.Def() + "\n")

						}
					}

				}

			} //end swicth
		}

		pylog.Printf("Spec: %s : Defs: \n%s\n", k, buf.String())

	} //end for-range

}

var spec_init bool = false

func ReadDemoSpecDefs() {

	if !spec_init {
		spec_map = make(map[string]*Iso8583MessageDef)
	} else {
		for k, _ := range spec_map {
			delete(spec_map, k)
		}
	}
	str := demo.Demo_Specs
	ReadSpecDefsFromBuf(bytes.NewBufferString(str))
	spec_init = true
	display_specs()
}

func ReadSpecDefs(file_name string) {

	if !spec_init {
		spec_map = make(map[string]*Iso8583MessageDef)
	} else {
		for k, _ := range spec_map {
			delete(spec_map, k)
		}
	}

	//initialize map of all spec's to their definitions
	//read all specs from the json file
	//lets just display details of all defined specs

	file, err := os.Open(file_name)
	if err != nil {
		pylog.Log("failed to open specs.json file", err.Error())
		return
	}

	tmp_data := make([]byte, 100)
	buf := bytes.NewBuffer([]byte{})
	for {
		count, err := file.Read(tmp_data)
		if err != nil && err.Error() == "EOF" {
			break
		} else if err != nil {
			pylog.Log("failed to read from specs.json file", err.Error())
			break
		}
		buf.Write(tmp_data[:count])
	}

	ReadSpecDefsFromBuf(buf)

	spec_init = true
	display_specs()

}

func ReadSpecDefsFromBuf(buf *bytes.Buffer) {

	spec_defs := new(JsonSpecDefs)
	err := json.Unmarshal(buf.Bytes(), spec_defs)
	if err != nil {
		pylog.Log("failed to parse specs.json", err.Error())
		return
	}

	for _, spec := range spec_defs.Specs {

		iso8583_msg_def := new(Iso8583MessageDef)
		iso8583_msg_def.spec_name = spec.SpecName
		iso8583_msg_def.field_seq = 0
		iso8583_msg_def.fields_def_list = list.New()

		for _, iso_field_def := range spec.Fields {

			if iso_field_def.Type == "Fixed" {
				fixed_field_def := construct_fixed_field_def(&iso_field_def)
				fixed_field_def.SetId(iso8583_msg_def.next_field_seq())
				if fixed_field_def != nil {
					iso8583_msg_def.add_field(fixed_field_def)
				}

			} else if iso_field_def.Type == "Variable" {
				var_field_def := construct_variable_field_def(&iso_field_def)
				var_field_def.SetId(iso8583_msg_def.next_field_seq())
				if var_field_def != nil {
					iso8583_msg_def.add_field(var_field_def)
				}

			} else if iso_field_def.Type == "Bitmapped" {

				bmp_field_def := construct_bmp_field_def(iso8583_msg_def, &iso_field_def)
				bmp_field_def.SetId(iso8583_msg_def.next_field_seq())
				if bmp_field_def != nil {
					iso8583_msg_def.add_field(bmp_field_def)
				}

			} else {
				log.Panic("unsupported field type - ", iso_field_def.Type)
			}
		}
		spec_map[iso8583_msg_def.spec_name] = iso8583_msg_def
	}

}

func construct_fixed_field_def(json_field_def *JsonFieldDef) *FixedFieldDef {

	/*attrs := strings.Split(json_field_def.Attrs, ";")
	if len(attrs) != 4 {
		log.Panic("invalid attribute spec for fixed field -" + json_field_def.Name)
		return nil
	}
	//TODO:: for now only look at 1 and 2 attrs i.e. lengt and encoding
	field_len, err := strconv.ParseUint(attrs[1], 10, 32)
	if err != nil {
		log.Panic("invalid field length -", json_field_def.Name)
		return nil
	}*/
	encoding_type := get_encoding(json_field_def, json_field_def.Attrs.DataEncoding)
	fixed_field_def := NewFixedFieldDef(json_field_def.Name, encoding_type, json_field_def.Attrs.FieldLength)
	fixed_field_def.SetBitPosition(json_field_def.BitPosition)

	pylog.Log("processed field def- " + json_field_def.Name)

	return fixed_field_def
}

func construct_variable_field_def(json_field_def *JsonFieldDef) *VariableFieldDef {

	/*attrs := strings.Split(json_field_def.Attrs, ";")
	if len(attrs) != 5 {
		log.Panic("invalid attribute spec for variable field -" + json_field_def.Name)
		return nil
	}
	//TODO:: for now only look at 1,2 and 3 attrs i.e. lengt and encoding
	field_len, err := strconv.ParseUint(attrs[1], 10, 32)
	if err != nil {
		log.Panic("invalid length indicator length -", json_field_def.Name)
		return nil
	}*/
	len_encoding_type := get_encoding(json_field_def, json_field_def.Attrs.FieldIndicatorEncoding)
	data_encoding_type := get_encoding(json_field_def, json_field_def.Attrs.DataEncoding)

	var_field_def := NewVariableFieldDef(json_field_def.Name, len_encoding_type, data_encoding_type, json_field_def.Attrs.FieldIndicatorLength)
	var_field_def.SetBitPosition(json_field_def.BitPosition)

	pylog.Log("processed variable field def- " + json_field_def.Name)

	return var_field_def
}

func construct_bmp_field_def(iso8583_msg_def *Iso8583MessageDef, json_field_def *JsonFieldDef) *BitMap {

	bmp := NewBitMap()
	//var bmp_field BitmappedField = bmp

	for _, child_field := range json_field_def.Children {
		switch child_field.Type {
		case "Fixed":
			{
				fld := construct_fixed_field_def(&child_field)
				fld.SetId(iso8583_msg_def.next_field_seq())
				bmp.sub_field_def[child_field.BitPosition] = fld

			}
		case "Variable":
			{
				fld := construct_variable_field_def(&child_field)
				fld.SetId(iso8583_msg_def.next_field_seq())
				bmp.sub_field_def[child_field.BitPosition] = fld
			}
		default:
			{
				log.Panic("unsupported child field for bitmapped parent field - ", child_field.Type)
			}
		}

	}
	pylog.Log("processed bitmapped field - " + json_field_def.Name)
	return bmp
}

func get_encoding(iso_field_def *JsonFieldDef, data string) int {

	encoding_type := 0
	switch data {
	case "ascii":
		{
			encoding_type = ascii_encoding
		}
	case "ebcdic":
		{
			encoding_type = ebcdic_encoding
		}
	case "binary":
		{
			encoding_type = binary_encoding
		}
	case "bcd":
		{
			encoding_type = bcd_encoding
		}
	default:
		{
			log.Panicf("invalid encoding [%s] on %s", data, iso_field_def.Name)

		}
	}
	return encoding_type

}
