package ebcdic

import (
	"encoding/hex"
	"testing"
)

func Test_Ebcdic(t *testing.T) {

	//t.Log(ebcdic_to_ascii)
	data, _ := hex.DecodeString("f0f1f2f3f420202020f1f9c2")
	str := EncodeToString(data)
	t.Log(str, "\n")

	data = Decode("AGNS0001")
	t.Log(hex.EncodeToString(data), "\n")

}
