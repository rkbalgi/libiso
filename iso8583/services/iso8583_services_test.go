package services

import (
	"testing"
)

func Test_GetSpecsTest(t *testing.T) {

	specs_json := GetSpecs()
	t.Log(specs_json)
}

func Test_GetSpecsLayoutTest(t *testing.T) {

	specs_json := GetSpecLayout("ISO8583_1 v1 (ASCII)")
	t.Log(specs_json)
}
