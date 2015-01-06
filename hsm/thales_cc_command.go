package hsm

import (
	"bytes"
	"fmt"
	"encoding/hex"
)

type thales_cc_req struct {
	_pro                   prologue
	source_zpk             string
	destination_zpk        string
	maximum_pin_len        uint
	source_pin_block       string //change type to []byte
	source_pinblk_fmt      uint
	destination_pinblk_fmt uint
	acct_no                string

	_epi epilogue
}

type thales_cc_resp struct {
	_pro prologue

	response_code          string
	error_code             string
	pin_len                []byte
	destination_pin_block  []byte
	destination_pinblk_fmt []byte
	end_message_delimiter  byte
	message_trailer        []byte
}

func (hsm_handle *ThalesHsm) handle_cc_command(msg_data []byte) []byte {

	msg_buf := bytes.NewBuffer(msg_data)

	req := new(thales_cc_req)
	resp := new(thales_cc_resp)

	if parse_prologue(msg_buf, &req._pro, hsm_handle.header_length) {
		parse_ok := read_key(msg_buf, &req.source_zpk)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = read_key(msg_buf, &req.destination_zpk)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = read_fixed_field(msg_buf, &req.maximum_pin_len, 2, DecimalInt)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = read_fixed_field(msg_buf, &req.source_pin_block, 16, String)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = read_fixed_field(msg_buf, &req.source_pinblk_fmt, 2, DecimalInt)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = read_fixed_field(msg_buf, &req.destination_pinblk_fmt, 2, DecimalInt)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		//TODO:should be 18H for formats 04
		parse_ok = read_fixed_field(msg_buf, &req.acct_no, 12, String)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = parse_epilogue(msg_buf, &req._epi)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
		parse_ok = parse_epilogue(msg_buf, &req._epi)
		if !parse_ok {
			return (req.invalid_data_response(resp))
		}
	} else {
		//no prolog, message should be dropped
		hsm_handle.log.Println("[CC] prolog could not be parsed, dropping message")
		return (nil)
	}

	in_key := decrypt_key(req.source_zpk, "001")
	out_key := decrypt_key(req.destination_zpk, "001")
	
	if(__hsm_debug_enabled){
		hsm_handle.log.Println(Dump(*req))
	}

	//handle message i.e translate pin block
	ph := new_pin_handler(req.acct_no, int(req.source_pinblk_fmt), int(req.destination_pinblk_fmt), in_key, out_key)
	in_pin_block, err := hex.DecodeString(req.source_pin_block)
	if err != nil {
		hsm_handle.log.Printf("invalid source pin block - %s", req.source_pin_block)
		return req.invalid_data_response(resp)
	}
	out_pin_blk := ph.translate(in_pin_block)
	if out_pin_blk != nil && len(out_pin_blk) == 8 {
		resp.error_code = HSM_OK
		if(__hsm_debug_enabled){
			hsm_handle.log.Println("clear pin -",ph.get_clear_pin())
			hsm_handle.log.Println("destination pin block - ",hex.EncodeToString(out_pin_blk));
		}
		resp.destination_pin_block = []byte(hex.EncodeToString(out_pin_blk))
		resp.pin_len = []byte(fmt.Sprintf("%02d", len(ph.get_clear_pin())))
	} else {
		return req.invalid_data_response(resp)
	}
	
	
	//generate response
	return req.generate_response(resp);

}

func (req *thales_cc_req) invalid_data_response(resp *thales_cc_resp) []byte {

	resp.error_code = HSM_PARSE_ERROR
	return (req.generate_response(resp))

}

func (req *thales_cc_req) generate_response(resp *thales_cc_resp) []byte {

	resp_buf := bytes.NewBuffer([]byte(req._pro.header))
	resp_cmd_code := []byte(req._pro.command_name)
	resp_cmd_code[1] = resp_cmd_code[1] + 1
	resp_buf.Write(resp_cmd_code)

	resp_buf.WriteString(resp.error_code)
	if resp.error_code == HSM_OK {
		resp_buf.Write(resp.pin_len)
		resp_buf.Write(resp.destination_pin_block)
		resp_buf.WriteString(fmt.Sprintf("%02d", req.destination_pinblk_fmt))

	}
	if req._epi.end_message_delimiter == 0x19 {
		resp_buf.WriteByte(req._epi.end_message_delimiter)
		resp_buf.Write(req._epi.message_trailer)
	}

	return resp_buf.Bytes()

}
