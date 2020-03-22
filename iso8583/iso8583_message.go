package iso8583

import (
	"bytes"
	"container/list"
	_ "encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"go/paysim"
	pylog "go/paysim/log"
	"log"
	"os"
	_ "reflect"
	_ "strconv"
)

type Iso8583Message struct {
	isoMsgDef     *MessageDef
	fieldDataList *list.List
	log           *log.Logger
	bitMap        *BitMap //for convenience
	nameToDataMap map[string]*FieldData
	idToDataMap   map[int]*FieldData
}

func (isoMsg *Iso8583Message) Bitmap() *BitMap {
	return isoMsg.bitMap
}

//GetMessageType returns the 'Message Type' as string
func (isoMsg *Iso8583Message) GetMessageType() string {
	return isoMsg.nameToDataMap["Message Type"].String()
}

//SpecName returns the name of the specification for this message
func (isoMsg *Iso8583Message) SpecName() string {
	return isoMsg.isoMsgDef.specName
}

func (isoMsg *Iso8583Message) ToWebMsg(isReq bool) *WebMsgData {

	jsonMsg := WebMsgData{}
	jsonMsg.Spec = isoMsg.isoMsgDef.specName
	if isReq {
		jsonMsg.Type = "Request"
	} else {
		jsonMsg.Type = "Response"
	}
	jsonMsg.DataArray = make([]string, isoMsg.isoMsgDef.fieldSeq)

	for l := isoMsg.isoMsgDef.fieldsDefList.Front(); l != nil; l = l.Next() {
		switch obj := l.Value.(type) {
		case IsoField:
			{
				isoField := isoMsg.GetFieldByName(obj.String())
				if isoField.fieldData != nil {
					jsonMsg.DataArray[isoField.fieldDef.GetId()] = isoField.String()
				}
			}
		case BitmappedField:
			{

				jsonMsg.DataArray[obj.GetId()] = isoMsg.bitMap.bitString()
				for fPos, fData := range isoMsg.bitMap.subFieldData {
					if fData != nil && fData.fieldData != nil && isoMsg.bitMap.IsOn(fPos) {
						jsonMsg.DataArray[fData.fieldDef.GetId()] = fData.String()
					}
				}

			} //end case

		} //end switch
	} //end for

	return &jsonMsg

}

//SetData sets data into individual fields by id
func (isoMsg *Iso8583Message) SetData(data []string) {

	for l := isoMsg.isoMsgDef.fieldsDefList.Front(); l != nil; l = l.Next() {

		switch obj := l.Value.(type) {
		case IsoField:
			{
				isoField := isoMsg.GetFieldByName(obj.String())
				isoField.SetData(data[isoField.fieldDef.GetId()])

			}
		case BitmappedField:
			{
				bitmapVal := data[obj.GetId()]
				for i := 0; i < len(bitmapVal); i++ {

					if bitmapVal[i:i+1] == "1" {
						isoMsg.bitMap.SetOn(i + 1)
					} else {
						isoMsg.bitMap.SetOff(i + 1)
					}
				}

				for fPos, fData := range isoMsg.bitMap.subFieldData {
					if fData != nil && isoMsg.bitMap.IsOn(fPos) {
						fData.SetData(data[fData.fieldDef.GetId()])
						//iso_msg.bit_map.SetOn(f_pos)
					}
				}

			}

		}
	}
}

//GetBinaryBitmap returns the 'Bitmap' as binary string
func (isoMsg *Iso8583Message) GetBinaryBitmap() string {

	binaryBmpStr := bytes.NewBufferString("")
	for i := 1; i < 129; i++ {
		if isoMsg.bitMap.IsOn(i) {
			binaryBmpStr.WriteString("1")
		} else {
			binaryBmpStr.WriteString("0")
		}
	}

	return binaryBmpStr.String()

}

//IsSelected returns a boolean indicating
//if the 'position' is selected in the bitmap
func (isoMsg *Iso8583Message) IsSelected(position int) bool {
	return isoMsg.bitMap.IsOn(position)
}

//GetFieldData returns the data associated with the 'position'
//in the iso_msg
func (isoMsg *Iso8583Message) GetFieldData(position int) (data string, err error) {
	fieldData, err := isoMsg.Field(position)
	if err == nil {
		data = fieldData.String()
	}
	//iso_msg.log.Println("len",field_data.field_def.String(),position,hex.EncodeToString(field_data.field_data));
	return data, err

}

func NewIso8583Message(specName string) *Iso8583Message {

	isoMsg := new(Iso8583Message)
	isoMsg.isoMsgDef = specMap[specName]
	isoMsg.fieldDataList = list.New()
	isoMsg.log = log.New(os.Stdout, "##iso_msg## ", log.LstdFlags)

	isoMsg.__init__()
	return isoMsg

}

//__init__ initilizes the data holding containers (list)
func (isoMsg *Iso8583Message) __init__() {

	isoMsg.nameToDataMap = make(map[string]*FieldData, 10)
	isoMsg.idToDataMap = make(map[int]*FieldData, 10)

	for l := isoMsg.isoMsgDef.fieldsDefList.Front(); l != nil; l = l.Next() {
		switch (l.Value).(type) {
		case IsoField:
			{
				var isoField IsoField = (l.Value).(IsoField)
				fdataPtr := &FieldData{fieldData: nil, fieldDef: isoField}
				isoMsg.fieldDataList.PushBack(fdataPtr)

				isoMsg.nameToDataMap[isoField.String()] = fdataPtr
				isoMsg.idToDataMap[isoField.GetId()] = fdataPtr

			}
		case BitmappedField:
			{
				var isoBmpField *BitMap = (l.Value).(*BitMap)
				isoMsg.bitMap = NewBitMap()
				for i, fDef := range isoBmpField.subFieldDef {
					if fDef != nil {
						fdataPtr := &FieldData{fieldData: nil, fieldDef: fDef}
						isoMsg.bitMap.subFieldData[i] = fdataPtr
						isoMsg.nameToDataMap[fDef.String()] = fdataPtr
						isoMsg.idToDataMap[fDef.GetId()] = fdataPtr
					}
				}
				isoMsg.fieldDataList.PushBack(isoMsg.bitMap)
				isoMsg.idToDataMap[isoBmpField.GetId()] = &FieldData{fieldData: nil, fieldDef: nil, bmpDef: isoMsg.bitMap}

			}
		default:
			{

				panic("unexpected type in iso8583 message definition!")
			}

		}
	}
}

func (isoMsg *Iso8583Message) fieldParseError(fieldName string, err error) {

	if err != nil {
		panic(fmt.Sprintf("parse_phase:error parsing field [%s] - error [%s]", fieldName, err.Error()))
	}
}

func (isoMsg *Iso8583Message) bufferUnderflowError(fieldName string) {
	panic(fmt.Sprintf("parse_phase: buffer underflow while parsing field [%s]", fieldName))
}

func (isoMsg *Iso8583Message) bufferOverflowError(data []byte) {
	isoMsg.log.Panic("parse_phase: buffer overflow -", hex.Dump(data))

}

func (isoMsg *Iso8583Message) handleError(err error) {

	if err != nil {
		panic(fmt.Sprintf("error [%s]", err.Error()))
	}
}

func (isoMsg *Iso8583Message) Field(pos int) (*FieldData, error) {

	if isoMsg.bitMap.IsOn(pos) {
		return isoMsg.bitMap.subFieldData[pos], nil
	} else {
		return &FieldData{}, errors.New("field not present")
	}

}

//set field
func (isoMsg *Iso8583Message) SetField(pos int, value string) {

	isoMsg.bitMap.SetOn(pos)
	isoMsg.bitMap.subFieldData[pos].SetData(value)

}

//set field
func (isoMsg *Iso8583Message) GetFieldByName(name string) *FieldData {

	fData := isoMsg.nameToDataMap[name]
	return fData

}

//copy all data from request to response message
func CopyRequestToResponse(isoReq *Iso8583Message, isoResp *Iso8583Message) {

	isoResp.bitMap.copyBits(isoReq.bitMap)
	for k, v := range isoReq.nameToDataMap {
		if v.fieldData != nil {
			data := make([]byte, len(v.fieldData))
			copy(data, v.fieldData)
			isoResp.nameToDataMap[k].fieldData = data
		} else {
			isoResp.nameToDataMap[k].fieldData = nil
		}
	}

}

//create a string dump of the iso message
func (isoMsg *Iso8583Message) Dump() string {

	msgBuf := bytes.NewBufferString("")
	for l := isoMsg.fieldDataList.Front(); l != nil; l = l.Next() {

		switch l.Value.(type) {
		case *FieldData:
			{

				var fData *FieldData = l.Value.(*FieldData)
				msgBuf.WriteString(fmt.Sprintf("\n%-25s: %s", fData.fieldDef.String(), fData.String()))
				break
			}

		case *BitMap:
			{

				var bmp = l.Value.(*BitMap)
				msgBuf.WriteString(fmt.Sprintf("\n%-25s: %s", "Bitmap", bmp.bitString()))

				for i, fData := range bmp.subFieldData {

					//if i == 0 || i == 1 || i == 65 || i == 129 {
					//skip invalid or bits that stand for position
					//that represents additional bitmap position
					//continue
					//}

					if fData != nil && bmp.IsOn(i) {
						msgBuf.WriteString(fmt.Sprintf("\n%-25s: %s", fData.fieldDef.String(), fData.String()))
					}
				}
				break
			}

		}
	}

	return msgBuf.String()
}

//create a string dump of the iso message
func (isoMsg *Iso8583Message) TabularFormat() *list.List {

	tabDataList := list.New()

	//msg_buf := bytes.NewBufferString("")
	for l := isoMsg.fieldDataList.Front(); l != nil; l = l.Next() {

		switch l.Value.(type) {
		case *FieldData:
			{

				var fData *FieldData = l.Value.(*FieldData)
				tabDataList.PushBack(paysim.NewTuple(fData.fieldDef.String(), fData.String()))
				//msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", f_data.field_def.String(), f_data.String()))
				break
			}

		case *BitMap:
			{

				var bmp *BitMap = l.Value.(*BitMap)
				tabDataList.PushBack(paysim.NewTuple("Bitmap", bmp.bitString()))
				//msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", "Bitmap", bmp.bit_string()))

				for i, fData := range bmp.subFieldData {

					//if i == 0 || i == 1 || i == 65 || i == 129 {
					//skip invalid or bits that stand for position
					//that represents additional bitmap position
					//continue
					//}

					if fData != nil && bmp.IsOn(i) {
						tabDataList.PushBack(paysim.NewTuple(fData.fieldDef.String(), fData.String()))
						//msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", f_data.field_def.String(), f_data.String()))
					}
				}
				break
			}

		}
	}

	return tabDataList

}

//parse the bytes from 'buf' and populate 'Iso8583Message'
func (isoMsg *Iso8583Message) Parse(buf *bytes.Buffer) (err error) {

	defer func() {
		str := recover()
		if str != nil {
			isoMsg.log.Printf("parse error. message: %s", str)
			err = errors.New("parse error")
		}
	}()

	for l := isoMsg.fieldDataList.Front(); l != nil; l = l.Next() {

		switch l.Value.(type) {
		case *FieldData:
			{

				var fData *FieldData = l.Value.(*FieldData)
				pylog.Log("parsing.. ", fData.fieldDef.Def())
				fData.fieldDef.Parse(isoMsg, fData, buf)
				break
			}

		case *BitMap:
			{

				var bmp *BitMap = l.Value.(*BitMap)
				bmp.Parse(isoMsg, buf)
				//parse sub fields of bitmap
				for i, fData := range bmp.subFieldData {

					//if i == 0 || i == 1 || i == 65 || i == 129 {
					//skip invalid or bits that stand for position
					//that represents additional bitmap position
					//continue
					//}

					if fData != nil && bmp.IsOn(i) {
						pylog.Log("parsing.. ", fData.fieldDef.Def())
						fData.fieldDef.Parse(isoMsg, fData, buf)
					}
				}
				break
			}

		}
	}

	if buf.Len() > 0 {
		isoMsg.bufferOverflowError(buf.Bytes())
	}

	return err

}

func (isoMsg *Iso8583Message) Bytes() []byte {

	msgBuf := bytes.NewBuffer(make([]byte, 0))

	for l := isoMsg.fieldDataList.Front(); l != nil; l = l.Next() {

		switch obj := l.Value.(type) {

		case *FieldData:
			{
				msgBuf.Write(obj.Bytes())
				break
			}
		case BitmappedField:
			{
				msgBuf.Write(isoMsg.bitMap.Bytes())
				bmp := obj.(*BitMap)

				for i, v := range bmp.subFieldData {
					if v != nil && v.fieldData != nil &&
						bmp.IsOn(i) {

						fData := v.Bytes()
						isoMsg.log.Printf("assembling: %s - len: %d data: %s final data: %s\n",
							v.fieldDef.String(), len(v.fieldData), hex.EncodeToString(v.fieldData),
							hex.EncodeToString(fData))
						msgBuf.Write(fData)
					}
				}
			}
		}

	}

	return msgBuf.Bytes()
}

func (isoMsg *Iso8583Message) SetFieldData(id int, fieldVal string) {

	isoMsg.idToDataMap[id].SetData(fieldVal)
}

func (isoMsg *Iso8583Message) GetFieldDataById(id int) *FieldData {

	return isoMsg.idToDataMap[id]

}
