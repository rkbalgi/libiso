package hsm

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type thalesCCReq struct {
	_pro                 prologue
	sourceZpk            string
	destinationZpk       string
	maximumPinLen        uint
	sourcePinBlock       string //change type to []byte
	sourcePinblkFmt      uint
	destinationPinblkFmt uint
	acctNo               string

	_epi epilogue
}

type thalesCCResp struct {
	_pro prologue

	responseCode         string
	errorCode            string
	pinLen               []byte
	destinationPinBlock  []byte
	destinationPinblkFmt []byte
	endMessageDelimiter  byte
	messageTrailer       []byte
}

func (th *ThalesHsm) handleCcCommand(msgData []byte) []byte {

	var (
		msgBuf = bytes.NewBuffer(msgData)
		err    error
	)
	req := new(thalesCCReq)
	resp := new(thalesCCResp)

	if parsePrologue(msgBuf, &req._pro, th.headerLength) {
		parseOk := readKey(msgBuf, &req.sourceZpk)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = readKey(msgBuf, &req.destinationZpk)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = readFixedField(msgBuf, &req.maximumPinLen, 2, DecimalInt)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = readFixedField(msgBuf, &req.sourcePinBlock, 16, String)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = readFixedField(msgBuf, &req.sourcePinblkFmt, 2, DecimalInt)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = readFixedField(msgBuf, &req.destinationPinblkFmt, 2, DecimalInt)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		//TODO:should be 18H for formats 04
		parseOk = readFixedField(msgBuf, &req.acctNo, 12, String)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = parseEpilogue(msgBuf, &req._epi)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
		parseOk = parseEpilogue(msgBuf, &req._epi)
		if !parseOk {
			return req.invalidDataResponse(resp)
		}
	} else {
		//no prolog, message should be dropped
		th.log.Println("[CC] prolog could not be parsed, dropping message")
		return nil
	}

	inKey, err := decryptKey(req.sourceZpk, "001")
	if err != nil {
		th.log.Printf("crypto error - %s", err)
		return req.invalidDataResponse(resp)
	}
	outKey, err := decryptKey(req.destinationZpk, "001")
	if err != nil {
		th.log.Printf("crypto error - %s", err)
		return req.invalidDataResponse(resp)
	}
	if hsmDebugEnabled {
		th.log.Println(Dump(*req))
	}

	//handle message i.e translate pin block
	ph := newPinHandler(req.acctNo, int(req.sourcePinblkFmt), int(req.destinationPinblkFmt), inKey, outKey)
	inPinBlock, err := hex.DecodeString(req.sourcePinBlock)
	if err != nil {
		th.log.Printf("invalid source pin block - %s", req.sourcePinBlock)
		return req.invalidDataResponse(resp)
	}
	outPinBlk, err := ph.translate(inPinBlock)
	if err != nil {
		th.log.Printf("crypto error - %s", err)
		return req.invalidDataResponse(resp)
	}
	if outPinBlk != nil && len(outPinBlk) == 8 {
		resp.errorCode = HSM_OK
		if hsmDebugEnabled {
			th.log.Println("clear pin -", ph.getClearPin())
			th.log.Println("destination pin block - ", hex.EncodeToString(outPinBlk))
		}
		resp.destinationPinBlock = []byte(hex.EncodeToString(outPinBlk))
		resp.pinLen = []byte(fmt.Sprintf("%02d", len(ph.getClearPin())))
	} else {
		return req.invalidDataResponse(resp)
	}

	//generate response
	return req.generateResponse(resp)

}

func (req *thalesCCReq) invalidDataResponse(resp *thalesCCResp) []byte {

	resp.errorCode = HSM_PARSE_ERROR
	return req.generateResponse(resp)

}

func (req *thalesCCReq) generateResponse(resp *thalesCCResp) []byte {

	respBuf := bytes.NewBuffer([]byte(req._pro.header))
	respCmdCode := []byte(req._pro.commandName)
	respCmdCode[1] = respCmdCode[1] + 1
	respBuf.Write(respCmdCode)

	respBuf.WriteString(resp.errorCode)
	if resp.errorCode == HSM_OK {
		respBuf.Write(resp.pinLen)
		respBuf.Write(resp.destinationPinBlock)
		respBuf.WriteString(fmt.Sprintf("%02d", req.destinationPinblkFmt))

	}
	if req._epi.endMessageDelimiter == 0x19 {
		respBuf.WriteByte(req._epi.endMessageDelimiter)
		respBuf.Write(req._epi.messageTrailer)
	}

	return respBuf.Bytes()

}
