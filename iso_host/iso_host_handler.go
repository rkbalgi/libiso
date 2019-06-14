package iso_host

import (
	"bytes"
	//"github.com/rkbalgi/go/iso8583"
	"errors"
	. "github.com/rkbalgi/go/iso8583"
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "##iso_handler##", log.LstdFlags)

//this method handles an incoming ISO8583 message, doing the parsing, processing
//and response creation
func Handle(spec_name string, buf *bytes.Buffer) (resp_iso_msg *Iso8583Message, err error) {

	req_iso_msg := NewIso8583Message(spec_name)

	//parse incoming message
	err = req_iso_msg.Parse(buf)
	if err != nil {
		return nil, err
	}

	logger.Println("parsed incoming message: ", req_iso_msg.Dump())

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

	f39, err := resp_iso_msg.Field(39)

	if f39.String() == ISO_RESP_DROP {
		logger.Println("drop response code, not sending response..")
		return nil, errors.New("dropped response!")
	}

	logger.Println("outgoing message: ", resp_iso_msg.Dump())

	return resp_iso_msg, err

}
