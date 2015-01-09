package hsm

var __hsm_debug_enabled bool = true

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
	header       string `size:"12"`
	command_name string `size:"2"`
}

type epilogue struct {
	delimiter             byte
	lmk_identifier        uint
	end_message_delimiter byte
	message_trailer       []byte
}

const (
	HSM_OK          = "00"
	HSM_PARSE_ERROR = "15"
)

const (
	ZMK_KEY_TYPE = "000"
	TMK_KEY_TYPE = "002"
)

