package main

import mynet "go/net"
import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {

	var port int = 0
	var ip string
	var mliStr, data string

	flag.StringVar(&ip, "ip", "127.0.0.1", "destination ip")
	flag.IntVar(&port, "port", 6766, "destination port")
	flag.StringVar(&mliStr, "mli", "2i", "mli if any to be attached to the data")
	flag.StringVar(&data, "data", "30303030", "data in hex")
	flag.Parse()

	/*for i:=0;i<flag.NFlag();i++{
		fmt.Println(flag.Arg(i));
	}*/

	var mli mynet.MliType
	if mliStr == string(mynet.Mli2i) {
		mli = mynet.Mli2i
	} else if mliStr == string(mynet.Mli2e) {
		mli = mynet.Mli2e
	}

	if len(ip) == 0 || port == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	connectionString := ip + ":" + strconv.FormatUint(uint64(port), 10)
	//fmt.Println(connectionString);

	nt := mynet.NewNetCatClient(connectionString, mli)
	err := nt.OpenConnection()
	mynet.HandleError(err)

	//fmt.Println("-->\n",hex.Dump(hex_ba),"\n")
	//write data
	binData, _ := hex.DecodeString(data)
	_ = nt.Write(binData)

	//read response
	readResponse(nt)
	nt.Close()
}

func readResponse(nt *mynet.NetCatClient) {

	var responseData [512]byte
	n, err := nt.Read(responseData[:])
	mynet.HandleError(err)

	fmt.Println("<--\n", hex.Dump(responseData[:n]))

}
