package hsm

import (
	"encoding/hex"
	"log"
)

//for double length keys
var variantDblLen1 byte = 0xa6
var variantDblLen2 byte = 0x5a

//for triple length keys
var variantTripleLen1 byte = 0x6a
var variantTripleLen2 byte = 0xde
var variantTripleLen3 byte = 0x2b
var (
	//lmk variants (position 0 is dummy - only 1 through 9 are used)
	__variants = [...]byte{0x00, 0xa6, 0x5a, 0x6a, 0xde, 0x2b, 0x50, 0x74, 0x9c, 0xfa}
)

//test lmk keys

var lmk01, _ = hex.DecodeString("01010101010101017902CD1FD36EF8BA")
var lmk23, _ = hex.DecodeString("20202020202020203131313131313131")
var lmk45, _ = hex.DecodeString("40404040404040405151515151515151")
var lmk67, _ = hex.DecodeString("61616161616161617070707070707070")
var lmk89, _ = hex.DecodeString("80808080808080809191919191919191")
var lmk1011, _ = hex.DecodeString("A1A1A1A1A1A1A1A1B0B0B0B0B0B0B0B0")
var lmk1213, _ = hex.DecodeString("C1C1010101010101D0D0010101010101")
var lmk1415, _ = hex.DecodeString("E0E0010101010101F1F1010101010101")
var lmk1617, _ = hex.DecodeString("1C587F1C13924FEF0101010101010101")
var lmk1819, _ = hex.DecodeString("01010101010101010101010101010101")
var lmk2021, _ = hex.DecodeString("02020202020202020404040404040404")
var lmk2223, _ = hex.DecodeString("07070707070707071010101010101010")
var lmk2425, _ = hex.DecodeString("13131313131313131515151515151515")
var lmk2627, _ = hex.DecodeString("16161616161616161919191919191919")
var lmk2829, _ = hex.DecodeString("1A1A1A1A1A1A1A1A1C1C1C1C1C1C1C1C")
var lmk3031, _ = hex.DecodeString("23232323232323232525252525252525")
var lmk3233, _ = hex.DecodeString("26262626262626262929292929292929")
var lmk3435, _ = hex.DecodeString("2A2A2A2A2A2A2A2A2C2C2C2C2C2C2C2C")
var lmk3637, _ = hex.DecodeString("2F2F2F2F2F2F2F2F3131313131313131")
var lmk3839, _ = hex.DecodeString("01010101010101010101010101010101")

var keyTypeTable map[string][]byte

func init() {
	keyTypeTable = make(map[string][]byte, 2)
	keyTypeTable["00"] = lmk45
	keyTypeTable["01"] = lmk67
	keyTypeTable["02"] = lmk1415
	keyTypeTable["03"] = lmk1617
	keyTypeTable["04"] = lmk1819
	keyTypeTable["05"] = lmk2021
	keyTypeTable["06"] = lmk2223
	keyTypeTable["07"] = lmk2425
	keyTypeTable["08"] = lmk2627
	keyTypeTable["09"] = lmk2829
	//
	keyTypeTable["0A"] = lmk3031
	keyTypeTable["0B"] = lmk3233
	keyTypeTable["0C"] = lmk3435
	keyTypeTable["0D"] = lmk3637
	keyTypeTable["0E"] = lmk3839

	log.Printf("[%d] lmk keys loaded.", len(keyTypeTable))
}
