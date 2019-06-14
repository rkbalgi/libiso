package paysim

import (
	"fmt"
)

type Tuple struct {
	data []interface{}
}

func (self *Tuple) Nth(pos int) interface{} {
	return self.data[pos]
}

func NewTuple(data ...interface{}) *Tuple {

	tuple := new(Tuple)
	tuple.data = make([]interface{}, len(data))
	copy(tuple.data, data)
	return tuple

}

func (self *Tuple) String() string {
	return fmt.Sprintln(self.data)
}
