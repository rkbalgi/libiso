package hsm

import (
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/net"
	"strings"
	"testing"
)

func Test_Thales_NC(t *testing.T) {

	cmd_str := "303030303030303030303032;4e43;"
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
