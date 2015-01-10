package hsm

import (
	"bytes"
	"encoding/hex"
	"github.com/rkbalgi/go/crypto"
	"github.com/rkbalgi/go/hsm/keys"
)

type thales_a0_req struct {
	_pro                  prologue
	mode                  uint //hex1
	key_type              string
	key_scheme_lmk        string
	derive_key_mode       string
	dukpt_master_key_type uint //hex1
	dukpt_master_key      string
	ksn                   string //15
	delimiter             string
	zmk_tmk_flag          uint //dec1
	zmk_tmk               string
	key_scheme_zmk        string //1
	attalla_variant       string

	_epi epilogue
}

type thales_a0_resp struct {
	_pro prologue

	response_code         string
	error_code            string
	key_under_lmk         []byte
	key_under_zmk         []byte
	key_check_value       []byte
	end_message_delimiter byte
	message_trailer       []byte
}

func (hsm_handle *ThalesHsm) handle_a0_command(msg_data []byte) []byte {

	zmk_tmk_present := false
	msg_buf := bytes.NewBuffer(msg_data)

	req := new(thales_a0_req)
	resp := new(thales_a0_resp)

	if parse_prologue(msg_buf, &req._pro, hsm_handle.header_length) {

		parse_ok := read_fixed_field(msg_buf, &req.mode, 1, HexadecimalInt)
		parse_ok = read_fixed_field(msg_buf, &req.key_type, 3, String)
		parse_ok = read_fixed_field(msg_buf, &req.key_scheme_lmk, 1, String)

		if req.mode == 0xa || req.mode == 0xb {
			parse_ok = read_fixed_field(msg_buf, &req.derive_key_mode, 1, String)
			if !parse_ok {
				return (req.invalid_data_response(resp))
			} else {
				//derive key mode should be 0
				if req.derive_key_mode != "0" {
					hsm_handle.log.Printf("invalid derive key mode - ", req.derive_key_mode)
					return (req.invalid_data_response(resp))
				} else {
					//read dupkt master key type and key
					parse_ok = read_fixed_field(msg_buf, &req.dukpt_master_key_type, 1, HexadecimalInt)
					if !parse_ok {
						return (req.invalid_data_response(resp))
					} else {
						if req.dukpt_master_key_type == 0x01 || req.dukpt_master_key_type == 0x02 {
							parse_ok = read_key(msg_buf, &req.dukpt_master_key)
							if !parse_ok {
								return (req.invalid_data_response(resp))
							} else {
								//read ksn
								parse_ok = read_fixed_field(msg_buf, &req.ksn, 15, String)
								if !parse_ok {
									return (req.invalid_data_response(resp))
								}
								//check if KSN is all hex, else throw error
								if !hex_regexp.MatchString(req.ksn) {
									hsm_handle.log.Printf("invalid ksn - ", req.ksn)
									return (req.invalid_data_response(resp))
								}

							}
						} else {
							hsm_handle.log.Printf("invalid dukpt master key type - ", req.dukpt_master_key_type)
							return (req.invalid_data_response(resp))
						}
					}
				}
			}
		}

		if msg_buf.Len() > 0 {
			if req.mode == 0x01 || req.mode == 0x0b {

				parse_ok = read_fixed_field(msg_buf, &req.delimiter, 1, String)
				if !parse_ok {
					return (req.invalid_data_response(resp))
				}
				if req.delimiter != ";" {
					hsm_handle.log.Printf("invalid delimiter - ", req.delimiter)
					return (req.invalid_data_response(resp))
				}
				parse_ok = read_fixed_field(msg_buf, &req.zmk_tmk_flag, 1, DecimalInt)
				if !parse_ok {
					return (req.invalid_data_response(resp))
				}
				if req.zmk_tmk_flag == 0 || req.zmk_tmk_flag == 1 {
					//ZMK or TMK
					parse_ok = read_key(msg_buf, &req.zmk_tmk)

					if !parse_ok {
						return (req.invalid_data_response(resp))
					}
					zmk_tmk_present = true
					parse_ok = read_fixed_field(msg_buf, &req.key_scheme_zmk, 1, String)
					if !parse_ok {
						return (req.invalid_data_response(resp))
					}
					//there may be attalla variant optionally.
					var b, b2 byte
					if msg_buf.Len() > 0 {
						b, _ = msg_buf.ReadByte()
						if b == byte('%') {
							msg_buf.UnreadByte()
						} else {
							//there is attala variant
							if msg_buf.Len() > 0 {
								b2, _ = msg_buf.ReadByte()
								if b2 == byte('%') {
									msg_buf.UnreadByte()
									//just a single digit variant
									req.attalla_variant = string([]byte{b})
								} else {
									req.attalla_variant = string([]byte{b, b2})
								}
							} else {
								//eob - single digit atalla variant
								req.attalla_variant = string([]byte{b})
							}
						}
					}

				} else {
					hsm_handle.log.Printf("invalid zmk/tmk flag - ", req.zmk_tmk_flag)
					return (req.invalid_data_response(resp))
				}

			}
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

	if __hsm_debug_enabled {
		hsm_handle.log.Println(Dump(*req))
	}

	if req.mode == 0x0a || req.mode == 0x0b {
		//we do not support it at the moment
		hsm_handle.log.Printf("derive mode (A, B) is not supported at the moment!")
		return (req.invalid_data_response(resp))
	}

	var zmk_tmk_key []byte

	if zmk_tmk_present {
		if req.zmk_tmk_flag == 0 {
			//zmk
			zmk_tmk_key = decrypt_key(req.zmk_tmk, ZMK_KEY_TYPE)
		} else {
			zmk_tmk_key = decrypt_key(req.zmk_tmk, TMK_KEY_TYPE)
		}
	}

	key_len := 0
	switch req.key_scheme_lmk {
	case keys.Z:
		key_len = 8
	case keys.U:
		key_len = 16
	case keys.T:
		key_len = 24
	default:
		{
			hsm_handle.log.Printf("invalid lmk key scheme ", req.key_scheme_lmk)
			return (req.invalid_data_response(resp))
		}
	}

	//generate the required key and its check value
	key := crypto.GenerateDesKey(key_len)
	//generate check value
	resp.key_check_value = gen_check_value(key)[:3]
	
	
	
	if(__hsm_debug_enabled){
		hsm_handle.log.Println("key value: ",hex.EncodeToString(key),"check value: ",hex.EncodeToString(resp.key_check_value));
	}

	//TODO:: odd parity enforcement
	if req.key_scheme_lmk == keys.Z {
		resp.key_under_lmk = encrypt_key(hex.EncodeToString(key), req.key_type)
	} else {
		resp.key_under_lmk = encrypt_key(req.key_scheme_lmk+hex.EncodeToString(key), req.key_type)
	}

	if zmk_tmk_present {
		//key should also be encrypted under ZMK/TMK
		switch {
		case req.key_scheme_zmk == keys.Z || req.key_scheme_zmk == keys.U || req.key_scheme_zmk == keys.T:
			{
				resp.key_under_zmk = encrypt_key_kek(req.key_scheme_zmk+hex.EncodeToString(key), zmk_tmk_key, req.key_type)
			}
		case req.key_scheme_zmk == keys.X || req.key_scheme_zmk == keys.Y:
			{
				resp.key_under_zmk = encrypt_key_kek_x917(hex.EncodeToString(key), zmk_tmk_key, req.key_type)
				hsm_handle.log.Println(hex.EncodeToString(resp.key_under_zmk),"??",hex.EncodeToString(key),"???",hex.EncodeToString(zmk_tmk_key))
			}

		default:
			{
				hsm_handle.log.Printf("invalid zmk key scheme ", req.key_scheme_lmk)
				return (req.invalid_data_response(resp))
			}
		}
	}

	//keys should be ascii encoded
	if req.key_scheme_lmk == keys.Z {
		//single length keys do not require
		//a scheme identifier
		resp.key_under_lmk = []byte(hex.EncodeToString(resp.key_under_lmk))
		if zmk_tmk_present {
			resp.key_under_zmk = []byte(hex.EncodeToString(resp.key_under_zmk))
		}
	} else {
		resp.key_under_lmk = []byte(req.key_scheme_lmk + hex.EncodeToString(resp.key_under_lmk))
		if zmk_tmk_present {
			resp.key_under_zmk = []byte(req.key_scheme_zmk + hex.EncodeToString(resp.key_under_zmk))
		}
	}
	resp.key_check_value = []byte(hex.EncodeToString(resp.key_check_value))

	resp.error_code = HSM_OK

	//generate response
	return req.generate_response(resp)

}

func (req *thales_a0_req) invalid_data_response(resp *thales_a0_resp) []byte {

	resp.error_code = HSM_PARSE_ERROR
	return (req.generate_response(resp))

}

func (req *thales_a0_req) generate_response(resp *thales_a0_resp) []byte {

	resp_buf := bytes.NewBuffer([]byte(req._pro.header))
	resp_cmd_code := []byte(req._pro.command_name)
	resp_cmd_code[1] = resp_cmd_code[1] + 1
	resp_buf.Write(resp_cmd_code)

	resp_buf.WriteString(resp.error_code)
	if resp.error_code == HSM_OK {

		resp_buf.Write(resp.key_under_lmk)
		if req.mode == 0x01 || req.mode == 0x0b {
			resp_buf.Write(resp.key_under_zmk)
		}
		resp_buf.Write(resp.key_check_value)
	}

	if req._epi.end_message_delimiter == 0x19 {
		resp_buf.WriteByte(req._epi.end_message_delimiter)
		resp_buf.Write(req._epi.message_trailer)
	}

	return resp_buf.Bytes()

}
