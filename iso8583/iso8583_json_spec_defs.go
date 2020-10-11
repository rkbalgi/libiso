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
	"bytes"
	"container/list"
	"encoding/json"
	"log"
	"os"
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

func displaySpecs() {

	for k, v := range specMap {
		spec := v
		buf := bytes.NewBufferString("")
		for l := spec.fieldsDefList.Front(); l != nil; l = l.Next() {

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

					for _, fDef := range bmp.subFieldDef {
						if fDef != nil {
							buf.WriteString(fDef.Def() + "\n")

						}
					}

				}

			} //end swicth
		}

		log.Printf("Spec: %s : Defs: \n%s\n", k, buf.String())

	} //end for-range

}

var specInit bool = false

func ReadSpecDefs(fileName string) {

	if !specInit {
		specMap = make(map[string]*MessageDef)
	} else {
		for k := range specMap {
			delete(specMap, k)
		}
	}

	//initialize map of all spec's to their definitions
	//read all specs from the json file
	//lets just display details of all defined specs

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("failed to open specs.json file", err.Error())
		return
	}

	tmpData := make([]byte, 100)
	buf := bytes.NewBuffer([]byte{})
	for {
		count, err := file.Read(tmpData)
		if err != nil && err.Error() == "EOF" {
			break
		} else if err != nil {
			log.Println("failed to read from specs.json file", err.Error())
			break
		}
		buf.Write(tmpData[:count])
	}

	ReadSpecDefsFromBuf(buf)

	specInit = true
	displaySpecs()

}

func ReadSpecDefsFromBuf(buf *bytes.Buffer) {

	specDefs := new(JsonSpecDefs)
	err := json.Unmarshal(buf.Bytes(), specDefs)
	if err != nil {
		log.Println("failed to parse specs.json", err.Error())
		return
	}

	for _, spec := range specDefs.Specs {

		iso8583MsgDef := new(MessageDef)
		iso8583MsgDef.specName = spec.SpecName
		iso8583MsgDef.fieldSeq = 0
		iso8583MsgDef.fieldsDefList = list.New()

		for _, isoFieldDef := range spec.Fields {

			if isoFieldDef.Type == "Fixed" {
				fixedFieldDef := constructFixedFieldDef(&isoFieldDef)
				fixedFieldDef.SetId(iso8583MsgDef.nextFieldSeq())
				if fixedFieldDef != nil {
					iso8583MsgDef.addField(fixedFieldDef)
				}

			} else if isoFieldDef.Type == "Variable" {
				varFieldDef := constructVariableFieldDef(&isoFieldDef)
				varFieldDef.SetId(iso8583MsgDef.nextFieldSeq())
				if varFieldDef != nil {
					iso8583MsgDef.addField(varFieldDef)
				}

			} else if isoFieldDef.Type == "Bitmapped" {

				bmpFieldDef := constructBmpFieldDef(iso8583MsgDef, &isoFieldDef)
				bmpFieldDef.SetId(iso8583MsgDef.nextFieldSeq())
				if bmpFieldDef != nil {
					iso8583MsgDef.addField(bmpFieldDef)
				}

			} else {
				log.Println("unsupported field type - ", isoFieldDef.Type)
			}
		}
		specMap[iso8583MsgDef.specName] = iso8583MsgDef
	}

}

func constructFixedFieldDef(jsonFieldDef *JsonFieldDef) *FixedFieldDef {

	encodingType := getEncoding(jsonFieldDef, jsonFieldDef.Attrs.DataEncoding)
	fixedFieldDef := NewFixedFieldDef(jsonFieldDef.Name, encodingType, jsonFieldDef.Attrs.FieldLength)
	fixedFieldDef.SetBitPosition(jsonFieldDef.BitPosition)

	log.Println("processed field def- " + jsonFieldDef.Name)

	return fixedFieldDef
}

func constructVariableFieldDef(jsonFieldDef *JsonFieldDef) *VariableFieldDef {

	lenEncodingType := getEncoding(jsonFieldDef, jsonFieldDef.Attrs.FieldIndicatorEncoding)
	dataEncodingType := getEncoding(jsonFieldDef, jsonFieldDef.Attrs.DataEncoding)

	varFieldDef := NewVariableFieldDef(jsonFieldDef.Name, lenEncodingType, dataEncodingType, jsonFieldDef.Attrs.FieldIndicatorLength)
	varFieldDef.SetBitPosition(jsonFieldDef.BitPosition)

	log.Println("processed variable field def- " + jsonFieldDef.Name)

	return varFieldDef
}

func constructBmpFieldDef(iso8583MsgDef *MessageDef, jsonFieldDef *JsonFieldDef) *BitMap {

	bmp := NewBitMap()
	//var bmp_field BitmappedField = bmp

	for _, childField := range jsonFieldDef.Children {
		switch childField.Type {
		case "Fixed":
			{
				fld := constructFixedFieldDef(&childField)
				fld.SetId(iso8583MsgDef.nextFieldSeq())
				bmp.subFieldDef[childField.BitPosition] = fld

			}
		case "Variable":
			{
				fld := constructVariableFieldDef(&childField)
				fld.SetId(iso8583MsgDef.nextFieldSeq())
				bmp.subFieldDef[childField.BitPosition] = fld
			}
		default:
			{
				log.Println("unsupported child field for bitmapped parent field - ", childField.Type)
			}
		}

	}
	log.Println("processed bitmapped field - " + jsonFieldDef.Name)
	return bmp
}

func getEncoding(isoFieldDef *JsonFieldDef, data string) int {

	encodingType := 0
	switch data {
	case "ascii":
		{
			encodingType = asciiEncoding
		}
	case "ebcdic":
		{
			encodingType = ebcdicEncoding
		}
	case "binary":
		{
			encodingType = binaryEncoding
		}
	case "bcd":
		{
			encodingType = bcdEncoding
		}
	default:
		{
			log.Panicf("invalid encoding [%s] on %s", data, isoFieldDef.Name)

		}
	}
	return encodingType

}
