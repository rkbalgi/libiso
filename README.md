# A library of utilities related to payments, crypto, ISO8583 etc 


## Creating ISO8583 messages

First, create a yaml file containing the spec definition (see v2\iso8583\testdata) and then list that under a file called specs.yaml
(ignore the .spec files - they're an older way of defining specs)

1. Read all the specs defined (the path should contain the file specs.yaml)

```go
if err := iso8583.ReadSpecs(filepath.Join(".", "testdata")); err != nil {
		log.Fatal(err)
		return
}```

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


## UPDATE (04/22/2020)
1. Renaming the repo to libiso
2. Deleting paysim (concentrating efforts on isosim for now)

## UPDATE (03/22/2020)
1. There will be no further development on paysim, this repo will be strictly be used as a module/library
2. If you're interested in a ISO8583 simulator, please check out [ISO WebSim](https://github.com/rkbalgi/isosim)

## UPDATE (06/16/2019)
1. Folks developing on Windows please see this - https://github.com/rkbalgi/go/wiki/Building-on-Windows
2. Doesn't follow standard go style (coding conventions etc) - WIP
2. This has not be subject to any kind of targeted tests (performance or otherwise), so use this with a bit of caution - It's at the moment perhaps best suited for simulators


# Paysim
An open ISO8583 Simulator

<ul>
<li>The application is built using go and GTK+2 bindings made available at github.com/mattn/go-gtk (Thanks a ton Yashuhiro Matsumoto!)</li>
<li>The entire source code is available at https://github.com/rkbalgi/go</li>
<li>The interesting packages would be github.com/rkbalgi/go/execs/paysim, github.com/rkbalgi/go/paysim and github.com/rkbalgi/iso8583</li>
<li>There are loads of other interesting things available in other packages – like a minimalist implementation of a Thales HSM for basic commands (A6, MS, M6 and the like)
</li>
</ul>

You can read more about paysim here - https://github.com/rkbalgi/go/wiki/Paysim
