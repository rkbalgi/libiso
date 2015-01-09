package hsm

import(
	"strings"
	"bytes"
	"encoding/hex"
)
//converts a series of semicolon separated hsm command strings into a
//array of bytes

//fields can also be hexadecimal digit strings enclosed between x' and ')
//example - "ABC;DEF;0002;x'0000F0F1' etc
func format_hsm_command(hsm_cmd_str string) []byte {
	sub_fields := strings.Split(hsm_cmd_str, ";")

	buf := bytes.NewBuffer(make([]byte, 0))

	for _, sub_field := range sub_fields {
		if strings.HasPrefix(sub_field, "x'") && strings.HasSuffix(sub_field, "'") {
			data, err := hex.DecodeString(sub_field[2 : len(sub_field)-1])
			if err != nil {
				panic(err.Error())
			}
			buf.Write(data)
		} else {
			buf.Write([]byte(sub_field))
		}
	}

	return buf.Bytes()

}