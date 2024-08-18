//  Copyright 2024 Pranav Singh

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Decrypter struct {
	privateKey string
	bytes []byte
}

func NewDecrypter(privateKey string, byts []byte) *Decrypter {
	return &Decrypter{
		privateKey: privateKey,
		bytes: byts,
	}
}

func (decrypt *Decrypter) Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
} 

func (decrypt *Decrypter) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(decrypt.privateKey))
	if err != nil {
		return "", err
	}
	cipherText, err := decrypt.Decode(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, decrypt.bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
