package net

import (
	"encoding/binary"
	"net"
	"time"
)

type MliType string

const (
	Mli2i MliType = "2i"
	Mli2e MliType = "2e"
)

type NetCatClient struct {
	connectionStr string
	mliType       MliType
	conn          net.Conn
}

func NewNetCatClient(connectionStr string, mliType MliType) *NetCatClient {
	var nt NetCatClient
	nt.mliType = mliType
	nt.connectionStr = connectionStr
	return &nt
}

func (nt *NetCatClient) OpenConnection() (err error) {
	nt.conn, err = net.Dial("tcp4", nt.connectionStr)
	return err
}

func (nt *NetCatClient) Close() {
	_ = nt.conn.Close()
}

func (nt *NetCatClient) Write(data []byte) (err error) {

	dataWithMli := AddMLI(nt.mliType, data)
	_, err = nt.conn.Write(dataWithMli)
	return err
}

func (nt *NetCatClient) Read(data []byte) (n int, err error) {

	n, err = nt.conn.Read(data)
	return n, err
}

func (nt *NetCatClient) IsConnected() bool {

	defer func() {
		_ = nt.conn.SetReadDeadline(time.Time{})
	}()
	_ = nt.conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))
	_, err := nt.conn.Read(make([]byte, 0))
	if err != nil {
		return false
	}

	return true
}

func (nt *NetCatClient) ReadNextPacket() ([]byte, error) {

	defer func() {
		_ = nt.conn.SetReadDeadline(time.Time{})
	}()
	_ = nt.conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))

	tmp := make([]byte, 2)
	_, err := nt.conn.Read(tmp)
	if err != nil {
		//if connection has been closed
		//return
		return nil, err
	}

	msgLen := binary.BigEndian.Uint16(tmp)
	if nt.mliType == Mli2i {
		msgLen -= 2
	}
	//read data
	msgData := make([]byte, msgLen)
	_, err = nt.conn.Read(msgData)
	if err != nil {
		//if connection has been closed
		//return
		return nil, err
	}

	return msgData, err

}
