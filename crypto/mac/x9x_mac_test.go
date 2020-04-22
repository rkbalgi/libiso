package mac

import (
	"bytes"
	"github.com/rkbalgi/libiso/utils"
	"testing"
)

func Test_X919Mac_Test0(t *testing.T) {

	var err error
	macData := utils.StringToHex("4E6F77206973207468652074696D6520666F7220616C6C20")
	expectedMac := utils.StringToHex("A1C72E74EA3FA9B6")
	computedMac, err := GenerateMacX919(macData, utils.StringToHex("0123456789ABCDEFFEDCBA9876543210"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expectedMac, computedMac) {
		t.Error(expectedMac, computedMac)
		t.Fail()
	}

}

func Test_X919Mac_Test1(t *testing.T) {

	var err error
	macData := utils.StringToHex("8155ADCC76B2FB0064F2C40037710477CE13C4BF75FD3DADF13B6D137AC1B915")
	expectedMac := utils.StringToHex("B2A45602664C486F")
	computedMac, err := GenerateMacX919(macData, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))

	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expectedMac, computedMac) {
		t.Error(expectedMac, computedMac)
		t.Fail()
	}

}

func Test_X919Mac_Test2(t *testing.T) {

	macData := utils.StringToHex("D3DADF13B6D137AC")
	expectedMac := utils.StringToHex("6A3307CAA14EB06A")
	if computedMac, err := GenerateMacX919(macData, utils.StringToHex("76850752AD7307ADE554D06D3BA73279")); err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(expectedMac, computedMac) {
			t.Error(expectedMac, computedMac)
			t.Fail()
		}
	}

}

func Test_X919Mac_Test3_WithPadding(t *testing.T) {

	var err error
	macData := utils.StringToHex("D3DADF13B6D137ACDF13B6")
	expectedMac := utils.StringToHex("8B34B5BEAF1087A0")
	computedMac, err := GenerateMacX919(macData, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expectedMac, computedMac) {
		t.Error(expectedMac, computedMac)
		t.Fail()
	}

}

func Test_X919Mac_Test3_WithPadding2(t *testing.T) {

	var err error
	macData := utils.StringToHex("D3DADF")
	expectedMac := utils.StringToHex("C66CA5DE6E324EDF")
	computedMac, err := GenerateMacX919(macData, utils.StringToHex("76850752AD7307ADE554D06D3BA73279"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expectedMac, computedMac) {
		t.Error(expectedMac, computedMac)
		t.Fail()
	}

}
