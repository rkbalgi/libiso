package net

import (
	"bytes"
	"encoding/hex"
	"errors"
	"go/iso8583"
	p_nc "go/net"
	pylog "go/paysim/log"
	//	"log"
	//	"os"
)

var ncMap map[string]*p_nc.NetCatClient

//var net_logger *log.Logger;

func init() {
	ncMap = make(map[string]*p_nc.NetCatClient, 10)
	//net_logger=log.New(os.Stdout,"#iso_net >> ",log.LstdFlags);
}

//SendIsoMsg will send a msg over to the host
//receive and return the parsed response message
//TODO:: there is need to keep a map of the requests/responses
//and manage them via a key in case there are multiple messages
func SendIsoMsg(connectionStr string,
	mli string,
	isoMsg *iso8583.Iso8583Message) (*iso8583.Iso8583Message, error) {

	var mliType p_nc.MliType

	if mli == "2I" {
		mliType = p_nc.Mli2i
	} else {
		mliType = p_nc.Mli2e
	}

	nc, ok := ncMap[connectionStr]
	//lets check if nc is still connected
	if ok && !nc.IsConnected() {
		pylog.Printf("an existing connection [%s] has been closed. opening new connection.", connectionStr)
		delete(ncMap, connectionStr)
		nc.Close()
		ok = false
	}

	if !ok {
		nc = p_nc.NewNetCatClient(connectionStr, mliType)
		err := nc.OpenConnection()
		if err != nil {
			return nil, err
		}
		pylog.Log("new tcp/ip connection opened to -", connectionStr)
		ncMap[connectionStr] = nc
	}

	//we have a client  now
	reqMsgData := isoMsg.Bytes()
	pylog.Log("sending data \n", hex.Dump(reqMsgData), "\n")
	_ = nc.Write(reqMsgData)
	respMsgData, err := nc.ReadNextPacket()
	if err != nil {
		return nil, err
	}
	pylog.Log("received data \n", hex.Dump(respMsgData), "\n")

	respIsoMsg := iso8583.NewIso8583Message(isoMsg.SpecName())
	msgBuf := bytes.NewBuffer(respMsgData)
	err = respIsoMsg.Parse(msgBuf)
	if err != nil {
		return nil, errors.New("error parsing response data")
	}

	return respIsoMsg, nil

}
