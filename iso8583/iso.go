package iso8583

import (
	"time"
)

const (
	// IsoMessageType is a constant that indicates the Message Type or the MTI
	// (This name has special meaning within the context of ISO8583 and cannot be named anything else. The same restrictions apply for 'Bitmap')
	IsoMessageType = "Message Type"
	// IsoBitmap is a constant for the bitmap in the ISO message
	IsoBitmap = "Bitmap"
)

// Iso is a handle into accessing the details of a ISO message(via the parsedMsg)
type Iso struct {
	parsedMsg *ParsedMsg
}

// MetaInfo provides additional information about an operation performed
// For example, in response to a parse or assemble op
type MetaInfo struct {
	// OpTime is time taken by an operation
	OpTime time.Duration
	// MessageKey is a key that can be used to uniquely identify a transaction
	MessageKey string
}

// NewIso returns a Iso instance that can be used to build messages
func (msg *Message) NewIso() *Iso {
	isoMsg := FromParsedMsg(&ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData)})
	return isoMsg
}

// FromParsedMsg constructs a new Iso from a parsedMsg
func FromParsedMsg(parsedMsg *ParsedMsg) *Iso {
	isoMsg := &Iso{parsedMsg: parsedMsg}

	bmpField := parsedMsg.Msg.fieldByName[IsoBitmap]

	//if the bitmap field is not set then initialize it to a empty bitmap
	if _, ok := parsedMsg.FieldDataMap[bmpField.ID]; !ok {
		bmpFieldData := &FieldData{Field: bmpField, Bitmap: emptyBitmap(parsedMsg)}
		isoMsg.parsedMsg.FieldDataMap[bmpField.ID] = bmpFieldData
	}

	return isoMsg

}

// Set sets a field to the supplied value
func (iso *Iso) Set(fieldName string, value string) error {

	field := iso.parsedMsg.Msg.Field(fieldName)
	if field == nil {
		return ErrUnknownField
	}

	bmpField := iso.parsedMsg.Get(IsoBitmap)
	if field.ParentId == bmpField.Field.ID {
		iso.Bitmap().SetOn(field.Position)
		iso.Bitmap().Set(field.Position, value)
	} else {
		fieldData, err := field.ValueFromString(value)
		if err != nil {
			return err
		}
		iso.parsedMsg.FieldDataMap[field.ID] = &FieldData{Field: field, Data: fieldData}

	}

	return nil

}

// Get returns a field by its name
func (iso *Iso) Get(fieldName string) *FieldData {

	field := iso.parsedMsg.Msg.Field(fieldName)
	return iso.parsedMsg.FieldDataMap[field.ID]

}

// Bitmap returns the Bitmap from the Iso message
func (iso *Iso) Bitmap() *Bitmap {
	field := iso.parsedMsg.Msg.Field(IsoBitmap)
	fieldData := iso.parsedMsg.FieldDataMap[field.ID].Bitmap
	if fieldData != nil && fieldData.parsedMsg == nil {
		fieldData.parsedMsg = iso.parsedMsg
	}
	return fieldData

}

// ParsedMsg returns the backing parsedMsg
func (iso *Iso) ParsedMsg() *ParsedMsg {
	return iso.parsedMsg
}
