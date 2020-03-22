package iso_host

func handleNetworkReq(isoResp *Iso8583Message) {

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(IsoMsg1814)

}
