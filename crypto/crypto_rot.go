package crypto


var english_alpha = []uint8{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}


/**
 This function rotates the input string by n
 (Caesars Cipher)
**/
func RotN(n int, val string) string {

	rot_n_val := make([]byte, len(val))
	n = n % 26

	for i := 0; i < len(val); i++ {
		tmp := val[i]
		found := false

		for j := 0; j < len(english_alpha); j++ {
			if tmp == english_alpha[j] {
				found = true
				if (j + n) < len(english_alpha) {
					rot_n_val[i] = english_alpha[j+n]
				} else {
					//how many beyond the length of english_alpha?
					iters := (j + n) - len(english_alpha)
					rot_n_val[i] = english_alpha[iters]
				}
				break

			}

		}
		if !found {
			panic(string(tmp) + " character found, only lowercase english alphabets expected.")
		}

	}

	return (string(rot_n_val[:]))
}
