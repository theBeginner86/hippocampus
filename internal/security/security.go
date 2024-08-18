package security

type Security struct {
	Encrypter *Encrypter
	Decrypter *Decrypter
}

func NewSecurity(privateKey string, bytes []byte) *Security {
	return &Security{
		Encrypter: NewEncrypter(privateKey, bytes),
		Decrypter: NewDecrypter(privateKey, bytes),
	}
}

