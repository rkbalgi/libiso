package hsm

import (
	"bytes"
)

//globals
var ncResponse []byte = []byte("ND0026860400000000001084-0906")

//handles Thales NC diagnostics command

func (th *ThalesHsm) HandleNC(msgData []byte) []byte {

	responseData := bytes.NewBuffer(msgData[:th.headerLength])

	responseData.Write(ncResponse)
	return responseData.Bytes()

}
