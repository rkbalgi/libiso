package net

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

type TcpMessageHandler interface {
	HandleMessage(clientCon *net.TCPConn, msgData []byte)
}

type TcpHost struct {
	mliType MliType
	addr    *net.TCPAddr
	handler TcpMessageHandler
	log     *log.Logger
}

func NewTcpHost(mliType MliType, tcpAddr *net.TCPAddr) *TcpHost {
	newTcpHost := new(TcpHost)
	newTcpHost.addr = tcpAddr
	newTcpHost.mliType = mliType
	newTcpHost.log = log.New(os.Stdout, "tcp_host ## ", log.LstdFlags)

	return newTcpHost

}

func (tcpHost *TcpHost) SetHandler(handler TcpMessageHandler) {
	tcpHost.handler = handler
}

func (tcpHost *TcpHost) Start() {

	tcpListener, err := net.ListenTCP("tcp4", tcpHost.addr)
	if err != nil {
		fmt.Println("error listening at port -", tcpHost.addr.String(), err.Error())
		return
	}
	logger.Printf("started tcp-ip host @ %s", tcpListener.Addr().String())

	for {
		clientConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("error accepting client connection - ", err.Error())
		}
		logger.Printf("new connection from %s", clientConn.RemoteAddr().String())

		go tcpHost.handleClient(clientConn)

	}

}

func (tcpHost *TcpHost) handleClient(clientConn *net.TCPConn) {

	//if any errors/panics handling the client, just recover and let
	//others connected continue
	defer func() {

		str := recover()
		if str != nil {
			tcpHost.log.Printf("panic recovered. client connection will be closed: %s", str)
			return
		}

	}()

	mli := make([]byte, 2)

	for {

		n, err := clientConn.Read(mli)
		if err != nil {
			handleNetworkError(err, clientConn.RemoteAddr().String())
			_ = clientConn.Close()
			return
		}

		reqLen := binary.BigEndian.Uint16(mli)
		if tcpHost.mliType == Mli2i {
			reqLen = reqLen - 2
		}
		tcpHost.log.Printf("reading incoming message with %d bytes...\n", reqLen)

		//all good, read rest of the message data
		msgData := make([]byte, reqLen)
		n, err = clientConn.Read(msgData)
		if err != nil {
			handleNetworkError(err, clientConn.RemoteAddr().String())
			_ = clientConn.Close()
			return
		}

		if uint16(n) != reqLen {
			logger.Printf("not enough data - required: %d != actual %d\n", reqLen, n)
			continue
		}

		go tcpHost.handler.HandleMessage(clientConn, msgData)

	}
}

func handleNetworkError(err error, refMsg string) {

	if err != nil {

		if err.Error() == "EOF" {
			log.Printf("client connection closed -[ref: %s]", refMsg)
		} else {

			logger.Printf("error on client connection. closing connection [Err: %s]\n", err.Error())
		}

		return
	}
}
