package net

import (
	"encoding/hex"
	//"fmt"
	"log"
	"net"
	"os"
	"bytes"
	bin "encoding/binary"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type EchoServ struct {
	TcpAddr *net.TCPAddr
}

//start listening to incoming client connections
//and start off a new goroutine for each client

func (echoServ *EchoServ) ListenAndAccept() (err error) {

	logger.Println("listening at -", echoServ.TcpAddr.String())
	listener, err := net.ListenTCP("tcp4", echoServ.TcpAddr)
	HandleError(err)

	for {
		clientCon, err := listener.Accept()
		HandleError(err)
		//start a new goroutine to handle the
		//client connection
		go handleClient(clientCon)
	}

}

func handleClient(clientCon net.Conn) {

	defer clientCon.Close()

	logger.Println("new client connection ", clientCon.RemoteAddr())

	for {
		var buf [512]byte
		n, err := clientCon.Read(buf[:])
		if err != nil {
			handleClientError(clientCon, err)
			return
		}
		if err == nil {
			tmp := buf[:n]
			logger.Println("received -\n", hex.Dump(tmp))
			logger.Println("writing  -\n", hex.Dump(tmp))
			clientCon.Write(tmp)
		} else {
			logger.Println("no data", err)
		}
	}

}

func handleClientError(clientCon net.Conn, err error) {

	if err.Error() == "EOF" {
		clientCon.Close()
		return
	}
	if err != nil {
		defer clientCon.Close()
		logger.Panicf("Error Occured - %s ", err)

	}
}

/**
This function adds a mli (length indicator) based on the type of mli - 2I, 2L etc
**/
func AddMLI(mliType MliType, data []byte) []byte {

	switch mliType {

	case MLI_2E:
		{
			var mli []byte = make([]byte, 2)
			bin.BigEndian.PutUint16(mli, uint16(len(data)))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return (buf.Bytes())
		}
	case MLI_2I:
		{
			var mli []byte = make([]byte, 2)
			bin.BigEndian.PutUint16(mli, uint16(len(data)+2))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return (buf.Bytes())
		}
		default:{
			return nil
		}
	}

}

func HandleError(err error) {

	if err != nil {
		logger.Panicf("panic! - %s ", err)
		os.Exit(1)
	}
}
