package keyblock

import (
	"encoding/hex"
	"strings"
)

var tripleDesLMK []byte //check value 165126
var aesLMK []byte       //9D04A0

func init() {
	tripleDesLMK, _ = hex.DecodeString(strings.Replace("01 23 45 67 89 AB CD EF 80 80 80 80 80 80 80 80 FE DC BA 98 76 54 32 10", " ", "", -1))
	aesLMK, _ = hex.DecodeString(strings.Replace("9B 71 33 3A 13 F9 FA E7 2F 9D 0E 2D AB 4A D6 78 47 18 01 2F 92 44 03 3F 3F 26 A2 DE 0C 8A A1 1A", " ", "", -1))
}
