# A library of utilities related to payments, crypto, ISO8583 etc 


## Creating ISO8583 messages

First, create a yaml file containing the spec definition (see v2\iso8583\testdata) and then list that under a file called specs.yaml
(ignore the .spec files - they're an older way of defining specs)

1. Read all the specs defined (the path should contain the file specs.yaml)

```go
if err := iso8583.ReadSpecs(filepath.Join(".", "testdata")); err != nil {
		log.Fatal(err)
		return
}
```

2. Once initialized you can construct ISO8583 messages like below (from iso_test.go) -

```go
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
```
Please checkout https://gthub.com/rkbalgi/isosim project which uses this library.


### Note:
* Message Type and Bitmap are reserved keywords within this library (i.e you cannot call the Bitmap as Primary Bitmap or Bmp etc)
* This library has not yet been subjected to any kind of targeted tests (performance or otherwise), so use this with a bit of caution - It's at the moment perhaps best suited for simulators



# Paysim
Paysim is an old application that uses this library. You can read more about paysim here - https://github.com/rkbalgi/go/wiki/Paysim
