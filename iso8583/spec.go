package iso8583 // github.com/rkbalgi/isosim/iso

import (
	"bytes"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

var specMapMu sync.RWMutex
var specMap = make(map[string]*Spec, 10)

// Spec represents an ISO8583 specification
type Spec struct {
	Name string `yaml:"name"`
	ID   int    `yaml:"id"`

	HeaderFields []*Field `yaml:"header_fields"`

	Messages []*Message `yaml:"messages"`
}

// GetOrAddMsg returns (or adds and returns) a msg - This is usually called
// during initialization
func (spec *Spec) GetOrAddMsg(msgId int, msgName string) (*Message, bool) {

	if msg := spec.MessageByID(msgId); msg != nil {
		return msg, false
	}

	specMapMu.Lock()
	defer specMapMu.Unlock()

	msg := &Message{Name: msgName, ID: msgId,
		Fields:       make([]*Field, 0, 10),
		fieldByIdMap: make(map[int]*Field, 10),
		fieldByName:  make(map[string]*Field),
	}
	spec.Messages = append(spec.Messages, msg)

	return msg, true

}

// GetMessages returns a list of all messages defined for the spec
func (spec *Spec) GetMessages() []*Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	msgs := make([]*Message, 0, len(spec.Messages))
	for _, msg := range spec.Messages {
		msgs = append(msgs, msg)
	}
	return msgs
}

// MessageByID returns a message given its id
func (spec *Spec) MessageByID(msgId int) *Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	for _, msg := range spec.Messages {
		if msg.ID == msgId {
			return msg
		}

	}

	return nil
}

// MessageByName returns a message given its name
func (spec *Spec) MessageByName(msgName string) *Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	for _, msg := range spec.Messages {
		if msg.Name == msgName {
			return msg
		}
	}

	return nil

}

// FindTargetMsg parses any defined header fields and returns a message
func (spec *Spec) FindTargetMsg(data []byte) *Message {

	if spec.HeaderFields == nil || len(spec.HeaderFields) == 0 {
		return nil
	}

	matchKey := ""
	buf := bytes.NewBuffer(data)
	parsedMsg := &ParsedMsg{
		IsRequest:    false,
		Msg:          nil,
		FieldDataMap: make(map[int]*FieldData),
		MessageKey:   ""}

	for _, field := range spec.HeaderFields {
		if err := parse(&ParserConfig{LogEnabled: false}, buf, parsedMsg, field); err != nil {
			log.Errorf("Failed to parse header fields: %s, Error: %v", field.Name, err)
			return nil
		}
	}
	for _, field := range spec.HeaderFields {
		matchKey += parsedMsg.GetById(field.ID).Value()
	}

	if matchKey == "" {
		log.Errorf("libiso: No match key found!")
		return nil
	}

	for _, msg := range spec.Messages {
		for _, hm := range msg.HeaderMatch {
			if hm == matchKey {
				return msg
			}
		}
	}

	return nil

}

func printAllSpecsInfo() {

	buf := bytes.NewBufferString("")

	for specName, spec := range specMap {

		buf.WriteString(fmt.Sprintf("\nSpec = %s\n", specName))
		for _, msg := range spec.Messages {
			buf.WriteString(fmt.Sprintf("Spec Message = %s\n", msg.Name))
			level := 0
			for _, field := range msg.Fields {
				displayField(buf, field, level)

			}
		}
		log.Debugln(buf.String() + "\n")
		buf.Reset()
	}
}

func displayField(buf *bytes.Buffer, field *Field, level int) {
	i := 0
	for ; i < level; i++ {
		buf.WriteString("--")
	}
	buf.WriteString(">  ")
	buf.WriteString(fmt.Sprintf("Field : %v\n", field))
	if field.HasChildren() {
		for _, childField := range field.Children {
			displayField(buf, childField, level+1)
		}
	}
}

// AllSpecs returns a list of all defined specs
func AllSpecs() []*Spec {

	specs := make([]*Spec, 0, len(specMap))
	for _, spec := range specMap {
		specs = append(specs, spec)
	}
	return specs

}

// SpecByID returns a spec given it's id
func SpecByID(specId int) *Spec {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	for _, spec := range specMap {
		if spec.ID == specId {
			return spec
		}
	}
	return nil

}

// SpecByName returns a spec given its name
func SpecByName(specName string) *Spec {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	return specMap[specName]

}

func getOrCreateNewSpec(specId int, specName string) (spec *Spec, ok bool, err error) {

	spec = SpecByID(specId)
	if spec != nil {
		return nil, false,
			fmt.Errorf("libiso: SpecID - %d cannot be used for spec - %s. Is already used by %s", specId, specName, spec.Name)
	}
	specMapMu.Lock()
	defer specMapMu.Unlock()

	spec, ok = specMap[specName]
	if !ok {
		spec = &Spec{Name: specName, ID: specId, Messages: make([]*Message, 0)}
		specMap[specName] = spec
		return spec, true, nil
	}
	return spec, false, nil

}
