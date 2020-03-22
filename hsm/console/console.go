package console

import (
	"bufio"
	"fmt"
	"go/hsm"
	//"io"
	_ "flag"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var startCmdRegexp *regexp.Regexp

func init() {
	var err error
	startCmdRegexp, err = regexp.Compile("start_hsm[ ]+-port[ ]+([0-9]+)")
	if err != nil {
		panic(err.Error())
	}
}

type Console struct {
	thalesHsm *hsm.ThalesHsm
}

const (
	EXIT = "exit"
	QUIT = "quit"
)

func New() *Console {
	return &Console{}
}

func (console *Console) Show(waitGroup *sync.WaitGroup) {
	var line string
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("console# ")
		line, _ = stdinReader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == QUIT || line == EXIT {
			break

		} else {
			console.handleCommand(line)

		}

	}
	waitGroup.Done()
}

func (console *Console) handleCommand(cmd string) {

	if len(cmd) == 0 {
		return
	}

	if startCmdRegexp.MatchString(cmd) {
		subMatches := startCmdRegexp.FindStringSubmatch(cmd)
		port, _ := strconv.ParseInt(subMatches[1], 10, 32)
		console.thalesHsm = hsm.NewThalesHsm("127.0.0.1", int(port), hsm.AsciiEncoding)
		go console.thalesHsm.Start()
		fmt.Println("done.")
	} else if cmd == "stop_hsm" {
		console.thalesHsm.Stop()
		fmt.Println("done.")
	} else {
		fmt.Println("bad command.")
	}

	//return "done."
}
