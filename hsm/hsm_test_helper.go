package hsm

import (
	"bytes"
	"encoding/hex"
	"log"
	"strings"
)

//converts a series of semicolon separated hsm command strings into a
//array of bytes

//fields can also be hexadecimal digit strings enclosed between x' and ')
//example - "ABC;DEF;0002;x'0000F0F1' etc
func formatHsmCommand(hsmCmdStr string) []byte {
	subFields := strings.Split(hsmCmdStr, ";")

	buf := bytes.NewBuffer(make([]byte, 0))

	for _, subField := range subFields {
		if strings.HasPrefix(subField, "x'") && strings.HasSuffix(subField, "'") {
			data, err := hex.DecodeString(subField[2 : len(subField)-1])
			if err != nil {
				log.Print(err)
				break
			}
			buf.Write(data)
		} else {
			buf.Write([]byte(subField))
		}
	}

	return buf.Bytes()

}
