package services

import (
	"encoding/json"
	"github.com/rkbalgi/go/iso8583"
	"log"
)

//GetSpecs returns all available specs
func GetSpecs() string {

	specs := iso8583.GetSpecs()
	json, err := json.Marshal(specs)
	if err != nil {
		log.Println("failed to marshall to JSON - ", err.Error())
		return ""
	}

	return string(json)

}

//GetSpecLayout returns the list of all fields in spec_name
//as a JSON string
func GetSpecLayout(spec_name string) string {

	fields := iso8583.GetSpecLayout(spec_name)

	json, err := json.Marshal(fields)
	if err != nil {
		log.Println("failed to marshall to JSON - ", err.Error())
		return ""
	}

	return string(json)

}
