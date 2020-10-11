package services

import (
	"github.com/rkbalgi/libiso/iso8583"
	"path/filepath"
	"testing"
)

func init() {
	iso8583.ReadSpecDefs(filepath.Join(".", "testdata", "sample_spec.json"))
}

func Test_GetSpecsTest(t *testing.T) {

	specsJson := GetSpecs()
	t.Log(specsJson)
}

func Test_GetSpecsLayoutTest(t *testing.T) {

	specsJson := GetSpecLayout("ISO8583_1_v1__DEMO_")
	t.Log(specsJson)
}
