package hsm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

type TestStruct struct {
	Key string
}

type TestStruct2 struct {
	StringField         string
	DecimalNumberField1 uint
	DecimalNumberField2 uint
	HexNumberField      uint
	RawData             []byte
}

func Test_EncryptTripleLen_003(t *testing.T) {

	var keyHolder TestStruct
	buf := bytes.NewBuffer([]byte("T76D01994A48FF758D55467BC2CB68343372A29CB49D66843"))

	readKey(buf, &keyHolder.Key)
	if eKey, err := encryptKey(keyHolder.Key, "003"); err != nil {
		t.Error(err)
	} else {
		expected, _ := hex.DecodeString("9bfb11644c48c173c22deecb0bbe57352f11bcacba5c3c6d")
		if !bytes.Equal(eKey, expected) {
			t.Fail()
		}
	}

}

func Test_Read_SingleLength(t *testing.T) {

	var key_holder TestStruct
	buf := bytes.NewBuffer([]byte("e045fe3ea2a2f47d007a3030"))
	readKey(buf, &key_holder.Key)

	if key_holder.Key != "e045fe3ea2a2f47d" {
		t.Fail()
	}

}

func Test_read_SingleLength2(t *testing.T) {

	var key_holder TestStruct
	buf := bytes.NewBuffer([]byte("Ze045fe3ea2a2f47d007a3030"))
	readKey(buf, &key_holder.Key)

	if key_holder.Key != "Ze045fe3ea2a2f47d" {
		t.Log(key_holder.Key)
		t.Fail()
	}

}

func Test_read_key_DoubleLength(t *testing.T) {

	var key_holder TestStruct
	buf := bytes.NewBuffer([]byte("Ue045fe3ea2a2f47d007afe3ea2a2f47d007a3030"))
	readKey(buf, &key_holder.Key)
	if key_holder.Key != "Ue045fe3ea2a2f47d007afe3ea2a2f47d" {
		t.Fail()
	}

}

func Test_read_key_TripleLength(t *testing.T) {

	var key_holder TestStruct
	buf := bytes.NewBuffer([]byte("Te045fe3ea2a2f47d007afe3ea2a2f47d007a30307ee5eae300"))
	readKey(buf, &key_holder.Key)

	if key_holder.Key != "Te045fe3ea2a2f47d007afe3ea2a2f47d007a30307ee5eae3" {
		t.Fail()
	}

}

func Test_read_key_Invalid(t *testing.T) {

	var keyHolder TestStruct
	buf := bytes.NewBuffer([]byte("Ke045fe3ea2a2f47d007afe3ea2a2f47d007a30307ee5eae300"))

	result := readKey(buf, &keyHolder.Key)

	if result {
		t.Fail()
	}

}

func TestReadFieldsFromStruct(t *testing.T) {

	var fields_struct TestStruct2
	data, _ := hex.DecodeString("303030303030303030303032323334303041310001020304")
	buf := bytes.NewBuffer(data)

	result := readFixedField(buf, &fields_struct.StringField, 12, String)
	if !result {
		t.Fail()
	}
	result = readFixedField(buf, &fields_struct.DecimalNumberField1, 1, DecimalInt)
	if !result {
		t.Fail()
	}
	result = readFixedField(buf, &fields_struct.DecimalNumberField2, 2, DecimalInt)
	if !result {
		t.Fail()
	}
	result = readFixedField(buf, &fields_struct.HexNumberField, 4, HexadecimalInt)
	if !result {
		t.Fail()
	}

	result = readFixedField(buf, &fields_struct.RawData, 5, Binary)
	if !result {
		t.Fail()
	}

	fmt.Println(Dump(fields_struct))

	if fields_struct.StringField == "000000000002" &&
		fields_struct.DecimalNumberField1 == 2 &&
		fields_struct.DecimalNumberField2 == 34 &&
		fields_struct.HexNumberField == 161 {

		if hex.EncodeToString(fields_struct.RawData) != "0001020304" {
			t.Fail()

		}
	} else {
		t.Fail()
	}

}
