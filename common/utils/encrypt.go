package utils

import (
	"bytes"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Encrypt(pwd string) (string, error) {
	password := []byte(pwd)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashAndPwd(hash, pwd string) bool {
	password := []byte(pwd)
	hashBytes := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashBytes, password)
	if err != nil {
		return false
	}
	return true
}

func GetRandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(allChars[rand.Int63()%int64(len(allChars))])
	}
	return s.String()
}
