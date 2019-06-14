package utils

import "testing"

func Test_HexConversions(t *testing.T) {
	inpVal := "00f0ffe1"

	tmp := StringToHex(inpVal)
	str := HexToString(tmp)
	if !(str == inpVal) {
		t.Errorf("%s Failed.", "Test_HexConversions")
	}
}
