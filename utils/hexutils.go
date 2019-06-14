package utils

//import "strconv"
//import "bytes"
//import "fmt"
import "encoding/hex"

/**
converts a string containing hexadecimal digits into a byte array of unsigned bytes
**/
func StringToHex(hexStr string) []byte {

	val, _ := hex.DecodeString(hexStr)
	return val

	/*var hex_data []byte = make([]byte, len(hex_str)/2)
	j := 0

	for i := 0; i < len(hex_str); i += 2 {
		//var err error
		var tmp uint64
		tmp, _ = strconv.ParseUint(hex_str[i:i+2], 16, 8)
		hex_data[j] = byte(tmp)
		j++
	}

	return (hex_data)*/
}

/**
convert a byte array of unsigned bytes to a String
**/
func HexToString(hexArray []byte) string {

	return hex.EncodeToString(hexArray)

	/*var buf bytes.Buffer

	for i := 0; i < len(hex_array); i++ {
		fmt.Fprintf(&buf, "%02x", uint8(hex_array[i]))
		//buf.WriteString(strconv.FormatUint(uint64(hex_array[i]), 16))

	}
	return buf.String()
	*/
}
