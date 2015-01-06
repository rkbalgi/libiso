package crypto

//implements padding schemes as defined in ISO/IEC 9797
// refer http://en.wikipedia.org/wiki/ISO/IEC_9797-1

import (
//	"encoding/hex"
//	"log"
//"fmt"
)

type PaddingType int

const (
	Iso9797M1Padding PaddingType = iota + 1
	Iso9797M2Padding
	DesBlockSize = 8
)

var __iso9797_pad_block []byte = []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func (paddingType PaddingType) Pad(data []byte) []byte {
	var paddedData []byte
	switch paddingType {
	case Iso9797M1Padding:
		{
			n := len(data)
			if n < DesBlockSize {
				n_pads := DesBlockSize - n
				padBytes := make([]byte, n_pads)
				//var padBytes [n_pads]byte;
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, padBytes...)
			} else if n == DesBlockSize || n%DesBlockSize == 0 {
				paddedData = data
			} else {
				n_pads := DesBlockSize - (n % DesBlockSize)
				padBytes := make([]byte, n_pads)
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, padBytes...)
			}

			break
		}

	case Iso9797M2Padding:
		{

			n := len(data)
			if n < DesBlockSize {
				n_pads := DesBlockSize - n
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, __iso9797_pad_block[:n_pads]...)
			} else {
				if n%DesBlockSize == 0 {
					paddedData = append(paddedData, data...)
					paddedData = append(paddedData, __iso9797_pad_block...)
				} else {
					n_pads := DesBlockSize - (n % DesBlockSize)
					paddedData = append(paddedData, data...)
					paddedData = append(paddedData, __iso9797_pad_block[:n_pads]...)
				}
			}

		}
	}

	return (paddedData)

}

func (pt PaddingType) RemovePad(padded_data []byte) []byte {

	var data []byte

	switch pt {
	case Iso9797M1Padding:
		{
			i := len(padded_data) - 1
			for padded_data[i] == 0x00 {
				i--
			}
			return padded_data[:i+1]

		}
	case Iso9797M2Padding:
		{
           i := len(padded_data) - 1
			for padded_data[i] != __iso9797_pad_block[0] {
				//fmt.Printf("%x\n",padded_data[i]);
				i--
			}
			return padded_data[:i]
		}

	}

	return (data)
}
