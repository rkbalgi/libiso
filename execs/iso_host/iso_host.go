package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/rkbalgi/go/iso8583"
	"github.com/rkbalgi/go/iso_host"
	bnet "github.com/rkbalgi/go/net"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stdout, "##iso_host## ", log.LstdFlags)

type IsoMessageHandler struct {
	spec_name string
}

func (iso_msg_handler *IsoMessageHandler) HandleMessage(client_conn *net.TCPConn, msg_data []byte) {

	logger.Println("handling request = \n", hex.Dump(msg_data))

	buf := bytes.NewBuffer(msg_data)
	resp_iso_msg, err := iso_host.Handle(iso_msg_handler.spec_name, buf)
	if err != nil {
		log.Printf("error handling message from client -[Err: %s]\n", err.Error())
		return
	}

	resp_data := resp_iso_msg.Bytes()

	//add mli + resp_data into buffer
	mli := make([]byte, 2)
	binary.BigEndian.PutUint16(mli, uint16(len(resp_data)+2))
	resp_buf := bytes.NewBuffer(mli)
	resp_buf.Write(resp_data)

	logger.Println("writing response = \n", hex.Dump(resp_buf.Bytes()))
	client_conn.Write(resp_buf.Bytes())

}

func main() {

	port := flag.Int("port", 5656, "port to listen at")
	spec_name := flag.String("spec", "ISO8583_1_v1__DEMO_", "specification from the spec file")
	spec_def_file_name := flag.String("spec-file", "", "file to read the specifications from")

	flag.Parse()
	if len(*spec_def_file_name) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	file_handle, err := os.Open(*spec_def_file_name)
	if err != nil {
		fmt.Println("unable to open spec-file")
		flag.Usage()
		os.Exit(1)
	}
	file_handle.Close()

	iso8583.ReadSpecDefs(*spec_def_file_name)

	tcp_addr := new(net.TCPAddr)
	tcp_addr.IP = net.ParseIP("")
	tcp_addr.Port = *port

	iso_host := bnet.NewTcpHost(bnet.Mli2i, tcp_addr)
	handler := new(IsoMessageHandler)
	handler.spec_name = *spec_name
	iso_host.SetHandler(handler)

	iso_host.Start()

}
