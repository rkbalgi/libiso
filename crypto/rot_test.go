package crypto

import "testing"

//import "github.com/rkbalgi/crypto"

func Test_RotN_3(t *testing.T) {

	tmp := RotN(3, "xyzragha")
	if !(tmp == "abcudjkd") {
		t.Error("Test_RotN1() failed")
	}
}

func Test_RotN_InvalidChar(t *testing.T) {
	defer func() {
		str := recover()
		if str == nil {
			t.Logf("Test_RotN_InvalidChar() failed - %v", str)
		}
	}()
	RotN(3, "xyzragha!")

}
