package iso8583

func handle_reversal_req(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	msg_type_field := iso_resp.get_field_by_name("Message Type")
	msg_type_field.SetData(ISO_MSG_1430)

}
