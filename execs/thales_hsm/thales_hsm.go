package main

import (
	_ "github.com/rkbalgi/go/hsm"
	"github.com/rkbalgi/go/hsm/console"
	"sync"
)

func main() {

	wait_group := new(sync.WaitGroup)
	wait_group.Add(1)

	thales_console := console.New()
	go thales_console.Show(wait_group)

	wait_group.Wait()

}
