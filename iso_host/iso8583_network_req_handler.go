package iso_host

import (
	. "github.com/rkbalgi/go/iso8583"
)

func handle_network_req(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	msg_type_field := iso_resp.GetFieldByName("Message Type")
	msg_type_field.SetData(ISO_MSG_1814)

}
