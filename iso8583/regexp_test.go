package iso8583

import (
	"strings"
	"testing"
)

func Test_ConstraintsRegExp(t *testing.T) {

	str := "constraints{minSize:10;maxSize:15;Content:Alpha;}"

	constraints := make(map[string]interface{}, 10)
	if constraintsRegExp1.MatchString(str) {
		targetString := str[strings.Index(str, "{")+1 : len(str)-1]
		matches := constraintsRegExp2.FindAllStringSubmatch(targetString, -1)
		for _, match := range matches {
			constraints[match[2]] = match[3]
			t.Log(match[2], match[3])
		}
	} else {
		t.Fail()
	}

}
