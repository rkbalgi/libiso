package hsm

import (
	"encoding/hex"
	"log"
	"os"
	"testing"
)

func Test_A0ParseTest_GenOnly(t *testing.T) {

	sHsmCmd := "000000000001;A0;0;003;U"
	hsmHandle := NewThalesHsm("127.0.0.1", 1500, AsciiEncoding)
	hsmHandle.log = log.New(os.Stdout, "##???## ", log.LstdFlags)
	hsmHandle.headerLength = 12

	respData := hsmHandle.handle_a0_command(formatHsmCommand(sHsmCmd))
	t.Log("response_data ", hex.EncodeToString(respData))
	t.Log("response_data (ascii)", string(respData))

}

func Test_A0ParseTest_GenAndExport(t *testing.T) {

	sHsmCmd := "000000000001;A0;1;008;U;x'3b';0;U0C999BC58C997CE279FC6427041AF9B7;X;%00"
	hsmHandle := NewThalesHsm("127.0.0.1", 1500, AsciiEncoding)
	hsmHandle.log = log.New(os.Stdout, "##???## ", log.LstdFlags)
	hsmHandle.headerLength = 12

	respData := hsmHandle.handle_a0_command(formatHsmCommand(sHsmCmd))
	t.Log("response_data ", hex.EncodeToString(respData))
	t.Log("response_data (ascii)", "\n", hex.Dump(respData))

}
