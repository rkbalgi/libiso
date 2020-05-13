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
	Mli4e MliType = "4e"
	Mli4i MliType = "4i"
)

// NetCatClient is network TCP client that can be used to send/receive length-delimited messages
// like 2E,2I, 4E, 4I etc
type NetCatClient struct {
	serverAddr string
	mliType    MliType
	conn       net.Conn
}

// ReadOptions are set of options that can be associated with a connection
type ReadOptions struct {
	Deadline time.Time
}

// NewNetCatClient returns a new netcat client associated with the given mli-type
// and connecting to the addr
func NewNetCatClient(addr string, mliType MliType) *NetCatClient {
	var nt NetCatClient
	nt.mliType = mliType
	nt.serverAddr = addr
	return &nt
}

// OpenConnection opens a connection to the server
func (nt *NetCatClient) OpenConnection() (err error) {
	nt.conn, err = net.Dial("tcp4", nt.serverAddr)
	return err
}

// Close closes the client side of the connection
func (nt *NetCatClient) Close() {
	_ = nt.conn.Close()
}

// Write writes data into the socket after adding the necessary length prefix as per the MLI type
// set on the netcat client
func (nt *NetCatClient) Write(data []byte) (err error) {

	dataWithMli := AddMLI(nt.mliType, data)
	_, err = nt.conn.Write(dataWithMli)
	return err
}

// ReadDirect reads requested data from the socket directly. Use this with caution because the caller is
// responsible for reading length prefix and knowing when one packet starts and ends
func (nt *NetCatClient) ReadDirect(data []byte) (n int, err error) {

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

// ReadNextPacket reads the next data segment (as per MLI type associated with nt)
//
// Deprecated:: Please use Read(*ReadOptions)
//
func (nt *NetCatClient) ReadNextPacket() ([]byte, error) {

	deadline := time.Now().Add(time.Duration(5) * time.Second)

	return nt.Read(&ReadOptions{
		Deadline: deadline,
	})

}

func (nt *NetCatClient) Read(opts *ReadOptions) ([]byte, error) {

	if opts != nil {
		defer nt.conn.SetDeadline(time.Time{})
		nt.conn.SetDeadline(opts.Deadline)
	}

	var mliByteLength uint32 = 0

	switch nt.mliType {
	case Mli2i, Mli2e:
		mliByteLength = 2
	case Mli4i, Mli4e:
		mliByteLength = 4
	}

	tmp := make([]byte, mliByteLength)

	_, err := nt.conn.Read(tmp)
	if err != nil {
		return nil, err
	}

	var msgLen uint32 = 0

	switch nt.mliType {
	case Mli2i, Mli2e:
		msgLen = uint32(binary.BigEndian.Uint16(tmp))
		if nt.mliType == Mli2i {
			msgLen -= mliByteLength
		}
	case Mli4i, Mli4e:
		msgLen = binary.BigEndian.Uint32(tmp)
		if nt.mliType == Mli4i {
			msgLen -= mliByteLength
		}
	}

	msgData := make([]byte, msgLen)
	_, err = nt.conn.Read(msgData)
	if err != nil {
		return nil, err
	}

	return msgData, err

}
