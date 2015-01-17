package iso8583

//handles 1100 message
func handle_auth_req(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	//copy all data from req to response
	//and then lets selective remove fields that are
	//not required in the response
	copy_iso_req_to_resp(iso_req, iso_resp)

	//set message type on response as 1110
	iso_resp.msg_type = ISO_MSG_1110

	//for demo purposes, we will simply base our responses
	//on the input amounts

	f_amount, err := iso_req.get_field(4)
	if err != nil {
		do_format_error_response(iso_resp)
		return
	}
	l_amount := str_to_uint64(f_amount.String())

	switch {

	case l_amount > 800 && l_amount < 900:
		{
			iso_resp.set_field(39, ISO_RESP_DECLINE)
		}
	default:
		{
			iso_resp.set_field(39, "APPISO")
			iso_resp.set_field(38, ISO_RESP_APPROVAL)
		}

	}

}

func do_format_error_response(iso_msg *Iso8583Message) {
	//iso_resp.set_field(39,"APPISO");
	iso_msg.set_field(38, ISO_FORMAT_ERROR)
}
