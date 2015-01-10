package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_Parity(t *testing.T) {

	//var b, i uint8
	b := []byte{0x47, 0x02, 0xe2}

	//fmt.Println(hex.EncodeToString(b))
	to_odd_parity(b)
	if(hex.EncodeToString(b)!="4602e3"){
		t.Fail();
	}

}

