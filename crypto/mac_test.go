package crypto

import (
	"github.com/rkbalgi/utils"
	"testing"
	"bytes"
)


func Test_X919Mac_Test0(t *testing.T) {

	mac_data := utils.StringToHex("4E6F77206973207468652074696D6520666F7220616C6C20")
	expected_mac := utils.StringToHex("A1C72E74EA3FA9B6")
	computed_mac := GenerateMac_X919(mac_data, utils.StringToHex("0123456789ABCDEFFEDCBA9876543210"))

	if  !bytes.Equal(expected_mac,computed_mac) {
		t.Error( expected_mac, computed_mac)
		t.Fail()
	}

}

func Test_X919Mac_Test1(t *testing.T) {

	mac_data := utils.StringToHex("8155ADCC76B2FB0064F2C40037710477CE13C4BF75FD3DADF13B6D137AC1B915")
	expected_mac := utils.StringToHex("B2A45602664C486F")
	computed_mac := GenerateMac_X919(mac_data, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))

	if  !bytes.Equal(expected_mac,computed_mac) {
		t.Error( expected_mac, computed_mac)
		t.Fail()
	}

}

func Test_X919Mac_Test2(t *testing.T) {

	mac_data := utils.StringToHex("D3DADF13B6D137AC")
	expected_mac := utils.StringToHex("6A3307CAA14EB06A")
	computed_mac := GenerateMac_X919(mac_data, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))

	if  !bytes.Equal(expected_mac,computed_mac) {
		t.Error( expected_mac, computed_mac)
		t.Fail()
	}

}

func Test_X919Mac_Test3_WithPadding(t *testing.T) {

	mac_data := utils.StringToHex("D3DADF13B6D137ACDF13B6")
	expected_mac := utils.StringToHex("8B34B5BEAF1087A0")
	computed_mac := GenerateMac_X919(mac_data, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))

	if  !bytes.Equal(expected_mac,computed_mac) {
		t.Error( expected_mac, computed_mac)
		t.Fail()
	}

}

func Test_X919Mac_Test3_WithPadding2(t *testing.T) {

	mac_data := utils.StringToHex("D3DADF")
	expected_mac := utils.StringToHex("C66CA5DE6E324EDF")
	computed_mac := GenerateMac_X919(mac_data, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))

	if  !bytes.Equal(expected_mac,computed_mac) {
		t.Error( expected_mac, computed_mac)
		t.Fail()
	}

}
