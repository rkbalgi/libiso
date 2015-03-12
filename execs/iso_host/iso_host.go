package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/rkbalgi/go/iso8583"
	bnet "github.com/rkbalgi/go/net"

	"log"
	"net"
	"os"
)

var logger = log.New(os.Stdout, "##iso8583## ", log.LstdFlags)

type IsoMessageHandler struct {
}

func (iso_msg_handler *IsoMessageHandler) HandleMessage(client_conn *net.TCPConn, msg_data []byte) {

	logger.Println("handling request = \n", hex.Dump(msg_data))

	buf := bytes.NewBuffer(msg_data)
	resp_iso_msg, err := iso8583.Handle("ISO8583_1 v1 (ASCII)",buf)
	if err != nil {
		log.Printf("error handling message from client -[Err: %s]\n", err.Error)
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
	
	iso8583.ReadDemoSpecDefs();

	tcp_addr := new(net.TCPAddr)
	tcp_addr.IP = net.ParseIP("127.0.0.1")
	tcp_addr.Port = 5656

	iso_host := bnet.NewTcpHost(bnet.MLI_2I, tcp_addr)
	iso_host.SetHandler(new(IsoMessageHandler))

	iso_host.Start()

}
