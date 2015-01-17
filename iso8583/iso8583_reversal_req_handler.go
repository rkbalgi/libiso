package iso8583

func handle_reversal_req(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	iso_resp.msg_type = ISO_MSG_1430
}
