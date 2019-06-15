package hsm

import (
	"encoding/hex"
	"github.com/rkbalgi/go/crypto"
	"golang.org/x/xerrors"
	"strconv"
)

//generate a check value for the given key
func genCheckValue(keyData []byte) ([]byte, error) {

	data := make([]byte, 8)
	if len(keyData) == 8 {
		return crypto.EncryptDes(data, keyData)
	} else {
		return crypto.EncryptTripleDes(data, keyData)
	}

}

//extract key data (as []byte) from a encoded string
func extractKeyData(keyStr string) []byte {

	var keyData []byte
	if keyStr[0] == 'U' || keyStr[0] == 'Z' || keyStr[0] == 'T' || keyStr[0] == 'X' || keyStr[0] == 'Y' {
		keyData, _ = hex.DecodeString(keyStr[1:])
	} else {
		keyData, _ = hex.DecodeString(keyStr[:])
	}
	return keyData
}

//encrypt a key under the lmk using thales variants

func encryptKey(keyStr string, keyType string) ([]byte, error) {

	var keyData []byte
	if keyStr[0] == 'U' || keyStr[0] == 'Z' || keyStr[0] == 'T' {
		keyData, _ = hex.DecodeString(keyStr[1:])
	} else {
		keyData, _ = hex.DecodeString(keyStr[:])
	}

	lmkKey := keyTypeTable[keyType[1:]]
	if lmkKey == nil {
		panic("unsupported key type" + keyType)
	}

	return encryptKeyThalesV(keyData, lmkKey, keyType)
}

//encrypt a key under the kek using thales variants
func encryptKeyKek(keyStr string, kek []byte, keyType string) ([]byte, error) {

	keyData := extractKeyData(keyStr)
	return encryptKeyThalesV(keyData, kek, keyType)
}

func encryptKeyThalesV(keyData []byte, lmkKey []byte, keyType string) ([]byte, error) {

	vKey := make([]byte, 16)
	copy(vKey, lmkKey)

	setVariants(vKey, keyType)

	if len(keyData) == 24 {
		//left
		vKey[8] = vKey[8] ^ variantTripleLen1
		left, err := crypto.EncryptTripleDes(keyData[0:8], vKey)
		if err != nil {
			return nil, err
		}

		//now middle
		vKey[8] = vKey[8] ^ variantTripleLen1
		vKey[8] = vKey[8] ^ variantTripleLen2
		middle, err := crypto.EncryptTripleDes(keyData[8:16], vKey)
		if err != nil {
			return nil, err
		}
		//right
		vKey[8] = vKey[8] ^ variantTripleLen2
		vKey[8] = vKey[8] ^ variantTripleLen3
		right, err := crypto.EncryptTripleDes(keyData[16:], vKey)
		if err != nil {
			return nil, err
		}
		clearKey := make([]byte, 0)
		clearKey = append(clearKey, left...)
		clearKey = append(clearKey, middle...)
		clearKey = append(clearKey, right...)
		return clearKey, nil

	} else if len(keyData) == 16 {
		vKey[8] = vKey[8] ^ variantDblLen1
		leftHalf, err := crypto.EncryptTripleDes(keyData[0:8], vKey)
		if err != nil {
			return nil, err
		}
		//now right half
		vKey[8] = vKey[8] ^ variantDblLen1
		vKey[8] = vKey[8] ^ variantDblLen2
		rightHalf, err := crypto.EncryptTripleDes(keyData[8:], vKey)
		if err != nil {
			return nil, err
		}
		clearKey := make([]byte, 0)
		clearKey = append(clearKey, leftHalf...)
		clearKey = append(clearKey, rightHalf...)

		return clearKey, nil

	} else if len(keyData) == 8 {
		//single length
		clearKey, _ := crypto.EncryptTripleDes(keyData, vKey)
		return clearKey, nil

	} else {
		return nil, xerrors.Errorf("gohsm: illegal key size - %s", len(keyData))
	}

}

//decrypt a key encrypted under the lmk using thales variants

func decryptKey(keyStr string, keyType string) ([]byte, error) {

	keyData := extractKeyData(keyStr)

	lmkKey := keyTypeTable[keyType[1:]]
	if lmkKey == nil {
		panic("unsupported key type" + keyType)
	}

	return decryptKeyThalesV(keyData, lmkKey, keyType)

}

//decrypt a key under the kek using thales variants

func decryptKeyThalesV(keyData []byte, lmkKey []byte, keyType string) ([]byte, error) {

	vKey := make([]byte, 16)
	copy(vKey, lmkKey)
	setVariants(vKey, keyType)

	if len(keyData) == 24 {
		//left
		vKey[8] = vKey[8] ^ variantTripleLen1
		left, err := crypto.DecryptTripleDes(keyData[0:8], vKey)
		if err != nil {
			return nil, err
		}
		//now middle
		vKey[8] = vKey[8] ^ variantTripleLen1
		vKey[8] = vKey[8] ^ variantTripleLen2
		middle, err := crypto.DecryptTripleDes(keyData[8:16], vKey)
		if err != nil {
			return nil, err
		}
		//right
		vKey[8] = vKey[8] ^ variantTripleLen2
		vKey[8] = vKey[8] ^ variantTripleLen3
		right, err := crypto.DecryptTripleDes(keyData[16:], vKey)
		if err != nil {
			return nil, err
		}
		clearKey := make([]byte, 0)
		clearKey = append(clearKey, left...)
		clearKey = append(clearKey, middle...)
		clearKey = append(clearKey, right...)
		return clearKey, nil

	} else if len(keyData) == 16 {
		vKey[8] = vKey[8] ^ variantDblLen1
		leftHalf, err := crypto.DecryptTripleDes(keyData[0:8], vKey)
		if err != nil {
			return nil, err
		}
		//now right half
		vKey[8] = vKey[8] ^ variantDblLen1
		vKey[8] = vKey[8] ^ variantDblLen2
		rightHalf, err := crypto.DecryptTripleDes(keyData[8:], vKey)
		if err != nil {
			return nil, err
		}
		clearKey := make([]byte, 0)
		clearKey = append(clearKey, leftHalf...)
		clearKey = append(clearKey, rightHalf...)

		return clearKey, nil

	} else if len(keyData) == 8 {
		//single length
		clearKey, err := crypto.DecryptTripleDes(keyData, vKey)
		if err != nil {
			return nil, err
		}
		return clearKey, nil

	} else {
		return nil, xerrors.Errorf("gohsm: illegal key size - %d", len(keyData))
	}

}

func setVariants(vKey []byte, keyType string) {

	if keyType[0] != '0' {
		//apply variant
		i, _ := strconv.ParseInt(keyType[:1], 10, 32)
		vKey[0] = vKey[0] ^ __variants[i]

	}
}
