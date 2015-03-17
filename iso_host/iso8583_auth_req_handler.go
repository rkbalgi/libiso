package iso_host

import (
	. "github.com/rkbalgi/go/iso8583"
	"strconv"
	"time"
)

//handles 1100 message
func handle_auth_req(iso_req *Iso8583Message, iso_resp *Iso8583Message) {

	//copy all data from req to response
	//and then lets selective remove fields that are
	//not required in the response
	CopyRequestToResponse(iso_req, iso_resp)

	//set message type on response as 1110

	msg_type_field := iso_resp.GetFieldByName("Message Type")
	msg_type_field.SetData(ISO_MSG_1110)

	//turn off fields not required in response
	iso_resp.Bitmap().SetOff(14)
	iso_resp.Bitmap().SetOff(35)

	//for demo purposes, we will simply base our responses
	//on the input amounts

	f_amount, err := iso_req.Field(4)
	if err != nil {
		do_format_error_response(iso_resp)
		return
	}

	l_amount, err := strconv.ParseUint(f_amount.String(), 10, 64)
	if err != nil {
		logger.Println("invalid amount -",f_amount.String());
		do_format_error_response(iso_resp)
		return
	}
	switch {

	case l_amount > 800 && l_amount < 900:
		{
			iso_resp.SetField(39, ISO_RESP_DECLINE)
		}
	case l_amount == 122:
		{
			iso_resp.SetField(39, ISO_RESP_DROP)
		}
	case l_amount == 123:
		{
			time.Sleep(30 * time.Second)
			iso_resp.SetField(39, ISO_RESP_PICKUP)
		}
	default:
		{
			iso_resp.SetField(38, "APPISO")
			iso_resp.SetField(39, ISO_RESP_APPROVAL)
		}

	}

}

func do_format_error_response(iso_msg *Iso8583Message) {
	iso_msg.SetField(38, ISO_FORMAT_ERROR)
}
