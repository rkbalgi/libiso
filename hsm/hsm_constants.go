package hsm

var hsmDebugEnabled bool = true

const (
	// String
	String = iota + 1
	// Binary
	Binary
	// DecimalInt
	DecimalInt
	// HexadecimalInt
	HexadecimalInt
)

type EncodingType int

const (
	// AsciiEncoding is ASCII encoding
	AsciiEncoding = iota + 1
	// EbcdicEncoding
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
	// HSM_OK is OK response from Thales HSM
	HSM_OK = "00"
	// HSM_PARSE_ERROR implies that the HSM command was malformed
	HSM_PARSE_ERROR = "15"
)

const (
	// ZMK_KEY_TYPE represents a Zone Master Key
	ZMK_KEY_TYPE = "000"
	//TMK_KEY_TYPE represents a Terminal Master Key
	TMK_KEY_TYPE = "002"
)
