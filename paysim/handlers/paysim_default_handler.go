package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/rkbalgi/go/iso8583"
	"github.com/rkbalgi/go/iso8583/services"
	"log"
	"net/http"
)

type PaysimDefaultHandler struct {
}
type ParseTraceHandlerHandler struct {
}

func (handler *PaysimDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	spec_name := req.FormValue("spec_name")
	log.Println("form values::", req.Form)
	log.Println("spec_name::", spec_name)
	layout_json := services.GetSpecLayout(spec_name)
	log.Println("response::  ", layout_json)
	w.Write([]byte(layout_json))
}

//parse_trace_req represents the
//data received from the paysim web application

type parse_trace_req struct {
	Spec_name string
	Data      string
}

func (handler *ParseTraceHandlerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	spec_name := req.FormValue("Spec")
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
	log.Println("parse trace req::  ", json_buf.String())
	req_obj := parse_trace_req{}
	err := json.Unmarshal(json_buf.Bytes(), &req_obj)
	if err != nil {
		log.Println("error parsing request ", err.Error())
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	log.Println("spec_name ", req_obj.Spec_name, " Trace ", req_obj.Data)

	data, err := hex.DecodeString(req_obj.Data)
	if err != nil {
		w.Write([]byte("Error: Invalid Trace Data"))
		return
	}

	in_buf := bytes.NewBuffer(data)
	iso_msg := iso8583.NewIso8583Message(spec_name)
	err = iso_msg.Parse(in_buf)
	if err != nil {
		w.Write([]byte("Error: Parse Failure"))
		return
	}

	field_def_exp_sl := iso8583.GetSpecLayout(req_obj.Spec_name)

	//0 and 1 should be 'Message Type' and 'Bitmap'
	field_def_exp_sl[0].Data = iso_msg.GetMessageType()
	field_def_exp_sl[1].Data = iso_msg.GetBinaryBitmap()

	for _, field_def_exp := range field_def_exp_sl {
		if field_def_exp.BitPosition > 0 {

			if iso_msg.IsSelected(field_def_exp.BitPosition) {
				fld_data, err := iso_msg.GetFieldData(field_def_exp.BitPosition)
				//log.Println(hex.EncodeToString(fld_data));
				if err != nil {
					w.Write([]byte("Error: Parse Failure"))
					return
				}
				field_def_exp.Data = fld_data
			}

		}

	}

	parsed_data_json, err := json.Marshal(field_def_exp_sl)
	if err != nil {
		w.Write([]byte("Error: Parse Failure"))
		return
	}

	log.Println("writing response ", string(parsed_data_json))

	w.Write(parsed_data_json)
}
