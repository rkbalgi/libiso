package iso8583

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_readSpecDef(t *testing.T) {

	specs, err := readSpecDef(filepath.Join(".", "testdata"), "iso_specs.yaml")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(specs))
	t.Log(specs)
}
