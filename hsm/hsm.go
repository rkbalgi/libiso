package hsm

import (
	"net"
	//"strconv"
	"bufio"
	//"bytes"
	"encoding/binary"
	"encoding/hex"
	_net "github.com/rkbalgi/libiso/net"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// ThalesHsm represents a software Thales HSM
type ThalesHsm struct {
	ip           string
	port         int
	encodingType EncodingType
	headerLength int
	stop         bool
	listener     *net.TCPListener
	log          *log.Logger
}

func NewThalesHsm(ip string, port int, encodingType EncodingType) *ThalesHsm {
	th := new(ThalesHsm)
	th.headerLength = 12
	th.port = port
	//th.ip = ip
	th.encodingType = encodingType
	th.log = log.New(os.Stdout, "## thales_hsm (8000/9000) ## ", log.LstdFlags)
	return th
}

func (th *ThalesHsm) Stop() {
	th.stop = true
	_ = th.listener.Close()
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
			if strings.Contains(err.Error(), "closed network connection") {
				th.log.Print("hsm being stopped..\n")
				return
			}
			th.log.Panicf("failed to start thales hsm - [%s]", err.Error())
			return
		}
		th.log.Printf("new client connection - %s", tcpConn.RemoteAddr().String())

		//start a new goroutine to read data off the
		//connection and create appropriate responses

		go th.msgReader(tcpConn)

	}

}

func (th *ThalesHsm) handleClientMsg(tcpConn *net.TCPConn, msgData []byte) {

	if msgData == nil || len(msgData) == 0 {
		th.log.Printf("invalid request from client")
		return
	}

	th.log.Printf("request from client - \n%s", hex.Dump(msgData))
	commandName := string(msgData[12 : 12+2])

	var responseData []byte

	switch commandName {
	case "NC":
		{
			responseData = th.HandleNC(msgData)
			break
		}
	case "MS":
		{
			responseData = th.HandleMS(msgData)
			break
		}
	case "CC":
		{
			responseData = th.handleCcCommand(msgData)
			break
		}
	default:
		{
			th.log.Printf("unsupported command [%s] received. message dropped", commandName)
			return
		}
	}

	if responseData != nil {
		th.log.Printf("writing response to client - \n%s\n", hex.Dump(responseData))
	} else {
		th.log.Printf("no response to write.")
	}
	responseData = _net.AddMLI(_net.Mli2e, responseData)
	_, err := tcpConn.Write(responseData)
	if th.checkError(err) {
		return
	}

}

func (th *ThalesHsm) msgReader(tcpConn *net.TCPConn) {

	var reader io.Reader = tcpConn
	defer func() {
		str := recover()
		if str != nil {
			th.log.Println("(recovered)", str)
		}
	}()

	for {
		//time.Sleep(time.Second * 2)

		//first read the 2E mli
		tmp := make([]byte, 2)
		_, err := reader.Read(tmp)
		if th.checkError(err) {
			//if connection has been closed
			//return
			return
		}

		msgLen := binary.BigEndian.Uint16(tmp)
		//read data
		msgData := make([]byte, msgLen)
		_, err = reader.Read(msgData)
		if th.checkError(err) {
			//if connection has been closed
			//return
			return
		}

		//new goroutine will handle the message
		go th.handleClientMsg(tcpConn, msgData)

	}

}

func (th *ThalesHsm) bufferedMsgReader(reader io.Reader) {
	bufMsgReader := bufio.NewReader(reader)

	for {
		time.Sleep(time.Second * 2)

		bufferedDataLen := bufMsgReader.Buffered()
		if bufferedDataLen > 2 {

			th.log.Println("buffered bytes3 ", bufferedDataLen)
			tmp, err := bufMsgReader.Peek(5)
			if err != nil {
				th.checkError(err)
			}
			th.log.Println(tmp)

			msgLen := binary.BigEndian.Uint16(tmp)
			th.log.Printf("message len - %d", msgLen)
			completeMsgLen := 2 + int(msgLen)
			if bufferedDataLen >= completeMsgLen {
				//we have enough bytes to make a
				//complete message
				msg := make([]byte, completeMsgLen)
				n, _ := bufMsgReader.Read(msg)
				th.log.Printf("new msg from client - \n%s", hex.Dump(msg))
				if n != completeMsgLen {
					th.log.Printf("read error, expected to read [%d] bytes., found only [%d]", completeMsgLen, n)
				}

			}
		}
	}

}

//check if err is not nil and return true if the client
//connection has been closed

func (th *ThalesHsm) checkError(err error) bool {
	if err != nil {

		if err.Error() == "EOF" {
			//closed connection, close silently
			th.log.Println("connection closed by client (EOF).")
			return true
		}

		th.log.Printf("error -%s", err.Error())
		if strings.Contains(err.Error(), "forcibly closed") {
			return true
		}

	}
	return false
}
