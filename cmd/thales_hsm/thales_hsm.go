package main

import (
	_ "libiso/hsm"
	"libiso/hsm/console"
	"sync"
)

func main() {

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	thalesConsole := console.New()
	go thalesConsole.Show(waitGroup)

	waitGroup.Wait()

}
