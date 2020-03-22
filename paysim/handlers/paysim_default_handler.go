package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"go/iso8583"
	"go/iso8583/services"
	"io"
	"log"
	"net/http"
)

type PaysimDefaultHandler struct {
}
type ParseTraceHandlerHandler struct {
}

func (handler *PaysimDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	specName := req.FormValue("spec_name")
	log.Println("form values::", req.Form)
	log.Println("spec_name::", specName)
	layoutJson := services.GetSpecLayout(specName)
	log.Println("response::  ", layoutJson)
	_, _ = w.(io.StringWriter).WriteString(layoutJson)
}

//parse_trace_req represents the
//data received from the paysim web application

type parseTraceReq struct {
	SpecName string
	Data     string
}

func (handler *ParseTraceHandlerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	specName := req.FormValue("Spec")
	log.Println("form values::", req.Form)
	log.Println("spec_name::", specName)

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
	log.Println("parse trace req::  ", jsonBuf.String())
	reqObj := parseTraceReq{}
	err := json.Unmarshal(jsonBuf.Bytes(), &reqObj)
	if err != nil {
		log.Println("error parsing request ", err.Error())
		_, _ = w.Write([]byte("Error: " + err.Error()))
		return
	}

	log.Println("spec_name ", reqObj.SpecName, " Trace ", reqObj.Data)

	data, err := hex.DecodeString(reqObj.Data)
	if err != nil {
		_, _ = w.Write([]byte("Error: Invalid Trace Data"))
		return
	}

	inBuf := bytes.NewBuffer(data)
	isoMsg := iso8583.NewIso8583Message(specName)
	err = isoMsg.Parse(inBuf)
	if err != nil {
		_, _ = w.Write([]byte("Error: Parse Failure"))
		return
	}

	fieldDefExpSl := iso8583.GetSpecLayout(reqObj.SpecName)

	//0 and 1 should be 'Message Type' and 'Bitmap'
	fieldDefExpSl[0].Data = isoMsg.GetMessageType()
	fieldDefExpSl[1].Data = isoMsg.GetBinaryBitmap()

	for _, fieldDefExp := range fieldDefExpSl {
		if fieldDefExp.BitPosition > 0 {

			if isoMsg.IsSelected(fieldDefExp.BitPosition) {
				fldData, err := isoMsg.GetFieldData(fieldDefExp.BitPosition)
				//log.Println(hex.EncodeToString(fld_data));
				if err != nil {
					_, _ = w.Write([]byte("Error: Parse Failure"))
					return
				}
				fieldDefExp.Data = fldData
			}

		}

	}

	parsedDataJson, err := json.Marshal(fieldDefExpSl)
	if err != nil {
		_, _ = w.Write([]byte("Error: Parse Failure"))
		return
	}

	log.Println("writing response ", string(parsedDataJson))

	_, _ = w.Write(parsedDataJson)
}
