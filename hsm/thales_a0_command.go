package hsm

import (
	"bytes"
	"encoding/hex"
	"go/crypto"
	"go/hsm/keys"
)

type thalesA0Req struct {
	_pro               prologue
	mode               uint //hex1
	keyType            string
	keySchemeLmk       string
	deriveKeyMode      string
	dukptMasterKeyType uint //hex1
	dukptMasterKey     string
	ksn                string //15
	delimiter          string
	zmkTmkFlag         uint //dec1
	zmkTmk             string
	keySchemeZmk       string //1
	attallaVariant     string

	_epi epilogue
}

type thalesA0Resp struct {
	_pro prologue

	responseCode        string
	errorCode           string
	keyUnderLmk         []byte
	keyUnderZmk         []byte
	keyCheckValue       []byte
	endMessageDelimiter byte
	messageTrailer      []byte
}

func (th *ThalesHsm) handleA0Command(msgData []byte) []byte {

	zmkTmkPresent := false
	msgBuf := bytes.NewBuffer(msgData)

	req := new(thalesA0Req)
	resp := new(thalesA0Resp)

	if parsePrologue(msgBuf, &req._pro, th.headerLength) {

		ok := readFixedField(msgBuf, &req.mode, 1, HexadecimalInt)
		ok = readFixedField(msgBuf, &req.keyType, 3, String)
		ok = readFixedField(msgBuf, &req.keySchemeLmk, 1, String)

		if req.mode == 0xa || req.mode == 0xb {
			ok = readFixedField(msgBuf, &req.deriveKeyMode, 1, String)
			if !ok {
				return req.invalidDataResponse(resp)
			} else {
				//derive key mode should be 0
				if req.deriveKeyMode != "0" {
					th.log.Printf("invalid derive key mode - %s", req.deriveKeyMode)
					return req.invalidDataResponse(resp)
				} else {
					//read dupkt master key type and key
					ok = readFixedField(msgBuf, &req.dukptMasterKeyType, 1, HexadecimalInt)
					if !ok {
						return req.invalidDataResponse(resp)
					} else {
						if req.dukptMasterKeyType == 0x01 || req.dukptMasterKeyType == 0x02 {
							ok = readKey(msgBuf, &req.dukptMasterKey)
							if !ok {
								return req.invalidDataResponse(resp)
							} else {
								//read ksn
								ok = readFixedField(msgBuf, &req.ksn, 15, String)
								if !ok {
									return req.invalidDataResponse(resp)
								}
								//check if KSN is all hex, else throw error
								if !hexRegexp.MatchString(req.ksn) {
									th.log.Printf("invalid ksn - %s", req.ksn)
									return req.invalidDataResponse(resp)
								}

							}
						} else {
							th.log.Printf("invalid dukpt master key type - %d", req.dukptMasterKeyType)
							return req.invalidDataResponse(resp)
						}
					}
				}
			}
		}

		if msgBuf.Len() > 0 {
			if req.mode == 0x01 || req.mode == 0x0b {

				ok = readFixedField(msgBuf, &req.delimiter, 1, String)
				if !ok {
					return req.invalidDataResponse(resp)
				}
				if req.delimiter != ";" {
					th.log.Printf("invalid delimiter - %s", req.delimiter)
					return req.invalidDataResponse(resp)
				}
				ok = readFixedField(msgBuf, &req.zmkTmkFlag, 1, DecimalInt)
				if !ok {
					return req.invalidDataResponse(resp)
				}
				if req.zmkTmkFlag == 0 || req.zmkTmkFlag == 1 {
					//ZMK or TMK
					ok = readKey(msgBuf, &req.zmkTmk)

					if !ok {
						return req.invalidDataResponse(resp)
					}
					zmkTmkPresent = true
					ok = readFixedField(msgBuf, &req.keySchemeZmk, 1, String)
					if !ok {
						return req.invalidDataResponse(resp)
					}
					//there may be attalla variant optionally.
					var b, b2 byte
					if msgBuf.Len() > 0 {
						b, _ = msgBuf.ReadByte()
						if b == byte('%') {
							_ = msgBuf.UnreadByte()
						} else {
							//there is attala variant
							if msgBuf.Len() > 0 {
								b2, _ = msgBuf.ReadByte()
								if b2 == byte('%') {
									_ = msgBuf.UnreadByte()
									//just a single digit variant
									req.attallaVariant = string([]byte{b})
								} else {
									req.attallaVariant = string([]byte{b, b2})
								}
							} else {
								//eob - single digit atalla variant
								req.attallaVariant = string([]byte{b})
							}
						}
					}

				} else {
					th.log.Printf("invalid zmk/tmk flag - %d", req.zmkTmkFlag)
					return req.invalidDataResponse(resp)
				}

			}
		}

		ok = parseEpilogue(msgBuf, &req._epi)
		if !ok {
			return req.invalidDataResponse(resp)
		}
	} else {
		//no prolog, message should be dropped
		th.log.Println("[A0] prolog could not be parsed, dropping message")
		return nil
	}

	if hsmDebugEnabled {
		th.log.Println(Dump(*req))
	}

	if req.mode == 0x0a || req.mode == 0x0b {
		//we do not support it at the moment
		th.log.Printf("derive mode (A, B) is not supported at the moment!")
		return req.invalidDataResponse(resp)
	}

	var zmkTmkKey []byte
	var err error

	if zmkTmkPresent {
		if req.zmkTmkFlag == 0 {
			//zmk
			zmkTmkKey, err = decryptKey(req.zmkTmk, ZMK_KEY_TYPE)
			if err != nil {
				th.log.Print("crypto error", err)
				return req.invalidDataResponse(resp)
			}
		} else {
			zmkTmkKey, err = decryptKey(req.zmkTmk, TMK_KEY_TYPE)
			if err != nil {
				th.log.Print("crypto error", err)
				return req.invalidDataResponse(resp)
			}
		}
	}

	keyLen := 0
	switch req.keySchemeLmk {
	case keys.Z:
		keyLen = 8
	case keys.U:
		keyLen = 16
	case keys.T:
		keyLen = 24
	default:
		{
			th.log.Printf("invalid lmk key scheme - %s", req.keySchemeLmk)
			return req.invalidDataResponse(resp)
		}
	}

	//generate the required key and its check value
	key, _ := crypto.GenerateDesKey(keyLen)
	//generate check value
	resp.keyCheckValue, err = genCheckValue(key)

	if err != nil {
		th.log.Print("crypto error", err)
		return req.invalidDataResponse(resp)
	}
	resp.keyCheckValue = resp.keyCheckValue[:3]

	if hsmDebugEnabled {
		th.log.Println("key value: ", hex.EncodeToString(key), "check value: ", hex.EncodeToString(resp.keyCheckValue))
	}

	//TODO:: odd parity enforcement
	if req.keySchemeLmk == keys.Z {
		resp.keyUnderLmk, err = encryptKey(hex.EncodeToString(key), req.keyType)
		if err != nil {
			th.log.Print("crypto error", err)
			return req.invalidDataResponse(resp)
		}
	} else {
		resp.keyUnderLmk, err = encryptKey(req.keySchemeLmk+hex.EncodeToString(key), req.keyType)
		if err != nil {
			th.log.Print("crypto error", err)
			return req.invalidDataResponse(resp)
		}
	}

	if zmkTmkPresent {
		//key should also be encrypted under ZMK/TMK
		switch {
		case req.keySchemeZmk == keys.Z || req.keySchemeZmk == keys.U || req.keySchemeZmk == keys.T:
			{
				resp.keyUnderZmk, err = encryptKeyKek(req.keySchemeZmk+hex.EncodeToString(key), zmkTmkKey, req.keyType)
				if err != nil {
					th.log.Print("crypto error", err)
					return req.invalidDataResponse(resp)
				}
			}
		case req.keySchemeZmk == keys.X || req.keySchemeZmk == keys.Y:
			{
				resp.keyUnderZmk, err = encryptKeyKekX917(hex.EncodeToString(key), zmkTmkKey)
				if err != nil {
					th.log.Print("crypto error", err)
					return req.invalidDataResponse(resp)
				}
				th.log.Println(hex.EncodeToString(resp.keyUnderZmk), "??", hex.EncodeToString(key), "???", hex.EncodeToString(zmkTmkKey))
			}

		default:
			{
				th.log.Printf("invalid zmk key scheme - %s ", req.keySchemeLmk)
				return req.invalidDataResponse(resp)
			}
		}
	}

	//keys should be ascii encoded
	if req.keySchemeLmk == keys.Z {
		//single length keys do not require
		//a scheme identifier
		resp.keyUnderLmk = []byte(hex.EncodeToString(resp.keyUnderLmk))
		if zmkTmkPresent {
			resp.keyUnderZmk = []byte(hex.EncodeToString(resp.keyUnderZmk))
		}
	} else {
		resp.keyUnderLmk = []byte(req.keySchemeLmk + hex.EncodeToString(resp.keyUnderLmk))
		if zmkTmkPresent {
			resp.keyUnderZmk = []byte(req.keySchemeZmk + hex.EncodeToString(resp.keyUnderZmk))
		}
	}
	resp.keyCheckValue = []byte(hex.EncodeToString(resp.keyCheckValue))

	resp.errorCode = HSM_OK

	//generate response
	return req.generateResponse(resp)

}

func (req *thalesA0Req) invalidDataResponse(resp *thalesA0Resp) []byte {

	resp.errorCode = HSM_PARSE_ERROR
	return req.generateResponse(resp)

}

func (req *thalesA0Req) generateResponse(resp *thalesA0Resp) []byte {

	respBuf := bytes.NewBuffer([]byte(req._pro.header))
	respCmdCode := []byte(req._pro.commandName)
	respCmdCode[1] = respCmdCode[1] + 1
	respBuf.Write(respCmdCode)

	respBuf.WriteString(resp.errorCode)
	if resp.errorCode == HSM_OK {

		respBuf.Write(resp.keyUnderLmk)
		if req.mode == 0x01 || req.mode == 0x0b {
			respBuf.Write(resp.keyUnderZmk)
		}
		respBuf.Write(resp.keyCheckValue)
	}

	if req._epi.endMessageDelimiter == 0x19 {
		respBuf.WriteByte(req._epi.endMessageDelimiter)
		respBuf.Write(req._epi.messageTrailer)
	}

	return respBuf.Bytes()

}
