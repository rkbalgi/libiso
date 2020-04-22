package iso_host

import (
	"bytes"
	"errors"
	. "libiso/iso8583"
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "##iso_handler##", log.LstdFlags)

//this method handles an incoming ISO8583 message, doing the parsing, processing
//and response creation
func Handle(specName string, buf *bytes.Buffer) (respIsoMsg *Iso8583Message, err error) {

	reqIsoMsg := NewIso8583Message(specName)

	//parse incoming message
	err = reqIsoMsg.Parse(buf)
	if err != nil {
		return nil, err
	}

	logger.Println("parsed incoming message: ", reqIsoMsg.Dump())

	//continue handling

	respIsoMsg = NewIso8583Message(specName)
	msgType := reqIsoMsg.GetMessageType()
	switch msgType {
	case IsoMsg1100:
		{
			handleAuthReq(reqIsoMsg, respIsoMsg)
		}
	case IsoMsg1804:
		{
			handleNetworkReq(respIsoMsg)
		}
	case IsoMsg1420:
		{
			handleReversalReq(respIsoMsg)
		}
	default:
		{
			err = errors.New("unsupported message type -" + reqIsoMsg.GetMessageType())

		}
	}

	f39, err := respIsoMsg.Field(39)
	if err != nil {
		return nil, err
	}

	if f39.String() == IsoRespDrop {
		logger.Println("drop response code, not sending response..")
		return nil, errors.New("dropped response")
	}

	logger.Println("outgoing message: ", respIsoMsg.Dump())

	return respIsoMsg, err

}
