package hsm

import (
	"encoding/hex"
	"log"
	"strconv"
	"github.com/rkbalgi/go/crypto"
)

//for double length keys
var __variant_dbl_len_1 byte = 0xa6
var __variant_dbl_len_2 byte = 0x5a

//for triple length keys
var __variant_triple_len_1 byte = 0x6a
var __variant_triple_len_2 byte = 0xde
var __variant_triple_len_3 byte = 0x2b

//lmk variants
var __variants = [...]byte{0x00, 0xa6, 0x5a, 0x6a, 0xde, 0x2b, 0x50, 0x74, 0x9c, 0xfa}


//test lmk keys
var __lmk_6_7__, _ = hex.DecodeString("61616161616161617070707070707070")
var __lmk_16_17__, _ = hex.DecodeString("1C587F1C13924FEF0101010101010101")
var __lmk_26_27__, _ = hex.DecodeString("16161616161616161919191919191919")


var key_type_table map[string][]byte

//decrypt a key encrypted under the lmk

func decrypt_key(key_str string, key_type string) []byte {

	var key_data []byte
	if key_str[0] == 'U' || key_str[0] == 'Z' || key_str[0] == 'T' {
		key_data, _ = hex.DecodeString(key_str[1:])
	} else {
		key_data, _ = hex.DecodeString(key_str[:])
	}

	lmk_key := key_type_table[key_type[1:]]
	if lmk_key == nil {
		panic("unsupported key type" + key_type)
	}

	//no variant to be applied, just the usual
	v_key := make([]byte, 16)
	copy(v_key, lmk_key)

	if key_type[0] != '0' {
		//apply variant
		i, _ := strconv.ParseInt(key_type[:1], 10, 32)
		v_key[0] = v_key[0] ^ __variants[i]

	}

	if len(key_data) == 16 {
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		left_half := crypto.DecryptTripleDes(key_data[0:8], v_key)

		//now right half
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		v_key[8] = v_key[8] ^ __variant_dbl_len_2
		right_half := crypto.DecryptTripleDes(key_data[8:], v_key)

		clear_key := make([]byte, 0)
		clear_key = append(clear_key, left_half...)
		clear_key = append(clear_key, right_half...)

		return (clear_key)

	} else if len(key_data) == 8{
		//single length
		clear_key:=crypto.DecryptTripleDes(key_data, v_key)
		return(clear_key);
		
	}else{
		panic("not implemented for triple length keys")
	}

}

func init() {
	key_type_table = make(map[string][]byte, 2)
	key_type_table["01"] = __lmk_6_7__
	key_type_table["03"] = __lmk_16_17__
	key_type_table["08"] = __lmk_26_27__
	log.Printf("[%d] lmk keys loaded.",len(key_type_table));
}


