package utils

import "encoding/hex"

//StringToHex coverts a hex string into a []byte (binary)
func StringToHex(hexStr string) []byte {
	val, _ := hex.DecodeString(hexStr)
	return val
}

//HexToString converts a binary slice into a string
func HexToString(hexArray []byte) string {
	return hex.EncodeToString(hexArray)
}
