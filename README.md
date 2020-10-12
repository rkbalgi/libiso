# libiso - A library of utilities related to payments, crypto, ISO8583 etc 

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rkbalgi/libiso?sort=semver)
![CircleCI](https://img.shields.io/circleci/build/github/rkbalgi/libiso/master?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/rkbalgi/libiso)](https://goreportcard.com/report/github.com/rkbalgi/libiso)
[![codecov](https://codecov.io/gh/rkbalgi/libiso/branch/master/graph/badge.svg)](https://codecov.io/gh/rkbalgi/libiso)
![GitHub go.mod Go version (branch & subfolder of monorepo)](https://img.shields.io/github/go-mod/go-version/rkbalgi/libiso/master?filename=go.mod)

## Creating ISO8583 messages

First, create a yaml file containing the spec definition (for example, see [this](https://github.com/rkbalgi/libiso/blob/master/v2/iso8583/testdata/iso_specs.yaml)) and then list that under a file called [specs.yaml](https://github.com/rkbalgi/libiso/blob/master/v2/iso8583/testdata/specs.yaml)
(ignore the .spec files - they're an older way of defining specs)

1. Read all the specs defined (the path should contain the file specs.yaml)

```go
if err := iso8583.ReadSpecs(filepath.Join(".", "testdata")); err != nil {
		log.Fatal(err)
		return
}
```

2. Once initialized you can construct ISO8583 messages like below (from https://github.com/rkbalgi/libiso/blob/master/v2/iso8583/iso_test.go#L20) -

```go
	specName := "ISO8583-Test"
	spec := iso8583.SpecByName(specName)
	if spec == nil {
		t.Fatal("Unable to find spec - " + specName)
	}

	// Parse a message using an existing hex-dump

	msgData, _ := hex.DecodeString("3131303070386000000080003136343736363937373635343332373737373030343030303030303030303030313039303636363535313230313333353035323239333131333336383236")

	msg := spec.FindTargetMsg(msgData) // if you know the kind of message you are parsing, you can do this - Example: spec.MessageByName("1100 - Authorization")
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
Please checkout https://github.com/rkbalgi/isosim project which uses this library.

## Benchmarks
With v2.0.1 you can turn off logging (and hence gain some speed and lower allocations) using the new parser API

```go
        parser := iso8583.NewParser(&iso8583.ParserConfig{LogEnabled: false})

        log.SetLevel(log.ErrorLevel)

	specName := "ISO8583-Test"
	spec := iso8583.SpecByName(specName)
	if spec == nil {
		b.Fatal("Unable to find spec - " + specName)
	}
	msgData, _ := hex.DecodeString("3131303070386000000080003136343736363937373635343332373737373030343030303030303030303030313039303636363535313230313333353035323239333131333336383236")

	msg := spec.FindTargetMsg(msgData) // if you know the kind of message you are parse, you can do this - Example: spec.MessageByName("1100 - Authorization")
	parsedMsg, err := parser.Parse(msg,msgData)
	iso := iso8583.FromParsedMsg(parsedMsg)
	assert.Equal(t, "000000001090", iso.Bitmap().Get(4).Value())

```
```
PS C:\Users\rkbal\IdeaProjects\libiso\v2\iso8583> go test -bench . -run Benchmark_Parse
time="2020-10-11T09:56:02+05:30" level=debug msg="Available spec files -  [isoSpecs.spec iso_specs.yaml sample_spec.yaml]"
time="2020-10-11T09:56:02+05:30" level=debug msg="Reading file .. isoSpecs.spec"
time="2020-10-11T09:56:02+05:30" level=debug msg="Reading file .. iso_specs.yaml"
time="2020-10-11T09:56:02+05:30" level=debug msg="Reading file .. sample_spec.yaml"
goos: windows
goarch: amd64
pkg: github.com/rkbalgi/libiso/v2/iso8583
Benchmark_ParseWithParserAPI-8            327625              3692 ns/op            4016 B/op         27 allocs/op
Benchmark_ParseWithMsg-8                   85014             14037 ns/op           12121 B/op        154 allocs/op
PASS
ok      github.com/rkbalgi/libiso/v2/iso8583    4.600s
PS C:\Users\rkbal\IdeaProjects\libiso\v2\iso8583>

```
Just to see the impact of logging , with log level turned to TRACE - 
```
Benchmark_ParseWithMsg-8                   502           2355728 ns/op           24749 B/op        409 allocs/op
```

Also, a new API for assembling
```go
			asm:=iso8583.NewAssembler(&iso8583.AssemblerConfig{
			  LogEnabled: false,
		    })

			iso := msg.NewIso()
			iso.Set("Message Type", "1100")
			iso.Bitmap().Set(3, "004000")
			iso.Bitmap().Set(4, "4766977654327777")
			iso.Bitmap().Set(3, "004000")

			iso.Bitmap().Set(49, "336")
			iso.Bitmap().Set(50, "826")

			_, _, err := asm.Assemble(iso)
```

### Note:
* Message Type and Bitmap are reserved keywords within this library (i.e you cannot call the Bitmap as Primary Bitmap or Bmp etc)
* This library has not yet been subjected to any kind of targeted tests (performance or otherwise), so use this with a bit of caution - It's at the moment perhaps best suited for simulators



# Paysim
Paysim is an old application that uses this library. You can read more about paysim here - https://github.com/rkbalgi/go/wiki/Paysim
