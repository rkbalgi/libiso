package iso_host

import (
	"strconv"
	"time"
)

//handles 1100 message
func handleAuthReq(isoReq *Iso8583Message, isoResp *Iso8583Message) {

	//copy all data from req to response
	//and then lets selective remove fields that are
	//not required in the response
	CopyRequestToResponse(isoReq, isoResp)

	//set message type on response as 1110

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(IsoMsg1110)

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
			isoResp.SetField(39, IsoRespDecline)
		}
	case lAmount == 122:
		{
			isoResp.SetField(39, IsoRespDrop)
		}
	case lAmount == 123:
		{
			time.Sleep(30 * time.Second)
			isoResp.SetField(39, IsoRespPickup)
		}
	default:
		{
			isoResp.SetField(38, "APPISO")
			isoResp.SetField(39, IsoRespApproval)
		}

	}

}

func doFormatErrorResponse(isoMsg *Iso8583Message) {
	isoMsg.SetField(38, IsoFormatError)
}
