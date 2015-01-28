package keyblock

const (
	tag_thales = "S"
	tag_tr31   = "R"

	keyblock_version_3des = "0"
	keyblock_version_aes  = "1"

	//key usage
	dek_generic         = "D0"
	dek                 = "21"
	zek                 = "22"
	tek                 = "23"
	tmk                 = "51"
	zmk                 = "52"
	mac_iso_16609       = "M0"
	mac_iso_9791_1_alg1 = "M1"
	mac_iso_9791_1_alg2 = "M2"
	mac_iso_9791_1_alg3 = "M3"
	mac_iso_9791_1_alg4 = "M4"
	mac_aes_cbc         = "M5"
	mac_aes_cmac        = "M6"
	tpk                 = "71"
	zpk                 = "72"

	//algorithm
	alg_aes  = "A"
	alg_des  = "D"
	alg_3des = "T"
	alg_ecc  = "E"
	alg_hmac = "H"
	alg_rsa  = "R"
	alg_dsa  = "S"

	//modes
	mode_enc_dec     = "B"
	mode_mac_gen_ver = "C"
	mode_dec         = "D"
	mode_enc         = "E"
	mode_mac_gen     = "G"
	mode_all         = "N"
	mode_dsig_gen    = "S"
	mode_dsig_verify = "V"
)

type ThalesKeyBlockHeader struct {
	version_id       byte
	key_block_length []byte //4 bytes
	key_usage        string
	algo             string
	mode_of_use      string
	key_version      string
	exportability    string
	n_opt_hdr_blocks string
	lmk_id           string
}

type KeyDataBlock struct {
	key_length []byte //in bits
	key        []byte
	padding    []byte
}

type ThalesKeyBlock struct {
	tag                     string
	header                  ThalesKeyBlockHeader
	key_data_block          KeyDataBlock
	key_block_authenticator []byte //leftmost 4 bytes of mac (as ascii)
}
