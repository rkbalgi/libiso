package handlers

import (
	"bytes"
	_ "encoding/hex"
	_ "encoding/json"
	_ "github.com/rkbalgi/go/iso8583"
	_ "github.com/rkbalgi/go/iso8583/services"
	"log"
	"net/http"
)

type SendMessageHandlerHandler struct {
}

func (handler *SendMessageHandlerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	spec_name := req.FormValue("spec_name")
	log.Println("form values::", req.Form)
	log.Println("spec_name::", spec_name)

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
	log.Println("send msg req::  ", json_buf.String())

	w.Write([]byte("done."))
}
