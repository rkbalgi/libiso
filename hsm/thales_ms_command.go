package hsm

import (
	"bytes"
	"encoding/hex"
	"github.com/rkbalgi/go/crypto/mac"
	"log"
	//_ "github.com/rkbalgi/hsm"
)

type ThalesMsRequest struct {
	_prologue          prologue
	messageBlockNumber uint
	keyType            uint
	keyLength          uint
	messageType        uint
	key                string
	iv                 string
	messageLength      uint
	messageBlock       []byte
	_epilogue          epilogue
}

type ThalesMsResponse struct {
	messageHeader       []byte
	responseCode        []byte
	errorCode           []byte
	MAB                 []byte
	endMessageDelimiter string
	messageTrailer      string
}

func (th *ThalesHsm) HandleMS(msgData []byte) []byte {

	defer func() {
		str := recover()
		if str != nil {
			th.log.Println("unexpected system error", str)
		}
	}()

	msReqStruct := new(ThalesMsRequest)
	msRespStruct := new(ThalesMsResponse)

	msgBuf := bytes.NewBuffer(msgData)
	//at every step this boolean will be updated
	//if false, respond with a 15 error code
	parseOk := true

	parsePrologue(msgBuf, &msReqStruct._prologue, th.headerLength)

	parseOk = readFixedField(msgBuf, &msReqStruct.messageBlockNumber, 1, DecimalInt)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	parseOk = readFixedField(msgBuf, &msReqStruct.keyType, 1, DecimalInt)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	parseOk = readFixedField(msgBuf, &msReqStruct.keyLength, 1, DecimalInt)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	//message type
	parseOk = readFixedField(msgBuf, &msReqStruct.messageType, 1, DecimalInt)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	//key
	parseOk = readKey(msgBuf, &msReqStruct.key)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	//iv
	if msReqStruct.messageBlockNumber == 2 || msReqStruct.messageBlockNumber == 3 {
		parseOk = readFixedField(msgBuf, &msReqStruct.iv, 16, String)
		if !parseOk {
			return msReqStruct.InvalidDataResponse(msRespStruct)
		}

	}

	//message length

	parseOk = readFixedField(msgBuf, &msReqStruct.messageLength, 4, HexadecimalInt)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	//message block

	parseOk = readFixedField(msgBuf, &msReqStruct.messageBlock, msReqStruct.messageLength, Binary)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	parseOk = parseEpilogue(msgBuf, &msReqStruct._epilogue)
	if !parseOk {
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	if hsmDebugEnabled {
		th.log.Printf(Dump(*msReqStruct))
	}

	//decrypt the key under the appropriate LMK
	keyType := "000"
	if msReqStruct.keyType == 0 {
		//TAK
		keyType = "003"
	} else if msReqStruct.keyType == 1 {
		//ZAK
		keyType = "008"
	} else {
		//error
		th.log.Printf("invalid key type - %s\n", keyType)
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	var (
		macKey []byte
		err    error
	)
	if macKey, err = decryptKey(msReqStruct.key, keyType); err != nil {
		log.Print("Error decrypting ", err)
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	if (msReqStruct.keyLength == 0 && (len(macKey) == 16 || len(macKey) == 24)) || (msReqStruct.keyLength == 1 && len(macKey) == 8) {
		th.log.Println("key length and actual size mismatch.")
		return msReqStruct.InvalidDataResponse(msRespStruct)
	}

	if hsmDebugEnabled {
		th.log.Println("mac key", hex.EncodeToString(macKey))
	}
	var genMac []byte

	if len(macKey) == 16 || len(macKey) == 24 {
		//generate X9.19 mac
		genMac, err = mac.GenerateMacX919(msReqStruct.messageBlock, macKey)
		if err != nil {
			th.log.Print("crypto error", err)
			msReqStruct.InvalidDataResponse(msRespStruct)
		}
	} else {
		//x9.9
		genMac, err = mac.GenerateMacX99(msReqStruct.messageBlock, macKey)
		if err != nil {
			th.log.Print("crypto error", err)
			msReqStruct.InvalidDataResponse(msRespStruct)
		}

	}
	msRespStruct.MAB = genMac
	msRespStruct.errorCode = []byte("00")

	if hsmDebugEnabled {
		th.log.Printf("MAC: %s", hex.EncodeToString(genMac))
	}

	//send response
	return msRespStruct.generateResponse(msReqStruct)

}

func (resp *ThalesMsRequest) InvalidDataResponse(msRespStruct *ThalesMsResponse) []byte {
	setFixedField(&msRespStruct.errorCode, 2, uint(15), DecimalInt)
	return msRespStruct.generateResponse(resp)

}

func (resp *ThalesMsResponse) generateResponse(msReqStruct *ThalesMsRequest) []byte {

	msRespBuf := bytes.NewBuffer([]byte(msReqStruct._prologue.header))
	respCmdCode := []byte(msReqStruct._prologue.commandName)
	respCmdCode[1] = respCmdCode[1] + 1
	msRespBuf.Write(respCmdCode)

	msRespBuf.Write(resp.errorCode)
	if resp.MAB != nil {
		msRespBuf.Write(toASCII(resp.MAB))
	}

	if msReqStruct._epilogue.endMessageDelimiter == 0x19 {
		msRespBuf.WriteByte(msReqStruct._epilogue.endMessageDelimiter)
		msRespBuf.Write(msReqStruct._epilogue.messageTrailer)
	}

	return msRespBuf.Bytes()

}
