package hsm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/net"
	"strings"
	"testing"
)

func fail_on_err(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
		t.Fail()
	}
}

//converts a series of semicolon separated hsm command strings into a
//array of bytes

//fields can also be hexadecimal digit strings enclosed between x' and ')
//example - "ABC;DEF;0002;x'0000F0F1' etc
func convert2ascii(hsm_cmd_str string) []byte {
	sub_fields := strings.Split(hsm_cmd_str, ";")

	buf := bytes.NewBuffer(make([]byte, 0))

	for _, sub_field := range sub_fields {
		if strings.HasPrefix(sub_field, "x'") && strings.HasSuffix(sub_field, "'") {
			data, err := hex.DecodeString(sub_field[2 : len(sub_field)-1])
			if err != nil {
				panic(err.Error())
			}
			buf.Write(data)
		} else {
			buf.Write([]byte(sub_field))
		}
	}

	return buf.Bytes()

}

func Test_ThalesTripleLength(t *testing.T) {

	hsm_cmd_str := "000000000011;MS;0;0;1;0;T9bfb11644c48c173c22deecb0bbe57352f11bcacba5c3c6d;000c;x'000102030506070809070809';%00;x'19';ABCDEFGH"
	msg_data := convert2ascii(hsm_cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.MLI_2E)
	err := hsm_client.OpenConnection()
	fail_on_err(t, err)
	defer hsm_client.Close();
	err = hsm_client.Write(msg_data)
	fail_on_err(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	fail_on_err(t, err)
	fmt.Println(hex.Dump(response_data))
	

}

func Test_Thales_MS(t *testing.T) {

	cmd_str := "303030303030303030303032;4D53;30;30;31;30;553831353541444343373642324642303036344632433430303337373130343737;30303043;000102030506070809070809"
	cmd_str = strings.Replace(cmd_str, ";", "", -1)
	msg_data, _ := hex.DecodeString(cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.MLI_2E)
	err := hsm_client.OpenConnection()
	fail_on_err(t, err)
	defer hsm_client.Close();
	err = hsm_client.Write(msg_data)
	fail_on_err(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	fail_on_err(t, err)
	fmt.Println(hex.Dump(response_data))
	//hsm_client.Close();

}

func Test_Thales_MS_SingleLengthKey(t *testing.T) {

	cmd_str := "303030303030303030303032;4D53;30;30;31;30;44324337314130324431394542343233;30303043;000102030506070809070809"
	cmd_str = strings.Replace(cmd_str, ";", "", -1)
	msg_data, _ := hex.DecodeString(cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.MLI_2E)
	err := hsm_client.OpenConnection()
	fail_on_err(t, err)
	defer hsm_client.Close();
	err = hsm_client.Write(msg_data)
	fail_on_err(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	fail_on_err(t, err)
	fmt.Println(hex.Dump(response_data))

}
