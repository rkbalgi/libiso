package main

import (
	mynet "go/net"
	"log"
)
import (
	"fmt"
	"net"
	//"testing"
	"os"
	"strconv"
)

var cmdLineArgs = make(map[string]string, 10)

const (
	portArg = "-port"
	ipArg   = "-ip"
)

func main() {

	cmdLineArgs[portArg] = ""
	cmdLineArgs[ipArg] = ""

	var argument string

	for i, v := range os.Args {

		if i == 0 {
			continue
		}

		if v[0] == '-' {
			if _, ok := cmdLineArgs[v]; ok {
				//so  val is a known argument type
				argument = v

			} else {
				//unknown argument
				printUsageAndQuit()
			}
		} else {
			cmdLineArgs[argument] = v
		}
	}

	//fmt.Println("arguments -", cmdLineArgs)

	if len(os.Args) != 2*2+1 {
		printUsageAndQuit()
	}

	echoServ := new(mynet.EchoServ)
	echoServ.TcpAddr = new(net.TCPAddr)
	echoServ.TcpAddr.IP = net.ParseIP(cmdLineArgs[ipArg])
	echoServ.TcpAddr.Port, _ = strconv.Atoi(cmdLineArgs[portArg])
	if err := echoServ.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

}

func printUsageAndQuit() {
	fmt.Printf("usage: %s [-port] [-ip]", os.Args[0])
	os.Exit(0)
}
