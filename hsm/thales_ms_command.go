package hsm

import (
	"bytes"
	"encoding/hex"
	"github.com/rkbalgi/go/crypto"
	//_ "github.com/rkbalgi/hsm"
)

type ThalesMsRequest struct {
	_prologue          prologue
	message_block_number uint
	key_type            uint
	key_length          uint
	message_type        uint
	key                string
	iv                 string
	message_length      uint
	message_block       []byte
	_epilogue          epilogue
}

type ThalesMsResponse struct {
	message_header       []byte
	response_code        []byte
	error_code           []byte
	MAB                 []byte
	end_message_delimiter string
	message_trailer      string
}

func (hsm_handle *ThalesHsm) Handle_MS(msg_data []byte) []byte {

	defer func() {
		str := recover()
		if str != nil {
			hsm_handle.log.Println("unexpected system error", str)
		}
	}()

	ms_req_struct := new(ThalesMsRequest)
	ms_resp_struct := new(ThalesMsResponse)

	msg_buf := bytes.NewBuffer(msg_data)
	//at every step this boolean will be updated
	//if false, respond with a 15 error code
	parse_ok := true

	parse_prologue(msg_buf, &ms_req_struct._prologue, hsm_handle.header_length)

	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.message_block_number, 1, DecimalInt)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.key_type, 1, DecimalInt)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.key_length, 1, DecimalInt)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	//message type
	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.message_type, 1, DecimalInt)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	//key
	parse_ok = read_key(msg_buf, &ms_req_struct.key)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	//iv
	if ms_req_struct.message_block_number == 2 || ms_req_struct.message_block_number == 3 {
		parse_ok = read_fixed_field(msg_buf, &ms_req_struct.iv, 16, String)
		if !parse_ok {
			return (ms_req_struct.invalid_data_response(ms_resp_struct))
		}

	}

	//message length

	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.message_length, 4, HexadecimalInt)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	//message block

	parse_ok = read_fixed_field(msg_buf, &ms_req_struct.message_block, ms_req_struct.message_length, Binary)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	parse_ok = parse_epilogue(msg_buf, &ms_req_struct._epilogue)
	if !parse_ok {
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	if __hsm_debug_enabled {
		hsm_handle.log.Printf(Dump(*ms_req_struct))
	}

	//decrypt the key under the appropriate LMK
	key_type := "000"
	if ms_req_struct.key_type == 0 {
		//TAK
		key_type = "003"
	} else if ms_req_struct.key_type == 1 {
		//ZAK
		key_type = "008"
	} else {
		//error
		hsm_handle.log.Printf("invalid key type - %s\n", key_type)
		return (ms_req_struct.invalid_data_response(ms_resp_struct))
	}

	mac_key := decrypt_key(ms_req_struct.key, key_type)

	
	var mac []byte
	if len(mac_key) == 16 {
		//generate X9.19 mac
		mac = crypto.GenerateMac_X919(ms_req_struct.message_block, mac_key)
	}else {
		//x9.9
		mac = crypto.GenerateMac_X99(ms_req_struct.message_block, mac_key)
		
	}
	ms_resp_struct.MAB = mac
	ms_resp_struct.response_code = []byte("00")

	if __hsm_debug_enabled {
		hsm_handle.log.Printf("MAC: %s", hex.EncodeToString(mac))
	}

	//send response
	return ms_resp_struct.generate_response(ms_req_struct)

}

func (ms_req_struct *ThalesMsRequest) invalid_data_response(ms_resp_struct *ThalesMsResponse) []byte {
	set_fixed_field(&ms_resp_struct.error_code, 2, uint(15), DecimalInt)
	return ms_resp_struct.generate_response(ms_req_struct)

}

func (ms_resp_struct *ThalesMsResponse) generate_response(ms_req_struct *ThalesMsRequest) []byte {

	ms_resp_buf := bytes.NewBuffer([]byte(ms_req_struct._prologue.header))
	resp_cmd_code := []byte(ms_req_struct._prologue.command_name)
	resp_cmd_code[1] = resp_cmd_code[1] + 1
	ms_resp_buf.Write(resp_cmd_code)

	ms_resp_buf.Write(ms_resp_struct.response_code)
	if ms_resp_struct.MAB != nil {
		ms_resp_buf.Write(to_ascii(ms_resp_struct.MAB))
	}
	return ms_resp_buf.Bytes()

}
