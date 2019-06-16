package paysim

import (
	"fmt"
)

type Tuple struct {
	data []interface{}
}

func (t *Tuple) Nth(pos int) interface{} {
	return t.data[pos]
}

func NewTuple(data ...interface{}) *Tuple {

	tuple := new(Tuple)
	tuple.data = make([]interface{}, len(data))
	copy(tuple.data, data)
	return tuple

}

func (t *Tuple) String() string {
	return fmt.Sprintln(t.data)
}
