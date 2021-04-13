package commons

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	mathrand "math/rand"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltMinLen = 8
	saltMaxLen = 48
	iter       = 1000
	keyLen     = 48
)

// EncryptPbkdf2 加密密码
func EncryptPbkdf2(str string) (string, error) {
	// 1、生成随机长度的盐值
	salt, err := randSalt()
	if err != nil {
		return "", err
	}

	// 2、生成加密串
	en := encryptPwdWithSalt([]byte(str), salt)
	en = append(en, salt...)

	// 3、合并盐值
	encryptStr := base64.StdEncoding.EncodeToString(en)

	return encryptStr, nil
}

// CheckEncryptPbkdf2 验证字符串是否与加密串匹配
func CheckEncryptPbkdf2(str, encryptStr string) bool {
	// 1、参数校验
	if len(encryptStr) == 0 {
		return false
	}

	enDecode, err := base64.StdEncoding.DecodeString(encryptStr)
	if err != nil {
		return false
	}

	// 2、截取加密串 固定长度
	salt := enDecode[keyLen:]

	// 3、比对
	enBase64 := base64.StdEncoding.EncodeToString(enDecode[0:keyLen])
	pwdEnBase64 := base64.StdEncoding.EncodeToString(encryptPwdWithSalt([]byte(str), salt))

	return enBase64 == pwdEnBase64
}

func randSalt() ([]byte, error) {
	// 生成8-32之间的随机数字
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptPwdWithSalt(pwd, salt []byte) []byte {
	pwd = append(pwd, salt...)
	return pbkdf2.Key(pwd, salt, iter, keyLen, sha256.New)
}
