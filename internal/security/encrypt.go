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

type Encrypter struct {
  privateKey string
  bytes []byte
}

func NewEncrypter(privateKey string, byts []byte) *Encrypter {
  return &Encrypter{
    privateKey: privateKey,
    bytes: byts,
  }
}

func (encrypt *Encrypter) Encode(b []byte) string {
  return base64.StdEncoding.EncodeToString(b)
}

func (encrypt *Encrypter) Encrypt(text string) (string, error) {
  block, err := aes.NewCipher([]byte(encrypt.privateKey))
  if err != nil {
    return "", err
  }
  plainText := []byte(text)
  cfb := cipher.NewCFBEncrypter(block, encrypt.bytes)
  cipherText := make([]byte, len(plainText))
  cfb.XORKeyStream(cipherText, plainText)
  return encrypt.Encode(cipherText), nil
}

// func Encrypter() {
//   StringToEncrypt := "Encrypting this string"

//   // To encrypt the StringToEncrypt
//   encText, err := Encrypt(StringToEncrypt, MySecret)
//   if err != nil {
//     fmt.Println("error encrypting your classified text: ", err)
//   }
//   fmt.Println(encText)
// }
