/*
google authenticator
*/
package for_game

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (this *GoogleAuth) un() int64 {
	return time.Now().Unix() / 30
}

func (this *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (this *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (this *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (this *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (this *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (this *GoogleAuth) oneTimePassword(key []byte, value []byte) uint32 {
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := this.toUint32(hashParts)
	pwd := number % 1000000

	return pwd
}

func (this *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, this.un())
	return strings.ToUpper(this.base32encode(this.hmacSha1(buf.Bytes(), nil)))
}

func (this *GoogleAuth) GetCode(secret string) (string, error) {
	inputNoSpaces := strings.Replace(secret, " ", "", -1)
	secretUpper := strings.ToUpper(inputNoSpaces)
	secretKey, err := this.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := this.oneTimePassword(secretKey, this.toBytes(this.un()))
	return fmt.Sprintf("%06d", number), nil
}

func (this *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

func (this *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := this.GetQrcode(user, secret)
	return fmt.Sprintf("http://www.google.com/chart?chs=200x200&chld=M%%7C0&cht=qr&chl=%s", qrcode)
}

func (this *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := this.GetCode(secret)
	fmt.Println(_code, code, err)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
