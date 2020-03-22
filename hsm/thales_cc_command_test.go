package hsm

import (
	"encoding/hex"
	"fmt"
	"go/net"
	"strings"
	"testing"
)

func Test_Thales_CC_1(t *testing.T) {

	cmdStr := "000000000002;CC;U86AF65D8C29DC08C75D13FBDD88ABB0B;UCBDB34FC28BCA2EECD92F932C4433EC2;12;7FE8132B2F7F0D57;01;01;111111111111;"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData := []byte(cmdStr)

	fmt.Println(hex.Dump(msgData))

	hsmClient := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsmClient.OpenConnection()
	failOnErr(t, err)
	err = hsmClient.Write(msgData)
	failOnErr(t, err)
	defer hsmClient.Close()

	responseData, err := hsmClient.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(responseData))

}

func Test_Thales_CC_2(t *testing.T) {

	cmdStr := "000000000002;CC;2E1AB3C9C6A56939;UCBDB34FC28BCA2EECD92F932C4433EC2;12;9BE87D27C9A6C1B6;01;01;111111111111;"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData := []byte(cmdStr)

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

func Test_Thales_CC_3(t *testing.T) {

	cmdStr := "000000000002;CC;U86AF65D8C29DC08C75D13FBDD88ABB0B;2E1AB3C9C6A56939;12;7FE8132B2F7F0D57;01;01;111111111111;"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData := []byte(cmdStr)

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

func Test_Thales_CC_4(t *testing.T) {

	cmd_str := "000000000002;CC;U86AF65D8C29DC08C75D13FBDD88ABB0B;UCBDB34FC28BCA2EECD92F932C4433EC2;12;7FE8132B2F7F0D57;01;03;111111111111;"
	cmd_str = strings.Replace(cmd_str, ";", "", -1)
	msg_data := []byte(cmd_str)

	fmt.Println(hex.Dump(msg_data))

	hsm_client := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsm_client.OpenConnection()
	failOnErr(t, err)
	defer hsm_client.Close()

	err = hsm_client.Write(msg_data)
	failOnErr(t, err)

	response_data, err := hsm_client.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(response_data))

}
