package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"go/iso8583"
	_ "go/iso8583/services"
	"go/iso_host"
	"log"
	"net/http"
)

type SendMessageHandlerHandler struct {
}

func (handler *SendMessageHandlerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	buf := make([]byte, 100)
	jsonBuf := bytes.NewBufferString("")
	for {
		n, err := req.Body.Read(buf)

		if n > 0 {
			//we have good data
			jsonBuf.Write(buf[:n])
		}

		if err != nil && n == 0 {
			break
		}

		if err != nil && err.Error() == "EOF" {
			//log.Println("error ::  ", err)
			break
		} else if err != nil {
			log.Println("error ::  ", err)
			return
		}

	}
	log.Println("send msg request::  ", jsonBuf.String())
	sndReq := iso8583.WebMsgData{Type: "Request"}
	_ = json.Unmarshal(jsonBuf.Bytes(), &sndReq)
	log.Println(sndReq.Spec, sndReq.DataArray)

	//iso_msg_def:=iso8583.GetMessageDefByName(snd_req.Spec);

	isoMsg := iso8583.NewIso8583Message(sndReq.Spec)

	isoMsg.SetData(sndReq.DataArray)

	log.Println("received request : ", isoMsg.Dump())

	//handle the incoming message
	reqData := isoMsg.Bytes()
	log.Println("req: \n", hex.EncodeToString(reqData))
	msgBuf := bytes.NewBuffer(reqData)
	respIsoMsg, err := iso_host.Handle(sndReq.Spec, msgBuf)
	if err != nil {
		_, _ = w.Write([]byte("error: " + err.Error()))
		return
	}
	log.Println("processed response: \n", respIsoMsg.Dump())

	webMsgData := respIsoMsg.ToWebMsg(false)
	jsonData, err := json.Marshal(webMsgData)
	if err == nil {
		log.Println("writing json ", string(jsonData))
		_, _ = w.Write(jsonData)
	} else {
		log.Println("error marshalling json -", err.Error())
		_, _ = w.Write([]byte("error"))

	}

}
