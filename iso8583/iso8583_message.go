package iso8583

import (
	"bytes"
	_ "encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	_ "strconv"
)

type Iso8583Message struct {
	iso_msg_def     *Iso8583MessageDef
	msg_type        string
	bit_map         *BitMap
	field_data_list []FieldData
	log             *log.Logger
}

//GetMessageType returns the 'Message Type' as string
func (iso_msg *Iso8583Message) GetMessageType() string {
	return iso_msg.msg_type
}


//GetMessageType returns the 'Message Type' as string
func (iso_msg *Iso8583Message) GetBinaryBitmap() string {
	
	
	binary_bmp_str:=bytes.NewBufferString("");
	for i:=1;i<129;i++{
		if iso_msg.bit_map.IsOn(i){
			binary_bmp_str.WriteString("1");
		}else{
			binary_bmp_str.WriteString("0");
		}
	}
	
	return binary_bmp_str.String();
	
	
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

func NewIso8583Message() *Iso8583Message {

	iso_msg := new(Iso8583Message)
	iso_msg.iso_msg_def = iso8583_msg_def
	iso_msg.bit_map = NewBitMap()
	iso_msg.field_data_list = make([]FieldData, len(iso8583_msg_def.fields))
	for i, f_def := range iso8583_msg_def.fields {
		if f_def != nil {
			//fmt.Println(i,f_def.String())
			iso_msg.field_data_list[i].field_def = f_def
		}
	}
	iso_msg.log = log.New(os.Stdout, "##iso_msg## ", log.LstdFlags)
	return iso_msg

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

func (iso_msg *Iso8583Message) get_field(pos int) (FieldData, error) {

	if iso_msg.bit_map.IsOn(pos) {
		log.Println("get_field ", hex.EncodeToString(iso_msg.field_data_list[pos].field_data))
		return iso_msg.field_data_list[pos], nil
	} else {
		return FieldData{}, errors.New("field not present")
	}

}

//set field
func (iso_msg *Iso8583Message) set_field(pos int, value string) {

	iso_msg.bit_map.SetOn(pos)
	iso_msg.field_data_list[pos].SetData(value)

}

//copy all data from req to response message
func copy_iso_req_to_resp(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	iso_resp.bit_map = iso_req.bit_map.Copy()
	//iso_resp.field_data_list = make([]FieldData, len(iso_req.field_data_list))

	for i := 1; i < 129+64; i++ {

		if iso_req.bit_map.IsOn(i) {
			iso_resp.field_data_list[i] = *iso_req.field_data_list[i].copy()
		}

	}

}

//this method handles an incoming ISO8583 message, doing the parsing, processing
//and response creation
func Handle(buf *bytes.Buffer) (resp_iso_msg *Iso8583Message, err error) {

	req_iso_msg := NewIso8583Message()

	//parse incoming message
	err = req_iso_msg.Parse(buf)
	if err != nil {
		return nil, err
	}

	req_iso_msg.log.Println("parsed incoming message: ", req_iso_msg.Dump())

	//continue handling

	resp_iso_msg = NewIso8583Message()
	switch req_iso_msg.msg_type {
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
			err = errors.New("unsupported message type -" + req_iso_msg.msg_type)

		}
	}

	req_iso_msg.log.Println("outgoing message: ", resp_iso_msg.Dump())

	return resp_iso_msg, err

}

//create a string dump of the iso message
func (iso_msg *Iso8583Message) Dump() string {

	msg_buf := bytes.NewBufferString(fmt.Sprintf("\n%-25s: %s\n", "Message Type", iso_msg.msg_type))
	msg_buf.WriteString(fmt.Sprintf("%-25s: %s\n", "BitMap", hex.EncodeToString(iso_msg.bit_map.Bytes())))
	for i, v := range iso_msg.field_data_list {
		if v.field_def != nil && iso_msg.bit_map.IsOn(i) {
			msg_buf.WriteString(fmt.Sprintf("%-25s: %s\n", v.field_def.String(), v.String()))
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

	v_data := make([]byte, 4)
	n, err := buf.Read(v_data)
	if n != 4 || err != nil {
		if n != 4 {
			iso_msg.buffer_underflow_error("Message Type")
		} else {
			iso_msg.field_parse_error("Message Type", err)
		}
	}
	iso_msg.msg_type = string(v_data)
	iso_msg.bit_map.Parse(iso_msg, buf)

	for i, fld_def := range iso_msg.iso_msg_def.fields {

		if i == 0 || i == 1 || i == 65 || i == 129 {
			//skip invalid or bits that stand for position
			//that represents additional bitmap position
			continue
		}

		if iso_msg.bit_map.IsOn(i) {
			//fmt.Println("parsing position",i);
			if fld_def != nil {
				//fmt.Println("parsing position",fld_def.String());
				iso_msg.field_data_list[i] = *fld_def.Parse(iso_msg, buf)
			} else {
				//not a defined field
				panic(fmt.Sprintf("no definition for bit position - %d\n", i))
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
	msg_buf.Write([]byte(iso_msg.msg_type))
	msg_buf.Write(iso_msg.bit_map.Bytes())

	for i, v := range iso_msg.field_data_list {
		if v.field_def != nil && iso_msg.bit_map.IsOn(i) {
			f_data := v.Bytes()
			iso_msg.log.Printf("assembling: %s - len: %d data: %s final data: %s\n",
				v.field_def.String(), len(v.field_data), hex.EncodeToString(v.field_data),
				hex.EncodeToString(f_data))
			msg_buf.Write(f_data)
		}
	}

	return msg_buf.Bytes()
}
