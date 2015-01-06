package hsm

import (
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/net"
	"strings"
	"testing"
)

func fail_on_err(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
		t.Fail()
	}
}
func Test_Thales_MS(t *testing.T) {

	cmd_str := "303030303030303030303032;4D53;30;30;31;30;553831353541444343373642324642303036344632433430303337373130343737;30303043;000102030506070809070809"
	cmd_str = strings.Replace(cmd_str, ";", "", -1)
	msg_data, _ := hex.DecodeString(cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.MLI_2E)
	err := hsm_client.OpenConnection()
	fail_on_err(t, err)
	err = hsm_client.Write(msg_data)
	fail_on_err(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	fail_on_err(t, err)
	fmt.Println(hex.Dump(response_data))

}

func Test_Thales_MS_SingleLengthKey(t *testing.T) {

	cmd_str := "303030303030303030303032;4D53;30;30;31;30;44324337314130324431394542343233;30303043;000102030506070809070809"
	cmd_str = strings.Replace(cmd_str, ";", "", -1)
	msg_data, _ := hex.DecodeString(cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.MLI_2E)
	err := hsm_client.OpenConnection()
	fail_on_err(t, err)
	err = hsm_client.Write(msg_data)
	fail_on_err(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	fail_on_err(t, err)
	fmt.Println(hex.Dump(response_data))

}
