// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	//	"log"
)

type DesEncrypt struct {
	PubKey string
}

type TripleDesEncrypt struct {
	PubKey string
}

func (this *DesEncrypt) getKey() []byte {
	strKey := this.PubKey
	keyLen := len(strKey)

	if keyLen < 8 {
		rs := []rune(tempkey)
		strKey = strKey + string(rs[0:8-keyLen])
	}
	arrKey := []byte(strKey)
	return arrKey[:8]
}

func (this *TripleDesEncrypt) getKey() []byte {
	strKey := this.PubKey
	keyLen := len(strKey)

	if keyLen < 24 {
		rs := []rune(tempkey)
		strKey = strKey + string(rs[0:24-keyLen])
	}
	arrKey := []byte(strKey)
	return arrKey[:24]
}

func (this *DesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	key := this.getKey()
	origData := []byte(strMesg)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (this *DesEncrypt) Decrypt(crypted []byte) (decrypted []byte, err error) {
	key := this.getKey()

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	crypted, err = base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	decrypted = make([]byte, len(crypted))
	blockMode.CryptBlocks(decrypted, crypted)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted, nil
}

// 3DES加密
func (this *TripleDesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	key := this.getKey()
	origData := []byte(strMesg)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil

}

// 3DES解密
func (this *TripleDesEncrypt) Decrypt(crypted []byte) ([]byte, error) {
	key := this.getKey()
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	crypted, err = base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil

}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
