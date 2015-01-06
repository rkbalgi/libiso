package hsm

import (
	"net"
	//"strconv"
	"bufio"
	//"bytes"
	"encoding/binary"
	"encoding/hex"
	_net "github.com/rkbalgi/go/net"
	"io"
	"log"
	"os"
	_ "strconv"
	"strings"
	"time"
)

type ThalesHsm struct {
	ip            string
	port          int
	encoding_type EncodingType
	header_length int
	stop          bool
	listener *net.TCPListener
	log           *log.Logger
}

func NewThalesHsm(ip string, port int, encoding_type EncodingType) *ThalesHsm {
	th := new(ThalesHsm)
	th.header_length = 12
	th.port = port
	//th.ip = ip
	th.encoding_type = encoding_type
	th.log = log.New(os.Stdout, "## thales_hsm (8000/9000) ## ", log.LstdFlags)
	return th
}

func (th *ThalesHsm) Stop() {
	th.stop = true
	th.listener.Close();
}
func (th *ThalesHsm) Start() {

	addr := new(net.TCPAddr)
	addr.IP = net.ParseIP(th.ip)
	addr.Port = th.port
	th.stop = false

	th.log.Println("starting thales hsm -", addr.String())
	var err error
	th.listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		th.log.Panicf("failed to start thales hsm - [%s]", err.Error())
		return
	}
	
	for !th.stop {

		tcpConn, err := th.listener.AcceptTCP()
		if err != nil {
			if strings.Contains(err.Error(),"closed network connection"){
				th.log.Print("hsm being stopped..\n");
				return;
			}
			th.log.Panicf("failed to start thales hsm - [%s]", err.Error())
			return
		}
		th.log.Printf("new client connection - %s", tcpConn.RemoteAddr().String())

		//start a new goroutine to read data off the
		//connection and create appropriate responses

		go th.msg_reader(tcpConn)

	}

}

func (hsm_handle *ThalesHsm) handle_client_msg(tcp_conn *net.TCPConn, msg_data []byte) {

	hsm_handle.log.Printf("request from client - \n%s", hex.Dump(msg_data))
	command_name := string(msg_data[12 : 12+2])

	var response_data []byte

	switch command_name {
	case "NC":
		{
			response_data = hsm_handle.Handle_NC(msg_data)
			break
		}
	case "MS":
		{
			response_data = hsm_handle.Handle_MS(msg_data)
			break
		}
	case "CC":
		{
			response_data = hsm_handle.handle_cc_command(msg_data)
			break
		}
	default:
		{
			hsm_handle.log.Printf("unsupported command [%s] received. message dropped", command_name)
			return
		}
	}

	if response_data != nil {
		hsm_handle.log.Printf("writing response to client - \n%s\n", hex.Dump(response_data))
	} else {
		hsm_handle.log.Printf("no response to write.")
	}
	response_data = _net.AddMLI(_net.MLI_2E, response_data)
	_, err := tcp_conn.Write(response_data)
	if hsm_handle.check_error(err) {
		return
	}

}

func (th *ThalesHsm) msg_reader(tcp_conn *net.TCPConn) {

	var reader io.Reader = tcp_conn
	for {
		//time.Sleep(time.Second * 2)

		//first read the 2E mli
		tmp := make([]byte, 2)
		_, err := reader.Read(tmp)
		if th.check_error(err) {
			//if connection has been closed
			//return
			return
		}

		msg_len := binary.BigEndian.Uint16(tmp)
		//read data
		msg_data := make([]byte, msg_len)
		_, err = reader.Read(msg_data)
		if th.check_error(err) {
			//if connection has been closed
			//return
			return
		}

		//new goroutine will handle the message
		go th.handle_client_msg(tcp_conn, msg_data)

	}

}

func (th *ThalesHsm) buffered_msg_reader(reader io.Reader) {
	buf_msg_reader := bufio.NewReader(reader)

	for {
		time.Sleep(time.Second * 2)

		buffered_data_len := buf_msg_reader.Buffered()
		if buffered_data_len > 2 {

			th.log.Println("buffered bytes3 ", buffered_data_len)
			tmp, err := buf_msg_reader.Peek(5)
			if err != nil {
				th.check_error(err)
			}
			th.log.Println(tmp)

			msg_len := binary.BigEndian.Uint16(tmp)
			th.log.Printf("message len - %s", msg_len)
			complete_msg_len := 2 + int(msg_len)
			if buffered_data_len >= complete_msg_len {
				//we have enough bytes to make a
				//complete message
				msg := make([]byte, complete_msg_len)
				n, _ := buf_msg_reader.Read(msg)
				th.log.Printf("new msg from client - \n%s", hex.Dump(msg))
				if n != complete_msg_len {
					th.log.Printf("read error, expected to read [%d] bytes., found only [%d]", complete_msg_len, n)
				}

			}
		}
	}

}

//check if err is not nil and return true if the client
//connection has been closed

func (th *ThalesHsm) check_error(err error) bool {
	if err != nil {
		th.log.Printf("error -%s", err.Error())
		if strings.Contains(err.Error(), "forcibly closed") {
			return true
		}

	}
	return false
}
