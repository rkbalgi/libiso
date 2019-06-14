package hsm

var hsmDebugEnabled bool = true

const (
	String = iota + 1
	Binary
	DecimalInt
	HexadecimalInt
)

type EncodingType int

const (
	AsciiEncoding = iota + 1
	EbcdicEncoding
)

type prologue struct {
	header      string `size:"12"`
	commandName string `size:"2"`
}

type epilogue struct {
	delimiter           byte
	lmkIdentifier       uint
	endMessageDelimiter byte
	messageTrailer      []byte
}

const (
	HSM_OK          = "00"
	HSM_PARSE_ERROR = "15"
)

const (
	ZMK_KEY_TYPE = "000"
	TMK_KEY_TYPE = "002"
)
