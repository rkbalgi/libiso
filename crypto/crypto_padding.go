package crypto

// Implements padding schemes as defined in ISO/IEC 9797
// refer http://en.wikipedia.org/wiki/ISO/IEC_9797-1

// PaddingType identifies the padding type used when MAC'ing
type PaddingType int

const (
	// Iso9797M1Padding
	Iso9797M1Padding PaddingType = iota + 1
	// Iso9797M2Padding
	Iso9797M2Padding
	// DesBlockSize
	DesBlockSize = 8
)

var iso9797PadBlock []byte = []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func (paddingType PaddingType) Pad(data []byte) []byte {
	var paddedData []byte
	switch paddingType {
	case Iso9797M1Padding:
		{
			n := len(data)
			if n < DesBlockSize {
				nPads := DesBlockSize - n
				padBytes := make([]byte, nPads)
				//var padBytes [n_pads]byte;
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, padBytes...)
			} else if n == DesBlockSize || n%DesBlockSize == 0 {
				paddedData = data
			} else {
				nPads := DesBlockSize - (n % DesBlockSize)
				padBytes := make([]byte, nPads)
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, padBytes...)
			}

			break
		}

	case Iso9797M2Padding:
		{

			n := len(data)
			if n < DesBlockSize {
				nPads := DesBlockSize - n
				paddedData = append(paddedData, data...)
				paddedData = append(paddedData, iso9797PadBlock[:nPads]...)
			} else {
				if n%DesBlockSize == 0 {
					paddedData = append(paddedData, data...)
					paddedData = append(paddedData, iso9797PadBlock...)
				} else {
					nPads := DesBlockSize - (n % DesBlockSize)
					paddedData = append(paddedData, data...)
					paddedData = append(paddedData, iso9797PadBlock[:nPads]...)
				}
			}

		}
	}

	return paddedData

}

func (paddingType PaddingType) RemovePad(paddedData []byte) []byte {

	var data []byte

	switch paddingType {
	case Iso9797M1Padding:
		{
			i := len(paddedData) - 1
			for paddedData[i] == 0x00 {
				i--
			}
			return paddedData[:i+1]

		}
	case Iso9797M2Padding:
		{
			i := len(paddedData) - 1
			for paddedData[i] != iso9797PadBlock[0] {
				//fmt.Printf("%x\n",padded_data[i]);
				i--
			}
			return paddedData[:i]
		}

	}

	return data
}
