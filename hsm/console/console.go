package console

import (
	"bufio"
	"fmt"
	"github.com/rkbalgi/go/hsm"
	//"io"
	_ "flag"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var start_cmd_regexp *regexp.Regexp

func init() {
	var err error
	start_cmd_regexp, err = regexp.Compile("start_hsm[ ]+-port[ ]+([0-9]+)")
	if err != nil {
		panic(err.Error())
	}
}

type Console struct {
	thales_hsm *hsm.ThalesHsm
}

const (
	EXIT = "exit"
	QUIT = "quit"
)

func New() *Console {
	return (&Console{})
}

func (console *Console) Show(wait_group *sync.WaitGroup) {
	var line string
	stdin_reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("console# ")
		line, _ = stdin_reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == QUIT || line == EXIT {
			break

		} else {
			console.handle_command(line)

		}

	}
	wait_group.Done()
}

func (console *Console) handle_command(cmd string) {
	
	if(len(cmd)==0){
		return;
	}

	if start_cmd_regexp.MatchString(cmd) {
		sub_matches := start_cmd_regexp.FindStringSubmatch(cmd)
		port, _ := strconv.ParseInt(sub_matches[1], 10, 32)
		console.thales_hsm = hsm.NewThalesHsm("127.0.0.1", int(port), hsm.AsciiEncoding)
		go console.thales_hsm.Start()
		fmt.Println("done.")
	} else if cmd == "stop_hsm" {
		console.thales_hsm.Stop()
		fmt.Println("done.")
	} else {
		fmt.Println("bad command.")
	}

	//return "done."
}
