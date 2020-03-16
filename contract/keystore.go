package contract

// KeyStore ...
type KeyStore struct {
	Address string `json:"address"`
	Crypto  Crypto `json:"crypto"`
	ID      string `json:"id"`
	Version int64  `json:"version"`
}

// Crypto ...
type Crypto struct {
	Cipher       string       `json:"cipher"`
	CipherText   string       `json:"ciphertext"`
	CipherParams CipherParams `json:"cipherparams"`
	Kdf          string       `json:"kdf"`
	Kdfparams    KdfParams    `json:"kdfparams"`
	MAC          string       `json:"mac"`
}

// CipherParams ...
type CipherParams struct {
	Iv string `json:"iv"`
}

// KdfParams ...
type KdfParams struct {
	Dklen int64  `json:"dklen"`
	N     int64  `json:"n"`
	P     int64  `json:"p"`
	R     int64  `json:"r"`
	Salt  string `json:"salt"`
}
