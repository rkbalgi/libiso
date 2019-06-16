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
func GetMessageDefByName(specName string) *MessageDef {
	return specMap[specName]
}

//GetSpecNames returns names of all available specs
func GetSpecNames() []string {

	specs := make([]string, len(specMap))
	i := 0
	for k, _ := range specMap {
		specs[i] = k
		i = i + 1
	}

	return specs
}

//GetSpecs returns all available specs
func GetSpecs() []*MessageDef {

	specs := make([]*MessageDef, len(specMap))
	i := 0
	for _, v := range specMap {
		specs[i] = v
		i = i + 1
	}

	return specs
}

//GetSpecLayout returns all fields associated with a spec
func GetSpecLayout(specName string) []*FieldDefExp {

	spec := specMap[specName]
	fields := list.New()
	fields.Init()

	fmt.Println(spec.fieldsDefList.Len())

	for l := spec.fieldsDefList.Front(); l != nil; l = l.Next() {
		switch (l.Value).(type) {
		case IsoField:
			{
				var isoField IsoField = (l.Value).(IsoField)
				fields.PushBack(&FieldDefExp{Id: isoField.GetId(), BitPosition: 0, Name: isoField.String(), JsSafeName: js_safe(isoField.String())})
				break
			}
		case BitmappedField:
			{
				var isoBmpField *BitMap = (l.Value).(*BitMap)
				fields.PushBack(&FieldDefExp{Id: isoBmpField.GetId(), BitPosition: 0, Name: "Bitmap", JsSafeName: "Bitmap"})

				for bPosition, fDef := range isoBmpField.subFieldDef {
					if fDef != nil {
						fields.PushBack(&FieldDefExp{Id: fDef.GetId(), BitPosition: bPosition, Name: fDef.String(), JsSafeName: js_safe(fDef.String())})
					}
				}

			}
		} //end of switch

	}

	fieldExpDefs := make([]*FieldDefExp, fields.Len())
	j := 0
	for i := fields.Front(); i != nil; i = i.Next() {
		fieldExpDefs[j] = i.Value.(*FieldDefExp)
		j++
	}

	return fieldExpDefs

}

func js_safe(in string) string {

	return strings.Replace(in, " ", "$", -1)

}
