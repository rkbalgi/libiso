package utils

import "encoding/hex"

/**
converts a string containing hexadecimal digits into a byte array of unsigned bytes
**/
func StringToHex(hexStr string) []byte {

	val, _ := hex.DecodeString(hexStr)
	return val

}

/**
convert a byte array of unsigned bytes to a String
**/
func HexToString(hexArray []byte) string {

	return hex.EncodeToString(hexArray)

}
