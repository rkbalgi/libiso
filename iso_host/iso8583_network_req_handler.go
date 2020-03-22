package iso_host

import "go/iso8583"

func handleNetworkReq(isoResp *iso8583.Iso8583Message) {

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(iso8583.IsoMsg1814)

}
