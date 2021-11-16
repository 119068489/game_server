package for_game

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"game_server/easygo"
	"math/big"
)

// copy from crypt/rsa/pkcs1v5.go
var hashPrefixes = map[crypto.Hash][]byte{
	crypto.MD5:       {0x30, 0x20, 0x30, 0x0c, 0x06, 0x08, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x02, 0x05, 0x05, 0x00, 0x04, 0x10},
	crypto.SHA1:      {0x30, 0x21, 0x30, 0x09, 0x06, 0x05, 0x2b, 0x0e, 0x03, 0x02, 0x1a, 0x05, 0x00, 0x04, 0x14},
	crypto.SHA224:    {0x30, 0x2d, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x04, 0x05, 0x00, 0x04, 0x1c},
	crypto.SHA256:    {0x30, 0x31, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01, 0x05, 0x00, 0x04, 0x20},
	crypto.SHA384:    {0x30, 0x41, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x02, 0x05, 0x00, 0x04, 0x30},
	crypto.SHA512:    {0x30, 0x51, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x03, 0x05, 0x00, 0x04, 0x40},
	crypto.MD5SHA1:   {}, // A special TLS case which doesn't use an ASN1 prefix.
	crypto.RIPEMD160: {0x30, 0x20, 0x30, 0x08, 0x06, 0x06, 0x28, 0xcf, 0x06, 0x03, 0x00, 0x31, 0x04, 0x14},
}

// copy from crypt/rsa/pkcs1v5.go
func encrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

// copy from crypt/rsa/pkcs1v5.go
func pkcs1v15HashInfo(hash crypto.Hash, inLen int) (hashLen int, prefix []byte, err error) {
	// Special case: crypto.Hash(0) is used to indicate that the data is
	// signed directly.
	if hash == 0 {
		return inLen, nil, nil
	}

	hashLen = hash.Size()
	if inLen != hashLen {
		return 0, nil, errors.New("crypto/rsa: input must be hashed message")
	}
	prefix, ok := hashPrefixes[hash]
	if !ok {
		return 0, nil, errors.New("crypto/rsa: unsupported hash function")
	}
	return
}

// copy from crypt/rsa/pkcs1v5.go
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}
func unLeftPad(input []byte) (out []byte) {
	n := len(input)
	t := 2
	for i := 2; i < n; i++ {
		if input[i] == 0xff {
			t = t + 1
		} else {
			if input[i] == input[0] {
				t = t + int(input[1])
			}
			break
		}
	}
	out = make([]byte, n-t)
	copy(out, input[t:])
	return
}

// copy&modified from crypt/rsa/pkcs1v5.go
func publicDecrypt(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte) (out []byte, err error) {
	hashLen, prefix, err := pkcs1v15HashInfo(hash, len(hashed))
	if err != nil {
		return nil, err
	}

	tLen := len(prefix) + hashLen
	k := (pub.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, fmt.Errorf("length illegal")
	}

	c := new(big.Int).SetBytes(sig)
	m := encrypt(new(big.Int), pub, c)
	em := leftPad(m.Bytes(), k)
	out = unLeftPad(em)

	err = nil
	return
}
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

//私钥加密
func PrivateEncrypt(privateKey []byte, data []byte) ([]byte, error) {
	//解密pem格式的私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	// 解析私钥
	privt, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//signData, err := rsa.SignPKCS1v15(nil, privt, crypto.Hash(0), data)
	if err != nil {
		return nil, err
	}
	partLen := privt.N.BitLen()/8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.SignPKCS1v15(rand.Reader, privt, crypto.Hash(0), chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

//公钥解密
func PublicDecrypt(publicKey []byte, data []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	easygo.PanicError(err)
	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("强制转换失败")
	}
	partLen := p.N.BitLen() / 8
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decData, err := publicDecrypt(p, crypto.Hash(0), nil, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decData)
	}
	return buffer.Bytes(), nil
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key decode error")
	}
	pkcs1PrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("private key parse error")
	}
	return pkcs1PrivateKey, nil
}

func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pkixPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pkixPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key type error")
	}
	return pub, nil
}

func PrivateKeySignAndBase64(privateKey []byte, data []byte) (string, error) {

	pkcs1PrivateKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)

	signPKCS1v15, err := rsa.SignPKCS1v15(nil, pkcs1PrivateKey, crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}
	base64EncodingData := base64.StdEncoding.EncodeToString(signPKCS1v15)
	return base64EncodingData, nil
}

func Base64DecodeAndPublicKeyVerifySign(publicKey []byte, data []byte, sign string) error {
	decodeSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	key, err := parsePublicKey(publicKey)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(key, crypto.SHA1, hashed, decodeSign)
	if err != nil {
		return err
	}
	return nil
}

func PublicKeyEncryptAndBase64(src []byte, publicKey []byte) (string, error) {
	key, err := parsePublicKey(publicKey)
	if err != nil {
		return "", err
	}
	encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, key, src)
	if err != nil {
		return "", err
	}
	base64EncodingData := base64.StdEncoding.EncodeToString(encryptPKCS1v15)
	return base64EncodingData, nil
}

func Base64DecodeAndPrivateKeyDecrypt(cipherText string, privateKey []byte) ([]byte, error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return []byte{}, err
	}
	pkcs1PrivateKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return []byte{}, err
	}
	decrypt, err := rsa.DecryptPKCS1v15(rand.Reader, pkcs1PrivateKey, decodeCipher)
	if err != nil {
		return []byte{}, err
	}
	return decrypt, nil
}
