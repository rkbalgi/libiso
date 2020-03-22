package hsm

import (
	//"bytes"
	"encoding/hex"
	"fmt"
	"go/net"
	"strings"
	"testing"
)

func failOnErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func Test_ThalesTripleLength(t *testing.T) {

	hsmCmdStr := "000000000011;MS;0;0;1;0;T9bfb11644c48c173c22deecb0bbe57352f11bcacba5c3c6d;000c;x'000102030506070809070809';%00;x'19';ABCDEFGH"
	msgData := formatHsmCommand(hsmCmdStr)

	fmt.Println(hex.Dump(msgData))

	hsmClient := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsmClient.OpenConnection()
	failOnErr(t, err)
	defer hsmClient.Close()
	err = hsmClient.Write(msgData)
	failOnErr(t, err)

	responseData, err := hsmClient.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(responseData))

}

func Test_Thales_MS(t *testing.T) {

	cmdStr := "303030303030303030303032;4D53;30;30;31;30;553831353541444343373642324642303036344632433430303337373130343737;30303043;000102030506070809070809"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData, _ := hex.DecodeString(cmdStr)

	fmt.Println(hex.Dump(msgData))

	hsmClient := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsmClient.OpenConnection()
	failOnErr(t, err)
	defer hsmClient.Close()
	err = hsmClient.Write(msgData)
	failOnErr(t, err)

	responseData, err := hsmClient.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(responseData))
	//hsm_client.Close();

}

func Test_Thales_MS_SingleLengthKey(t *testing.T) {

	cmdStr := "303030303030303030303032;4D53;30;30;31;30;44324337314130324431394542343233;30303043;000102030506070809070809"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData, _ := hex.DecodeString(cmdStr)

	fmt.Println(hex.Dump(msgData))

	hsmClient := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsmClient.OpenConnection()
	failOnErr(t, err)
	defer hsmClient.Close()
	err = hsmClient.Write(msgData)
	failOnErr(t, err)

	responseData, err := hsmClient.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(responseData))

}
