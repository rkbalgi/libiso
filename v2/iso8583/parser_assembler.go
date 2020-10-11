package iso8583

import (
	"bytes"
	"time"
)

// Parser defines a ISO8583 message parser
type Parser struct {
	cfg *ParserConfig
}

// ParserConfig defines the various configurations of the Parser
type ParserConfig struct {
	LogEnabled bool
}

// NewParser creates and returns a new parser
func NewParser(parserCfg *ParserConfig) *Parser {
	return &Parser{cfg: parserCfg}
}

// Parse parses msgData to the structure defined by msg
func (p *Parser) Parse(msg *Message, msgData []byte) (*ParsedMsg, *MetaInfo, error) {
	metaInfo := &MetaInfo{}
	tStart := time.Now()
	parsedMsg, err := parseWithConfig(p.cfg, msg, msgData)
	if err != nil {
		return nil, nil, err
	}
	metaInfo.OpTime = time.Since(tStart)
	return parsedMsg, metaInfo, nil
}

// Assembler defines a ISO8583 message assembler
type Assembler struct {
	cfg *AssemblerConfig
}

// AssemblerConfig defines all configuration options for the assembler
type AssemblerConfig struct {
	LogEnabled bool
}

// NewAssembler creates a new assembler
func NewAssembler(assemblerCfg *AssemblerConfig) *Assembler {
	return &Assembler{cfg: assemblerCfg}
}

func (asm *Assembler) Assemble(iso *Iso) ([]byte, *MetaInfo, error) {

	msg := iso.parsedMsg.Msg
	buf := new(bytes.Buffer)
	meta := &MetaInfo{}
	t1 := time.Now()
	for _, field := range msg.Fields {
		if err := asm.assemble(buf, meta, iso.parsedMsg, iso.parsedMsg.FieldDataMap[field.ID]); err != nil {
			return nil, nil, err
		}
	}

	meta.OpTime = time.Since(t1)

	return buf.Bytes(), meta, nil
}

// Assemble assembles the raw form of the message
func (iso *Iso) Assemble() ([]byte, *MetaInfo, error) {

	msg := iso.parsedMsg.Msg
	buf := new(bytes.Buffer)
	meta := &MetaInfo{}
	t1 := time.Now()
	for _, field := range msg.Fields {
		if err := NewAssembler(&AssemblerConfig{
			LogEnabled: false,
		}).assemble(buf, meta, iso.parsedMsg, iso.parsedMsg.FieldDataMap[field.ID]); err != nil {
			return nil, nil, err
		}
	}

	meta.OpTime = time.Since(t1)

	return buf.Bytes(), meta, nil

}
