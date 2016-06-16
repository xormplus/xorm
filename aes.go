package xorm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const (
	tempkey = "1234567890!@#$%^&*()_+-="
)

type AesEncrypt struct {
	PubKey string
}

func (this *AesEncrypt) getKey() []byte {
	strKey := this.PubKey

	keyLen := len(strKey)

	if keyLen < 16 {
		rs := []rune(tempkey)
		strKey = strKey + string(rs[0:16-keyLen])
	}

	if keyLen > 16 && keyLen < 24 {
		rs := []rune(tempkey)
		strKey = strKey + string(rs[0:24-keyLen])
	}

	if keyLen > 24 && keyLen < 32 {
		rs := []rune(tempkey)
		strKey = strKey + string(rs[0:32-keyLen])
	}

	arrKey := []byte(strKey)
	if keyLen >= 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[:24]
	}

	return arrKey[:16]
}

//加密字符串
func (this *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (decrypted []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err = base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return nil, err
	}
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted = make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return decrypted, nil
}
