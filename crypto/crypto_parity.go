package crypto

func to_odd_parity(b []byte) {
	for i,_ := range b {


		_c := 0
		var j uint = 0
		for ; j < 9; j++ {
			if b[i]>>j&0x01 == 0x01 {
				_c++
			}
		}
		if _c%2 == 0 {

			if b[i]&0x01 == 0x01 {
				//last bit is set
				b[i] = b[i] & 0xfe
			} else {
				b[i] = b[i] | 0x01
			}

		}

	}
}
