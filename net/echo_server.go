package net

import (
	"encoding/hex"
	//"fmt"
	"bytes"
	bin "encoding/binary"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type EchoServ struct {
	TcpAddr *net.TCPAddr
}

//start listening to incoming client connections
//and start off a new goroutine for each client

func (echoServer *EchoServ) ListenAndAccept() (err error) {

	logger.Println("listening at -", echoServer.TcpAddr.String())
	listener, err := net.ListenTCP("tcp4", echoServer.TcpAddr)
	HandleError(err)

	for {
		clientCon, err := listener.Accept()
		HandleError(err)
		//start a new goroutine to handle the
		//client connection
		go handleClient(clientCon)
	}

}

func handleClient(conn net.Conn) {

	defer conn.Close()

	logger.Println("new client connection ", conn.RemoteAddr())

	for {
		var buf [512]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			handleClientError(conn, err)
			return
		}
		if err == nil {
			tmp := buf[:n]
			logger.Println("received -\n", hex.Dump(tmp))
			logger.Println("writing  -\n", hex.Dump(tmp))
			conn.Write(tmp)
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
		logger.Panicf("Error Occurred - %s ", err)

	}
}

// AddMLI adds a MLI to the payload
func AddMLI(mliType MliType, data []byte) []byte {

	switch mliType {

	case Mli2e:
		{
			var mli []byte = make([]byte, 2)
			bin.BigEndian.PutUint16(mli, uint16(len(data)))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return buf.Bytes()
		}
	case Mli2i:
		{
			var mli []byte = make([]byte, 2)
			bin.BigEndian.PutUint16(mli, uint16(len(data)+2))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return buf.Bytes()
		}
	case Mli4e:
		{
			var mli []byte = make([]byte, 4)
			bin.BigEndian.PutUint32(mli, uint32(len(data)))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return buf.Bytes()
		}
	case Mli4i:
		{
			var mli []byte = make([]byte, 4)
			bin.BigEndian.PutUint32(mli, uint32(len(data)+4))
			buf := bytes.NewBuffer(mli)
			buf.Write(data)
			return buf.Bytes()
		}
	default:
		{
			return nil
		}
	}

}

func HandleError(err error) {

	if err != nil {
		logger.Panicf("panic! - %s ", err)
	}
}
