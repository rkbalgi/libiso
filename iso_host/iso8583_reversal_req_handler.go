package iso_host

func handleReversalReq(isoResp *Iso8583Message) {

	msgTypeField := isoResp.GetFieldByName("Message Type")
	msgTypeField.SetData(IsoMsg1430)

}
