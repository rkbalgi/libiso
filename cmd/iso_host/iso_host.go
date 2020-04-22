package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"libiso/iso8583"
	"libiso/iso_host"
	bnet "libiso/net"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stdout, "##iso_host## ", log.LstdFlags)

type IsoMessageHandler struct {
	specName string
}

func (isoMsgHandler *IsoMessageHandler) HandleMessage(clientConn *net.TCPConn, msgData []byte) {

	logger.Println("handling request = \n", hex.Dump(msgData))

	buf := bytes.NewBuffer(msgData)
	respIsoMsg, err := iso_host.Handle(isoMsgHandler.specName, buf)
	if err != nil {
		log.Printf("error handling message from client -[Err: %s]\n", err.Error())
		return
	}

	respData := respIsoMsg.Bytes()

	//add mli + resp_data into buffer
	mli := make([]byte, 2)
	binary.BigEndian.PutUint16(mli, uint16(len(respData)+2))
	respBuf := bytes.NewBuffer(mli)
	respBuf.Write(respData)

	logger.Println("writing response = \n", hex.Dump(respBuf.Bytes()))
	_, _ = clientConn.Write(respBuf.Bytes())

}

func main() {

	port := flag.Int("port", 5656, "port to listen at")
	specName := flag.String("spec", "ISO8583_1_v1__DEMO_", "specification from the spec file")
	specDefFileName := flag.String("spec-file", "", "file to read the specifications from")

	flag.Parse()
	if len(*specDefFileName) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fileHandle, err := os.Open(*specDefFileName)
	if err != nil {
		fmt.Println("unable to open spec-file")
		flag.Usage()
		os.Exit(1)
	}
	err = fileHandle.Close()

	iso8583.ReadSpecDefs(*specDefFileName)

	tcpAddr := new(net.TCPAddr)
	tcpAddr.IP = net.ParseIP("")
	tcpAddr.Port = *port

	isoHost := bnet.NewTcpHost(bnet.Mli2i, tcpAddr)
	handler := new(IsoMessageHandler)
	handler.specName = *specName
	isoHost.SetHandler(handler)

	isoHost.Start()

}
