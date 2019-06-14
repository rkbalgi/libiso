package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/rkbalgi/go/iso8583"
	_ "github.com/rkbalgi/go/iso8583/services"
	"github.com/rkbalgi/go/iso_host"
	"log"
	"net/http"
)

type SendMessageHandlerHandler struct {
}

func (handler *SendMessageHandlerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	buf := make([]byte, 100)
	json_buf := bytes.NewBufferString("")
	for {
		n, err := req.Body.Read(buf)

		if n > 0 {
			//we have good data
			json_buf.Write(buf[:n])
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
	log.Println("send msg request::  ", json_buf.String())
	snd_req := iso8583.WebMsgData{Type: "Request"}
	json.Unmarshal(json_buf.Bytes(), &snd_req)
	log.Println(snd_req.Spec, snd_req.DataArray)

	//iso_msg_def:=iso8583.GetMessageDefByName(snd_req.Spec);

	iso_msg := iso8583.NewIso8583Message(snd_req.Spec)

	iso_msg.SetData(snd_req.DataArray)

	log.Println("received request : ", iso_msg.Dump())

	//handle the incoming message
	req_data := iso_msg.Bytes()
	log.Println("req: \n", hex.EncodeToString(req_data))
	msg_buf := bytes.NewBuffer(req_data)
	resp_iso_msg, err := iso_host.Handle(snd_req.Spec, msg_buf)
	log.Println("processed response: \n", resp_iso_msg.Dump())

	if err != nil {
		w.Write([]byte("error"))
		return
	}

	web_msg_data := resp_iso_msg.ToWebMsg(false)
	json_data, err := json.Marshal(web_msg_data)
	if err == nil {
		log.Println("writing json ", string(json_data))
		w.Write(json_data)
	} else {
		log.Println("error marshalling json -", err.Error())
		w.Write([]byte("error"))

	}

}
