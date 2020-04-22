package iso_host

import (
	"github.com/rkbalgi/libiso/iso8583"
	"strconv"
	"time"
)

//handles 1100 message
func handleAuthReq(isoReq *iso8583.Iso8583Message, isoResp *iso8583.Iso8583Message) {

	//copy all data from req to response
	//and then lets selective remove fields that are
	//not required in the response
	iso8583.CopyRequestToResponse(isoReq, isoResp)

	//set message type on response as 1110

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(iso8583.IsoMsg1110)

	//turn off fields not required in response
	isoResp.Bitmap().SetOff(14)
	isoResp.Bitmap().SetOff(35)

	//for demo purposes, we will simply base our responses
	//on the input amounts

	fAmount, err := isoReq.Field(4)
	if err != nil {
		doFormatErrorResponse(isoResp)
		return
	}

	lAmount, err := strconv.ParseUint(fAmount.String(), 10, 64)
	if err != nil {
		logger.Println("invalid amount -", fAmount.String())
		doFormatErrorResponse(isoResp)
		return
	}
	switch {

	case lAmount > 800 && lAmount < 900:
		{
			isoResp.SetField(39, iso8583.IsoRespDecline)
		}
	case lAmount == 122:
		{
			isoResp.SetField(39, iso8583.IsoRespDrop)
		}
	case lAmount == 123:
		{
			time.Sleep(30 * time.Second)
			isoResp.SetField(39, iso8583.IsoRespPickup)
		}
	default:
		{
			isoResp.SetField(38, "APPISO")
			isoResp.SetField(39, iso8583.IsoRespApproval)
		}

	}

}

func doFormatErrorResponse(isoMsg *iso8583.Iso8583Message) {
	isoMsg.SetField(38, iso8583.IsoFormatError)
}
