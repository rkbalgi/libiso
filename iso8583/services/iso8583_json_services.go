package services

import (
	"encoding/json"
	"go/iso8583"
	"log"
)

//GetSpecs returns all available specs
func GetSpecs() string {

	specs := iso8583.GetSpecs()
	jsonContent, err := json.Marshal(specs)
	if err != nil {
		log.Println("failed to marshall to JSON - ", err.Error())
		return ""
	}

	return string(jsonContent)

}

//GetSpecLayout returns the list of all fields in spec_name
//as a JSON string
func GetSpecLayout(specName string) string {

	fields := iso8583.GetSpecLayout(specName)

	jsonContent, err := json.Marshal(fields)
	if err != nil {
		log.Println("failed to marshall to JSON - ", err.Error())
		return ""
	}

	return string(jsonContent)

}
