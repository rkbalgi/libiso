package services

import (
	"testing"
)

func Test_GetSpecsTest(t *testing.T) {

	specsJson := GetSpecs()
	t.Log(specsJson)
}

func Test_GetSpecsLayoutTest(t *testing.T) {

	specsJson := GetSpecLayout("ISO8583_1 v1 (ASCII)")
	t.Log(specsJson)
}
