package iso8583

import (
	"container/list"
	"strings"
)

//FieldDefExp represents the ISO8583 field information required
//to be sent as JSON to paysim application
type FieldDefExp struct {
	BitPosition int
	Name        string
	JsSafeName  string
	Data        string
}

//GetMessageDefByName returns a spec definition for the given spec_name
func GetMessageDefByName(spec_name string) *Iso8583MessageDef {
	return spec_map[spec_name]
}

//GetSpecs returns all available specs
func GetSpecs() []string {

	specs := make([]string, len(spec_map))
	i := 0
	for k, _ := range spec_map {
		specs[i] = k
		i = i + 1
	}

	return specs
}

//GetSpecLayout returns all fields associated with a spec
func GetSpecLayout(spec_name string) []*FieldDefExp {

	fields := list.New()
	fields.Init()
	fields.PushBack(&FieldDefExp{BitPosition: 0, Name: "Message Type", JsSafeName: "Message$Type"})
	fields.PushBack(&FieldDefExp{BitPosition: 0, Name: "Bitmap", JsSafeName: "Bitmap"})

	spec := spec_map[spec_name]
	//fmt.Println(spec)
	for i, v := range spec.fields {
		if v != nil {
			//fmt.Println(v.String())
			fields.PushBack(&FieldDefExp{BitPosition: i, Name: v.String(), JsSafeName: js_safe(v.String())})
		}

	}

	field_exp_defs := make([]*FieldDefExp, fields.Len())
	j := 0
	for i := fields.Front(); i != nil; i = i.Next() {
		field_exp_defs[j] = i.Value.(*FieldDefExp)
		j++
	}

	return field_exp_defs

}

func js_safe(in string) string {

	return strings.Replace(in, " ", "$", -1)

}
