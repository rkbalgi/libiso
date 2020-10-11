package crypto

import "log"

var englishAlpha = []uint8{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

// RotN function rotates the input string by n (Caesars Cipher)
func RotN(n int, val string) string {

	rotNVal := make([]byte, len(val))
	n = n % 26

	for i := 0; i < len(val); i++ {
		tmp := val[i]
		found := false

		for j := 0; j < len(englishAlpha); j++ {
			if tmp == englishAlpha[j] {
				found = true
				if (j + n) < len(englishAlpha) {
					rotNVal[i] = englishAlpha[j+n]
				} else {
					//how many beyond the length of englishAlpha?
					iters := (j + n) - len(englishAlpha)
					rotNVal[i] = englishAlpha[iters]
				}
				break

			}

		}
		if !found {
			log.Printf(string(tmp) + " character found, only lowercase english alphabets expected.")
		}

	}

	return string(rotNVal[:])
}
