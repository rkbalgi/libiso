package net

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

type TcpMessageHandler interface {
	HandleMessage(client_con *net.TCPConn, msg_data []byte)
}

type TcpHost struct {
	mli_type MliType
	addr     *net.TCPAddr
	handler  TcpMessageHandler
	log      *log.Logger
}

func NewTcpHost(mli_type MliType, tcp_addr *net.TCPAddr) *TcpHost {
	new_tcp_host := new(TcpHost)
	new_tcp_host.addr = tcp_addr
	new_tcp_host.mli_type = mli_type
	new_tcp_host.log = log.New(os.Stdout, "tcp_host ## ", log.LstdFlags)

	return new_tcp_host

}

func (tcp_host *TcpHost) SetHandler(handler TcpMessageHandler) {
	tcp_host.handler = handler
}

func (tcp_host *TcpHost) Start() {

	tcp_listener, err := net.ListenTCP("tcp4", tcp_host.addr)
	if err != nil {
		fmt.Println("error listening at port -", tcp_host.addr.String(), err.Error())
		return
	}
	logger.Printf("started tcp-ip host @ %s", tcp_listener.Addr().String())

	for {
		client_conn, err := tcp_listener.AcceptTCP()
		if err != nil {
			fmt.Println("error accepting client connection - ", err.Error())
		}
		logger.Printf("new connection from %s", client_conn.RemoteAddr().String())

		go tcp_host.handle_client(client_conn)

	}

}

func (tcp_host *TcpHost) handle_client(client_conn *net.TCPConn) {

	//if any errors/panics handling the client, just recover and let
	//others connected continue
	defer func() {

		str := recover()
		if str != nil {
			tcp_host.log.Printf("panic recovered. client connection will be closed.")
			return
		}

	}()

	mli := make([]byte, 2)

	for {

		n, err := client_conn.Read(mli)
		if err != nil {
			handle_network_error(err, client_conn.RemoteAddr().String())
			client_conn.Close()
			return
		}

		req_len := binary.BigEndian.Uint16(mli)
		if tcp_host.mli_type == MLI_2I {
			req_len = req_len - 2
		}
		tcp_host.log.Printf("reading incoming message with %d bytes...\n", req_len)

		//all good, read rest of the message data
		msg_data := make([]byte, req_len)
		n, err = client_conn.Read(msg_data)
		if err != nil {
			handle_network_error(err, client_conn.RemoteAddr().String())
			client_conn.Close()
			return
		}

		if uint16(n) != req_len {
			logger.Printf("not enough data - required: %d != actual %d\n", req_len, n)
			continue
		}

		go tcp_host.handler.HandleMessage(client_conn, msg_data)

	}
}

func handle_network_error(err error, ref_msg string) {

	if err != nil {

		if err.Error() == "EOF" {
			log.Printf("client connection closed -[ref: %s]", ref_msg)
		} else {

			logger.Printf("error on client connection. closing connection [Err: %s]\n", err.Error())
		}

		return
	}
}
