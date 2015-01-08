package hsm

import (
	"encoding/hex"
	"log"
)

//for double length keys
var __variant_dbl_len_1 byte = 0xa6
var __variant_dbl_len_2 byte = 0x5a

//for triple length keys
var __variant_triple_len_1 byte = 0x6a
var __variant_triple_len_2 byte = 0xde
var __variant_triple_len_3 byte = 0x2b

//lmk variants (position 0 is dummy - only 1 through 9 are used)
var __variants = [...]byte{0x00, 0xa6, 0x5a, 0x6a, 0xde, 0x2b, 0x50, 0x74, 0x9c, 0xfa}

//test lmk keys

var __lmk_0_1__, _ = hex.DecodeString("01010101010101017902CD1FD36EF8BA")
var __lmk_2_3__, _ = hex.DecodeString("20202020202020203131313131313131")
var __lmk_4_5__, _ = hex.DecodeString("40404040404040405151515151515151")
var __lmk_6_7__, _ = hex.DecodeString("61616161616161617070707070707070")
var __lmk_8_9__, _ = hex.DecodeString("80808080808080809191919191919191")
var __lmk_10_11__, _ = hex.DecodeString("A1A1A1A1A1A1A1A1B0B0B0B0B0B0B0B0")
var __lmk_12_13__, _ = hex.DecodeString("C1C1010101010101D0D0010101010101")
var __lmk_14_15__, _ = hex.DecodeString("E0E0010101010101F1F1010101010101")
var __lmk_16_17__, _ = hex.DecodeString("1C587F1C13924FEF0101010101010101")
var __lmk_18_19__, _ = hex.DecodeString("01010101010101010101010101010101")
var __lmk_20_21__, _ = hex.DecodeString("02020202020202020404040404040404")
var __lmk_22_23__, _ = hex.DecodeString("07070707070707071010101010101010")
var __lmk_24_25__, _ = hex.DecodeString("13131313131313131515151515151515")
var __lmk_26_27__, _ = hex.DecodeString("16161616161616161919191919191919")
var __lmk_28_29__, _ = hex.DecodeString("1A1A1A1A1A1A1A1A1C1C1C1C1C1C1C1C")
var __lmk_30_31__, _ = hex.DecodeString("23232323232323232525252525252525")
var __lmk_32_33__, _ = hex.DecodeString("26262626262626262929292929292929")
var __lmk_34_35__, _ = hex.DecodeString("2A2A2A2A2A2A2A2A2C2C2C2C2C2C2C2C")
var __lmk_36_37__, _ = hex.DecodeString("2F2F2F2F2F2F2F2F3131313131313131")
var __lmk_38_39__, _ = hex.DecodeString("01010101010101010101010101010101")

var key_type_table map[string][]byte

func init() {
	key_type_table = make(map[string][]byte, 2)
	key_type_table["00"] = __lmk_4_5__
	key_type_table["01"] = __lmk_6_7__
	key_type_table["02"] = __lmk_14_15__
	key_type_table["03"] = __lmk_16_17__
	key_type_table["04"] = __lmk_18_19__
	key_type_table["05"] = __lmk_20_21__
	key_type_table["06"] = __lmk_22_23__
	key_type_table["07"] = __lmk_24_25__
	key_type_table["08"] = __lmk_26_27__
	key_type_table["09"] = __lmk_28_29__
	//
	key_type_table["0A"] = __lmk_30_31__
	key_type_table["0B"] = __lmk_32_33__
	key_type_table["0C"] = __lmk_34_35__
	key_type_table["0D"] = __lmk_36_37__
	key_type_table["0E"] = __lmk_38_39__

	log.Printf("[%d] lmk keys loaded.", len(key_type_table))
}
