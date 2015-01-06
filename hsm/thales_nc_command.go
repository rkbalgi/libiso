package hsm

import(
	"bytes"
)
//globals
var nc_response []byte=[]byte("ND0026860400000000001084-0906");

//handles Thales NC diagnostics command

func (hsm_handle *ThalesHsm)Handle_NC(msg_data []byte) []byte{
	
	response_data:=bytes.NewBuffer(msg_data[:hsm_handle.header_length]);
	
	response_data.Write(nc_response);
	return(response_data.Bytes());
	
}