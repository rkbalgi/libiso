package iso8583

import (
	"container/list"
	"fmt"
	"strings"
)

//FieldDefExp represents the ISO8583 field information required
//to be sent as JSON to paysim application
type FieldDefExp struct {
	BitPosition int
	Name        string
	JsSafeName  string
	Data        string
	Id          int
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

	spec := spec_map[spec_name]
	fields := list.New()
	fields.Init()

	fmt.Println(spec.fields_def_list.Len())

	for l := spec.fields_def_list.Front(); l != nil; l = l.Next() {
		switch (l.Value).(type) {
		case IsoField:
			{
				var iso_field IsoField = (l.Value).(IsoField)
				fields.PushBack(&FieldDefExp{Id: iso_field.GetId(), BitPosition: 0, Name: iso_field.String(), JsSafeName: js_safe(iso_field.String())})
				break
			}
		case BitmappedField:
			{
				var iso_bmp_field *BitMap = (l.Value).(*BitMap)
				fields.PushBack(&FieldDefExp{Id: iso_bmp_field.GetId(), BitPosition: 0, Name: "Bitmap", JsSafeName: "Bitmap"})

				for b_position, f_def := range iso_bmp_field.sub_field_def {
					if f_def != nil {
						fields.PushBack(&FieldDefExp{Id: f_def.GetId(), BitPosition: b_position, Name: f_def.String(), JsSafeName: js_safe(f_def.String())})
					}
				}

			}
		} //end of switch

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
