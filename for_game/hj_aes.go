/*
描述 :  golang  AES/ECB/PKCS5  加密解密
date : 2016-04-08
*/

package for_game

import (
	"crypto/aes"
	"encoding/base64"
	"game_server/easygo"
)

func HJAesDecrypt(src string, key []byte) string {
	crypted, err := base64.StdEncoding.DecodeString(src)
	easygo.PanicError(err)
	block, err := aes.NewCipher(key)
	easygo.PanicError(err)
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}

func HJAesEncrypt(src, key string) string {
	if src == "" {
		return ""
	}
	block, err := aes.NewCipher([]byte(key))
	easygo.PanicError(err)
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted)
}
