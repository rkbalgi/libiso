package utils
import "testing"

func Test_HexConversions(t *testing.T) {
	inp_val := "00f0ffe1"

	tmp := StringToHex(inp_val)
	str := HexToString(tmp)
	if !(str == inp_val) {
		t.Errorf("%s Failed.","Test_HexConversions");
	}
}

