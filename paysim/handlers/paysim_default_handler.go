package handlers

import (
	"github.com/rkbalgi/go/iso8583/services"
	"log"
	"net/http"
)

type PaysimDefaultHandler struct {
}

func (handler *PaysimDefaultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	spec_name := req.FormValue("spec_name")
	log.Println("form values::",req.Form);
	log.Println("spec_name::", spec_name)
	layout_json := services.GetSpecLayout(spec_name)
	log.Println("response::  ", layout_json)
	w.Write([]byte(layout_json))
}
