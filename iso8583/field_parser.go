package iso8583

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// ErrInsufficientData is an error when there is not enough data in the raw message to parse it
var ErrInsufficientData = errors.New("libiso: Insufficient data to parse field")

// ErrLargeLengthIndicator is an error that could happen when a large length indicator is used in a variable field
var ErrLargeLengthIndicator = errors.New("libiso: Too large length indicator size. ")

// ErrInvalidEncoding is when an unsupported encoding is used for a field
var ErrInvalidEncoding = errors.New("libiso: Invalid encoding")

// ParsedMsg is a type that represents a parsed form of a ISO8583 message
type ParsedMsg struct {
	IsRequest bool
	Msg       *Message

	//A map of Id to FieldData
	FieldDataMap map[int]*FieldData

	// MessageKey is a value that unique identifies a transaction
	// (usually a combination of fields like STAN, PAN etc)
	MessageKey string
}

// Get returns the field-data from the parsed message given its name
func (pMsg *ParsedMsg) Get(name string) *FieldData {

	field := pMsg.Msg.Field(name)
	if field != nil {
		return pMsg.FieldDataMap[field.ID]
	}

	return nil

}

// GetById returns field data given its id
func (pMsg *ParsedMsg) GetById(id int) *FieldData {
	return pMsg.FieldDataMap[id]
}

// Copy returns a deep copy of the ParsedMsg
func (pMsg *ParsedMsg) Copy() *ParsedMsg {

	newParsedMsg := &ParsedMsg{IsRequest: false}
	newParsedMsg.FieldDataMap = make(map[int]*FieldData, len(pMsg.FieldDataMap))
	for id, fieldData := range pMsg.FieldDataMap {
		newParsedMsg.FieldDataMap[id] = fieldData.Copy()
	}

	newParsedMsg.Msg = pMsg.Msg
	newParsedMsg.MessageKey = pMsg.MessageKey

	return newParsedMsg

}

func parse(parserCfg *ParserConfig, buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	var err error
	switch field.Type {

	case FixedType:
		err = parseFixed(parserCfg, buf, parsedMsg, field)
	case VariableType:
		err = parseVariable(parserCfg, buf, parsedMsg, field)
	case BitmappedType:
		err = parseBitmap(parserCfg, buf, parsedMsg, field)
	default:
		return fmt.Errorf("libiso: Unsupported field type - %v", field.Type)

	}

	if err != nil {
		return err
	}

	switch field.Type {
	case FixedType, VariableType:

	case BitmappedType:
		{
			bitmap := parsedMsg.FieldDataMap[field.ID].Bitmap
			for _, cf := range field.Children {
				if bitmap.IsOn(cf.Position) {
					if err := parse(parserCfg, buf, parsedMsg, cf); err != nil {
						return err
					}
					bitmap.childData[cf.Position] = parsedMsg.FieldDataMap[cf.ID]
				}
			}
		}
	}

	return nil

}

func parseFixed(parserCfg *ParserConfig, buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	fieldData := &FieldData{Field: field}
	var err error

	if fieldData.Data, err = NextBytes(buf, field.Size); err != nil {
		return err
	}
	if field.Key {
		parsedMsg.MessageKey += fieldData.Value()
	}
	if parserCfg.LogEnabled {
		log.WithFields(log.Fields{"component": "parser"}).Debugf("Field %s, Length: %d, Value: %s\n", field.Name, field.Size, hex.EncodeToString(fieldData.Data))
	}
	parsedMsg.FieldDataMap[field.ID] = fieldData

	if field.HasChildren() {
		newBuf := bytes.NewBuffer(parsedMsg.Get(field.Name).Data)
		for _, cf := range field.Children {
			if err := parse(parserCfg, newBuf, parsedMsg, cf); err != nil {
				return err
			}
		}
	}

	return nil

}

func parseVariable(parserCfg *ParserConfig, buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	lenData, err := NextBytes(buf, field.LengthIndicatorSize)
	if err != nil {
		return err
	}
	var length uint64

	switch field.LengthIndicatorEncoding {
	case BINARY:
		{
			if field.LengthIndicatorSize > 4 {
				return ErrLargeLengthIndicator
			}

			switch field.LengthIndicatorSize {
			case 1:
				var byteLength uint8
				if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
					return err
				}
				length = uint64(byteLength)
			case 2:
				var byteLength uint16
				if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
					return err
				}
				length = uint64(byteLength)
			case 4:
				var byteLength uint32
				if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
					return err
				}
				length = uint64(byteLength)
			case 8:
				var byteLength uint64
				if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
					return err
				}
				length = byteLength
			default:
				{
					return errors.New(fmt.Sprint("Invalid length indicator size for binary (max 8) -", field.LengthIndicatorSize))

				}

			}

		}
	case BCD:
		if length, err = strconv.ParseUint(hex.EncodeToString(lenData), 10, 64); err != nil {
			return err
		}
	case ASCII:
		if length, err = strconv.ParseUint(string(lenData), 10, 64); err != nil {
			return err
		}
	case EBCDIC:
		if length, err = strconv.ParseUint(EBCDIC.EncodeToString(lenData), 10, 64); err != nil {
			return err
		}
	default:
		{
			return ErrInvalidEncoding
		}
	}

	if field.LengthIndicatorMultiplier == 2 && field.DataEncoding == BINARY {
		//special bcd field handling - https://github.com/rkbalgi/isosim/wiki/Variable-Fields
		if length%2 != 0 {
			length++
			length = length / 2
		}
	}

	fieldData := &FieldData{Field: field}
	if fieldData.Data, err = NextBytes(buf, int(length)); err != nil {
		return err
	}
	if field.Key {
		parsedMsg.MessageKey += fieldData.Value()
	}

	if parserCfg.LogEnabled {
		log.WithFields(log.Fields{"component": "parser"}).Debugf("Field %s, Length: %d, Value: %s\n", field.Name, length, hex.EncodeToString(fieldData.Data))
	}
	parsedMsg.FieldDataMap[field.ID] = fieldData

	if field.HasChildren() {
		newBuf := bytes.NewBuffer(parsedMsg.Get(field.Name).Data)
		for _, cf := range field.Children {
			if err := parse(parserCfg, newBuf, parsedMsg, cf); err != nil {
				return err
			}
		}
	}

	return nil

}

func parseBitmap(parserCfg *ParserConfig, buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	bitmap := NewBitmap()
	bitmap.field = field
	err := bitmap.parse(buf, parsedMsg, field)
	if err != nil {
		return err
	}
	if parserCfg.LogEnabled {
		log.WithFields(log.Fields{"component": "parser"}).Debugf("Field %s, Length: -, Value: %s\n", field.Name, bitmap.BinaryString())
	}
	parsedMsg.FieldDataMap[field.ID] = &FieldData{Field: field, Bitmap: bitmap}

	return nil

}

// NextBytes returns the next 'n' bytes from the Buffer
func NextBytes(buf *bytes.Buffer, n int) ([]byte, error) {

	if buf.Len() < n {
		return nil, ErrInsufficientData
	}
	nextData := buf.Next(n)
	return nextData, nil

}
