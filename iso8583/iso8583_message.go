package iso8583

import (
	"bytes"
	"container/list"
	_ "encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	_ "reflect"
	_ "strconv"
)

type Iso8583Message struct {
	iso_msg_def      *Iso8583MessageDef
	field_data_list  *list.List
	log              *log.Logger
	bit_map          *BitMap //for convenience
	name_to_data_map map[string]*FieldData
}

//GetMessageType returns the 'Message Type' as string
func (iso_msg *Iso8583Message) GetMessageType() string {
	return iso_msg.name_to_data_map["Message Type"].String()
}

func (iso_msg *Iso8583Message) ToWebMsg(is_req bool) *WebMsgData {

	json_msg := WebMsgData{}
	json_msg.Spec = iso_msg.iso_msg_def.spec_name
	if is_req {
		json_msg.Type = "Request"
	} else {
		json_msg.Type = "Response"
	}
	json_msg.DataArray = make([]string, iso_msg.iso_msg_def.field_seq)

	for l := iso_msg.iso_msg_def.fields_def_list.Front(); l != nil; l = l.Next() {
		switch obj := l.Value.(type) {
		case IsoField:
			{
				iso_field := iso_msg.get_field_by_name(obj.String())
				if iso_field.field_data != nil {
					json_msg.DataArray[iso_field.field_def.GetId()] = iso_field.String()
				}
			}
		case BitmappedField:
			{

				json_msg.DataArray[obj.GetId()] = iso_msg.bit_map.bit_string()
				for f_pos, f_data := range iso_msg.bit_map.sub_field_data {
					if f_data != nil && f_data.field_data != nil && iso_msg.bit_map.IsOn(f_pos) {
						json_msg.DataArray[f_data.field_def.GetId()] = f_data.String()
					}
				}

			} //end case

		} //end switch
	} //end for

	return &json_msg

}

//SetData sets data into individual fields by id
func (iso_msg *Iso8583Message) SetData(data []string) {

	for l := iso_msg.iso_msg_def.fields_def_list.Front(); l != nil; l = l.Next() {

		switch obj := l.Value.(type) {
		case IsoField:
			{
				iso_field := iso_msg.get_field_by_name(obj.String())
				iso_field.SetData(data[iso_field.field_def.GetId()])

			}
		case BitmappedField:
			{
				bitmap_val := data[obj.GetId()]
				for i := 0; i < len(bitmap_val); i++ {

					if bitmap_val[i:i+1] == "1" {
						iso_msg.bit_map.SetOn(i + 1)
					} else {
						iso_msg.bit_map.SetOff(i + 1)
					}
				}

				for f_pos, f_data := range iso_msg.bit_map.sub_field_data {
					if f_data != nil && iso_msg.bit_map.IsOn(f_pos) {
						f_data.SetData(data[f_data.field_def.GetId()])
						//iso_msg.bit_map.SetOn(f_pos)
					}
				}

			}

		}
	}
}

//GetBinaryBitmap returns the 'Bitmap' as binary string
func (iso_msg *Iso8583Message) GetBinaryBitmap() string {

	binary_bmp_str := bytes.NewBufferString("")
	for i := 1; i < 129; i++ {
		if iso_msg.bit_map.IsOn(i) {
			binary_bmp_str.WriteString("1")
		} else {
			binary_bmp_str.WriteString("0")
		}
	}

	return binary_bmp_str.String()

}

//IsSelected returns a boolean indicating
//if the 'position' is selected in the bitmap
func (iso_msg *Iso8583Message) IsSelected(position int) bool {
	return iso_msg.bit_map.IsOn(position)
}

//GetFieldData returns the data associated with the 'position'
//in the iso_msg
func (iso_msg *Iso8583Message) GetFieldData(position int) (data string, err error) {
	field_data, err := iso_msg.get_field(position)
	if err == nil {
		data = field_data.String()
	}
	//iso_msg.log.Println("len",field_data.field_def.String(),position,hex.EncodeToString(field_data.field_data));
	return data, err

}

func NewIso8583Message(spec_name string) *Iso8583Message {

	iso_msg := new(Iso8583Message)
	iso_msg.iso_msg_def = spec_map[spec_name]
	iso_msg.field_data_list = list.New()
	iso_msg.log = log.New(os.Stdout, "##iso_msg## ", log.LstdFlags)

	iso_msg.__init__()
	return iso_msg

}

//__init__ initilizes the data holding containers (list)
func (iso_msg *Iso8583Message) __init__() {

	iso_msg.name_to_data_map = make(map[string]*FieldData, 10)

	for l := iso_msg.iso_msg_def.fields_def_list.Front(); l != nil; l = l.Next() {
		switch (l.Value).(type) {
		case IsoField:
			{
				var iso_field IsoField = (l.Value).(IsoField)
				fdata_ptr := &FieldData{field_data: nil, field_def: iso_field}
				iso_msg.field_data_list.PushBack(fdata_ptr)
				iso_msg.name_to_data_map[iso_field.String()] = fdata_ptr

			}
		case BitmappedField:
			{
				var iso_bmp_field *BitMap = (l.Value).(*BitMap)
				iso_msg.bit_map = NewBitMap()
				for i, f_def := range iso_bmp_field.sub_field_def {
					if f_def != nil {
						fdata_ptr := &FieldData{field_data: nil, field_def: f_def}
						iso_msg.name_to_data_map[f_def.String()] = fdata_ptr
						iso_msg.bit_map.sub_field_data[i] = fdata_ptr
					}
				}
				iso_msg.field_data_list.PushBack(iso_msg.bit_map)

			}
		default:
			{

				panic("unexpected type in iso8583 message definition!")
			}

		}
	}
}

func (iso_msg *Iso8583Message) field_parse_error(field_name string, err error) {

	if err != nil {
		panic(fmt.Sprintf("parse_phase:error parsing field [%s] - error [%s]", field_name, err.Error()))
	}
}

func (iso_msg *Iso8583Message) buffer_underflow_error(field_name string) {
	panic(fmt.Sprintf("parse_phase: buffer underflow while parsing field [%s]", field_name))
}

func (iso_msg *Iso8583Message) buffer_overflow_error(data []byte) {
	iso_msg.log.Panic("parse_phase: buffer overflow -", hex.Dump(data))

}

func (iso_msg *Iso8583Message) handle_error(err error) {

	if err != nil {
		panic(fmt.Sprintf("error [%s]", err.Error()))
	}
}

func (iso_msg *Iso8583Message) get_field(pos int) (*FieldData, error) {

	if iso_msg.bit_map.IsOn(pos) {
		return iso_msg.bit_map.sub_field_data[pos], nil
	} else {
		return &FieldData{}, errors.New("field not present")
	}

}

//set field
func (iso_msg *Iso8583Message) set_field(pos int, value string) {

	iso_msg.bit_map.SetOn(pos)
	iso_msg.bit_map.sub_field_data[pos].SetData(value)

}

//set field
func (iso_msg *Iso8583Message) get_field_by_name(name string) *FieldData {

	f_data := iso_msg.name_to_data_map[name]
	return f_data

}

//copy all data from request to response message
func copy_iso_req_to_resp(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	iso_resp.bit_map.copy_bits(iso_req.bit_map)
	for k, v := range iso_req.name_to_data_map {
		if v.field_data != nil {
			data := make([]byte, len(v.field_data))
			copy(data, v.field_data)
			iso_resp.name_to_data_map[k].field_data = data
		} else {
			iso_resp.name_to_data_map[k].field_data = nil
		}
	}


}

//this method handles an incoming ISO8583 message, doing the parsing, processing
//and response creation
func Handle(spec_name string,buf *bytes.Buffer) (resp_iso_msg *Iso8583Message, err error) {

	req_iso_msg := NewIso8583Message(spec_name)

	//parse incoming message
	err = req_iso_msg.Parse(buf)
	if err != nil {
		return nil, err
	}

	req_iso_msg.log.Println("parsed incoming message: ", req_iso_msg.Dump())

	//continue handling

	resp_iso_msg = NewIso8583Message(spec_name)
	msg_type := req_iso_msg.GetMessageType()
	switch msg_type {
	case ISO_MSG_1100:
		{
			handle_auth_req(req_iso_msg, resp_iso_msg)
		}
	case ISO_MSG_1804:
		{
			handle_network_req(req_iso_msg, resp_iso_msg)
		}
	case ISO_MSG_1420:
		{
			handle_reversal_req(req_iso_msg, resp_iso_msg)
		}
	default:
		{
			err = errors.New("unsupported message type -" + req_iso_msg.GetMessageType())

		}
	}

	req_iso_msg.log.Println("outgoing message: ", resp_iso_msg.Dump())

	return resp_iso_msg, err

}

//create a string dump of the iso message
func (iso_msg *Iso8583Message) Dump() string {

	msg_buf := bytes.NewBufferString("")
	for l := iso_msg.field_data_list.Front(); l != nil; l = l.Next() {

		switch l.Value.(type) {
		case *FieldData:
			{

				var f_data *FieldData = l.Value.(*FieldData)
				msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", f_data.field_def.String(), f_data.String()))
				break
			}

		case *BitMap:
			{

				var bmp *BitMap = l.Value.(*BitMap)
				msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", "Bitmap", bmp.bit_string()))

				for i, f_data := range bmp.sub_field_data {

					//if i == 0 || i == 1 || i == 65 || i == 129 {
					//skip invalid or bits that stand for position
					//that represents additional bitmap position
					//continue
					//}

					if f_data != nil && bmp.IsOn(i) {
						msg_buf.WriteString(fmt.Sprintf("\n%-25s: %s", f_data.field_def.String(), f_data.String()))
					}
				}
				break
			}

		}
	}

	return msg_buf.String()
}

//parse the bytes from 'buf' and populate 'Iso8583Message'
func (iso_msg *Iso8583Message) Parse(buf *bytes.Buffer) (err error) {

	defer func() {
		str := recover()
		if str != nil {
			iso_msg.log.Printf("parse error. message: %s", str)
			err = errors.New("parse error")
		}
	}()

	for l := iso_msg.field_data_list.Front(); l != nil; l = l.Next() {

		switch l.Value.(type) {
		case *FieldData:
			{

				var f_data *FieldData = l.Value.(*FieldData)
				f_data.field_def.Parse(iso_msg, f_data, buf)
				break
			}

		case *BitMap:
			{

				var bmp *BitMap = l.Value.(*BitMap)
				bmp.Parse(iso_msg, buf)
				//parse sub fields of bitmap
				for i, f_data := range bmp.sub_field_data {

					//if i == 0 || i == 1 || i == 65 || i == 129 {
					//skip invalid or bits that stand for position
					//that represents additional bitmap position
					//continue
					//}

					if f_data != nil && bmp.IsOn(i) {
						f_data.field_def.Parse(iso_msg, f_data, buf)
					}
				}
				break
			}

		}
	}

	if buf.Len() > 0 {
		iso_msg.buffer_overflow_error(buf.Bytes())
	}

	return err

}

func (iso_msg *Iso8583Message) Bytes() []byte {

	msg_buf := bytes.NewBuffer(make([]byte, 0))

	for l := iso_msg.field_data_list.Front(); l != nil; l = l.Next() {

		switch obj := l.Value.(type) {

		case *FieldData:
			{
				msg_buf.Write(obj.Bytes())
				break
			}
		case BitmappedField:
			{
				msg_buf.Write(iso_msg.bit_map.Bytes())
				bmp := obj.(*BitMap)

				for i, v := range bmp.sub_field_data {
					if v != nil && v.field_data != nil &&
						bmp.IsOn(i) {

						f_data := v.Bytes()
						iso_msg.log.Printf("assembling: %s - len: %d data: %s final data: %s\n",
							v.field_def.String(), len(v.field_data), hex.EncodeToString(v.field_data),
							hex.EncodeToString(f_data))
						msg_buf.Write(f_data)
					}
				}
			}
		}

	}

	return msg_buf.Bytes()
}
