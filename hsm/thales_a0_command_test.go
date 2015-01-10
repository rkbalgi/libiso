package hsm


import(
	"testing"
	"encoding/hex"
	"log"
	"os"
	
)


func Test_A0ParseTest_GenOnly(t *testing.T){
	
	s_hsm_cmd:="000000000001;A0;0;003;U"
	hsm_handle:=NewThalesHsm("127.0.0.1",1500,AsciiEncoding);
	hsm_handle.log=log.New(os.Stdout,"##???## ",log.LstdFlags)
	hsm_handle.header_length=12
	
	resp_data:=hsm_handle.handle_a0_command(format_hsm_command(s_hsm_cmd));
	t.Log("response_data ",hex.EncodeToString(resp_data))
	t.Log("response_data (ascii)",string(resp_data))
	
}

func Test_A0ParseTest_GenAndExport(t *testing.T){
	
	s_hsm_cmd:="000000000001;A0;1;008;U;x'3b';0;U0C999BC58C997CE279FC6427041AF9B7;X;%00"
	hsm_handle:=NewThalesHsm("127.0.0.1",1500,AsciiEncoding);
	hsm_handle.log=log.New(os.Stdout,"##???## ",log.LstdFlags)
	hsm_handle.header_length=12
	
	resp_data:=hsm_handle.handle_a0_command(format_hsm_command(s_hsm_cmd));
	t.Log("response_data ",hex.EncodeToString(resp_data))
	t.Log("response_data (ascii)","\n",hex.Dump(resp_data))
	
}