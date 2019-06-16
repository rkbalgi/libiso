package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

type BitMap struct {
	//primary secondary and tertiary bitmaps
	_1bmp        *big.Int
	_2bmp        *big.Int
	_3bmp        *big.Int
	subFieldDef  []IsoField
	subFieldData []*FieldData
	id           int
	isoMsgDef    *MessageDef
}

func NewBitMap() *BitMap {

	bitmap := new(BitMap)
	bitmap._1bmp = big.NewInt(0)
	bitmap._2bmp = big.NewInt(0)
	bitmap._3bmp = big.NewInt(0)
	bitmap.subFieldDef = make([]IsoField, 128+1)
	bitmap.subFieldData = make([]*FieldData, 128+1)

	return bitmap
}

type BitmappedField interface {
	Parse(isoMsg *Iso8583Message, buf *bytes.Buffer)
	IsOn(int) bool
	SetOn(int)
	SetOff(int)
	SetId(int)
	GetId() int
	SetSpec(isoMsgDef *MessageDef)
	Bytes() []byte
	Def() string
}

func (bmp *BitMap) Def() string {
	return fmt.Sprintf("Name: Bitmap ; Id: %04d", bmp.GetId())
}

//add a new  fixed field to the bitmap
func (bmp *BitMap) addFixedField(
	pBmpPos int,
	pName string,
	pDataEncoding int,
	pFieldLen int) {

	field := NewFixedFieldDef(pName, pDataEncoding, pFieldLen)
	field.SetId(bmp.isoMsgDef.nextFieldSeq())
	bmp.subFieldDef[pBmpPos] = field

}

//add a new  variable field to the bitmap
func (bmp *BitMap) addVariableField(pBmpPos int, pName string,
	pLenEncoding int,
	pDataEncoding int,
	pFieldLen int) {

	field := NewVariableFieldDef(pName, pLenEncoding, pDataEncoding, pFieldLen)
	field.SetId(bmp.isoMsgDef.nextFieldSeq())
	bmp.subFieldDef[pBmpPos] = field

}

//check if the bit at position 'pos' is
//set as 1
func (bmp *BitMap) IsOn(pos int) bool {
	pos = pos - 1
	iBmp, iPos := bmp.getBmpAndPos(pos)
	if iBmp.Bit(iPos) == 1 {
		return true
	} else {
		return false
	}

}

func (bmp *BitMap) getBmpAndPos(pos int) (*big.Int, int) {

	var targetBmp *big.Int
	var iPos int

	switch {
	case pos >= 0 && pos < 64:
		{
			targetBmp = bmp._1bmp
			iPos = 63 - pos

		}
	case pos >= 64 && pos < 128:
		{

			targetBmp = bmp._2bmp
			iPos = 127 - pos

		}
	case pos >= 128 && pos < 192:
		{
			targetBmp = bmp._3bmp
			iPos = 191 - pos
		}
	default:
		{
			log.Print(fmt.Sprint("error: invalid position in bitmap", pos))
		}

	}

	return targetBmp, iPos

}

func (bmp *BitMap) SetOff(pos int) {
	pos = pos - 1
	targetBmp, iPos := bmp.getBmpAndPos(pos)
	targetBmp.SetBit(targetBmp, iPos, 0)
}

func (bmp *BitMap) SetOn(pos int) {
	pos = pos - 1
	targetBmp, iPos := bmp.getBmpAndPos(pos)
	targetBmp.SetBit(targetBmp, iPos, 1)

}

func toOctet(inData []byte) []byte {

	//fmt.Println("to_octet -", hex.EncodeToString(in_data))
	if len(inData) != 8 {
		nPads := 8 - len(inData)
		withPads := make([]byte, nPads)
		withPads = append(withPads, inData...)

		return withPads
	}

	return inData

}

func (bmp *BitMap) Parse(isoMsg *Iso8583Message, buf *bytes.Buffer) {

	if buf.Len() >= 8 {
		tmp := make([]byte, 8)
		_, err := buf.Read(tmp)
		isoMsg.handleError(err)
		bmp._1bmp.SetBytes(tmp)

		if tmp[0]&0x80 == 0x80 {
			if buf.Len() >= 8 {
				tmp := make([]byte, 8)
				_, err := buf.Read(tmp)
				isoMsg.handleError(err)
				bmp._2bmp.SetBytes(tmp)

				if tmp[0]&0x80 == 0x80 {
					if buf.Len() >= 8 {
						tmp := make([]byte, 8)
						_, err := buf.Read(tmp)
						isoMsg.handleError(err)
						bmp._3bmp.SetBytes(tmp)

					} else {
						isoMsg.bufferUnderflowError("bitmap")
					}
				}

			} else {
				isoMsg.bufferUnderflowError("bitmap")
			}
		}

	}

	isoMsg.log.Printf("parsed bitmap: %s\n", hex.EncodeToString(bmp.Bytes()))

}

func (bmp *BitMap) Bytes() []byte {

	var bmpData []byte

	if bmp._2bmp.Uint64() != 0 {
		//second bmp exists
		bmp._1bmp.SetBit(bmp._1bmp, 63, 1)
	} else {
		bmp._1bmp.SetBit(bmp._1bmp, 63, 0)
	}
	if bmp._3bmp.Uint64() != 0 {
		//second bmp exists
		bmp._2bmp.SetBit(bmp._2bmp, 63, 1)
	} else {
		bmp._2bmp.SetBit(bmp._2bmp, 63, 0)
	}
	bmpData = toOctet(bmp._1bmp.Bytes())

	if bmp._1bmp.Bit(63) == 1 {
		tmp := toOctet(bmp._2bmp.Bytes())
		bmpData = append(bmpData, tmp...)
	}
	if bmp._2bmp.Bit(63) == 1 {
		tmp := toOctet(bmp._3bmp.Bytes())
		bmpData = append(bmpData, tmp...)
	}

	return bmpData

}

func (bmp *BitMap) bitString() string {

	buf := bytes.NewBufferString("")
	for i := 1; i < 129; i++ {
		if bmp.IsOn(i) {
			buf.Write([]byte("1"))
		} else {
			buf.Write([]byte("0"))
		}
	}

	return buf.String()

}

func (bmp *BitMap) copyBits(srcBmp *BitMap) {

	bmp._1bmp = big.NewInt(0).Set(srcBmp._1bmp)
	bmp._2bmp = big.NewInt(0).Set(srcBmp._2bmp)
	bmp._3bmp = big.NewInt(0).Set(srcBmp._3bmp)

}

func (bmp *BitMap) SetId(id int) {
	bmp.id = id
}

func (bmp *BitMap) GetId() int {
	return bmp.id
}

func (bmp *BitMap) SetSpec(isoMsgDef *MessageDef) {
	bmp.isoMsgDef = isoMsgDef
}
