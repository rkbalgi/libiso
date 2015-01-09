package crypto

import (
	//"crypto/rand"
	"encoding/hex"
)

//596b0abc958120f62a7edc6896c99144

func GenerateDesKey(key_len int) []byte {
	
	if(key_len%8!=0){
		panic("invalid keylen for generation")
	}
	
	key := make([]byte, key_len)
	tmp,_:=hex.DecodeString("596b0abc958120f62a7edc6896c991442a7edc6896c99144");
	if(key_len==8){
		return tmp[0:8] 
	}else if(key_len==16){
		return tmp[:16]
	}else{
		return tmp[:24]
	}
	
	return key

}


/*
func GenerateDesKey(key_len int) []byte {
	
	if(key_len%8!=0){
		panic("invalid keylen for generation")
	}
	
	key := make([]byte, key_len)
	n, err := rand.Read(key)
	if n != key_len || err != nil {
		panic("key gen failure" + err.Error())
	}

	return key

}
*/