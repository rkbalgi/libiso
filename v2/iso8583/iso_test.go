package iso8583_test

import (
	"encoding/hex"
	"github.com/rkbalgi/libiso/v2/iso8583"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
	if err := iso8583.ReadSpecs(filepath.Join(".", "testdata")); err != nil {
		log.Fatal(err)
		return
	}
}

func Test_ParseAndAssemble_Iso8585_Test_Spec(t *testing.T) {

	specName := "ISO8583-Test"
	spec := iso8583.SpecByName(specName)
	if spec == nil {
		t.Fatal("Unable to find spec - " + specName)
	}

	// Parse a message using an existing hex-dump

	msgData, _ := hex.DecodeString("3131303070386000000080003136343736363937373635343332373737373030343030303030303030303030313039303636363535313230313333353035323239333131333336383236")

	msg := spec.FindTargetMsg(msgData) // if you know the kind of message you are parse, you can do this - Example: spec.MessageByName("1100 - Authorization")
	if msg != nil {
		parsedMsg, err := msg.Parse(msgData)
		if err != nil {
			t.Fatal(err)
		} else {
			iso := iso8583.FromParsedMsg(parsedMsg)
			assert.Equal(t, "000000001090", iso.Bitmap().Get(4).Value())
			assert.Equal(t, "666551", iso.Bitmap().Get(11).Value())
		}
	} else {
		t.Fatal("Unable to derive the type of message the data represents")
	}

	// OR
	// build a message from scratch

	msg = spec.MessageByName("1100 - Authorization")
	iso := msg.NewIso()
	iso.Set("Message Type", "1100")
	iso.Bitmap().Set(3, "004000")
	iso.Bitmap().Set(4, "4766977654327777") // or iso.Set("PAN","4766977654327777")
	iso.Bitmap().Set(3, "004000")

	iso.Bitmap().Set(49, "336")
	iso.Bitmap().Set(50, "826")

	msgData, _, err := iso.Assemble()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "31313030300000000000c00030303430303034373636393737363534333237373737333336383236", hex.EncodeToString(msgData))

}

func Test_ParseMsg(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	msgData, _ := hex.DecodeString("31313030fe000000000003100000002000000000313233F4F5F6123678abcdef3132333435363738313568656c6c6f0003616263776f726c6400120102030405060708090a0b0c000568656c6c6ff0f0f60009f1a3b2c13032f1f2")

	spec := iso8583.SpecByName("TestSpec")
	if spec != nil {
		defaultMsg := spec.MessageByName("Default Message")
		parsedMsg, err := defaultMsg.Parse(msgData)
		if err != nil {
			t.Fatal("Test Failed. Error = " + err.Error())
			return
		}

		assert.Equal(t, "1100", parsedMsg.Get("Message Type").Value())
		bmp := parsedMsg.Get(iso8583.IsoBitmap).Bitmap
		assert.Equal(t, "hello", bmp.Get(56).Value())

		//sub fields of fixed field

		assert.Equal(t, "1234", parsedMsg.Get("SF6_1").Value())
		assert.Equal(t, "12", parsedMsg.Get("SF6_1_1").Value())
		assert.Equal(t, "34", parsedMsg.Get("SF6_1_2").Value())
		assert.Equal(t, "56", parsedMsg.Get("SF6_2").Value())
		assert.Equal(t, "78", parsedMsg.Get("SF6_3").Value())

		assert.Equal(t, "68656c6c6f0003616263776f726c64", bmp.Get(7).Value())

		//sub fields of variable field
		assert.Equal(t, "hello", parsedMsg.Get("SF7_1").Value())
		assert.Equal(t, "abc", parsedMsg.Get("SF7_2").Value())
		assert.Equal(t, "world", parsedMsg.Get("SF7_3").Value())

	} else {
		t.Fatal("No spec : TestSpec")
	}

}

func Test_AssembleMsg(t *testing.T) {

	spec := iso8583.SpecByName("TestSpec")
	if spec != nil {

		msg := spec.MessageByName("Default Message")
		if msg == nil {
			t.Fatal("msg is nil")
			return
		}
		isoMsg := msg.NewIso()

		//setting directly
		isoMsg.Set(iso8583.IsoMessageType, "1100")
		isoMsg.Set("Fixed2_ASCII", "123")
		isoMsg.Set("Fixed3_EBCDIC", "456")
		isoMsg.Set("FxdField6_WithSubFields", "12345678")
		isoMsg.Set("VarField7_WithSubFields", "68656c6c6f0003616263776f726c64")

		//setting via bitmap
		isoMsg.Bitmap().Set(56, "hello_iso")
		isoMsg.Bitmap().Set(60, "0987aefe")
		isoMsg.Bitmap().Set(91, "field91")

		if assembledMsg, _, err := isoMsg.Assemble(); err != nil {
			t.Fatal(err)
			return
		} else {
			t.Log(hex.EncodeToString(assembledMsg))
			assert.Equal(t, "31313030e6000000000001100000002000000000313233f4f5f63132333435363738313568656c6c6f0003616263776f726c64000968656c6c6f5f69736ff0f0f40987aefe30378689859384f9f1",
				hex.EncodeToString(assembledMsg))
			assert.True(t, isoMsg.Bitmap().IsOn(6))
			assert.True(t, isoMsg.Bitmap().IsOn(7))
			assert.True(t, isoMsg.Bitmap().IsOn(56))
			assert.True(t, isoMsg.Bitmap().IsOn(60))
			assert.True(t, isoMsg.Bitmap().IsOn(91))
		}
	} else {
		t.Fatal("No spec : TestSpec")
	}

}
