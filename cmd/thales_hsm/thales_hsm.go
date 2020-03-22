package main

import (
	_ "go/hsm"
	"go/hsm/console"
	"sync"
)

func main() {

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	thalesConsole := console.New()
	go thalesConsole.Show(waitGroup)

	waitGroup.Wait()

}
