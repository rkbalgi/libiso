package log

import (
	"fmt"
	"github.com/rkbalgi/go/paysim/ui"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "## paysim.ui >>", log.LstdFlags)

func Log(log_msg ...interface{}) {

	logger.Println(log_msg)

	if ui.PaysimConsole != nil {
		msg := fmt.Sprintln(log_msg)
		ui.PaysimConsole.Log(msg)
	}
}

func Printf(fmt_str string, log_msg ...interface{}) {

	logger.Printf(fmt_str, log_msg)

	if ui.PaysimConsole != nil {
		msg := fmt.Sprintf(fmt_str, log_msg)
		ui.PaysimConsole.Log(msg)
	}

}
