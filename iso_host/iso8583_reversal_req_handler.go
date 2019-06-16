package iso_host

import (
	. "github.com/rkbalgi/go/iso8583"
)

func handleReversalReq(isoResp *Iso8583Message) {

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(IsoMsg1430)

}
