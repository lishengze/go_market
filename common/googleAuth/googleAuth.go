package googleAuth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strconv"
	"time"
)

func prefix0(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := 6 - len(otp); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

// GenerateTwoFactorCode  根据google auth secret 生成 6位验证码
func GenerateTwoFactorCode(keyStr string) string {
	// sign the value using HMAC-SHA1
	epochSeconds := time.Now().Unix()
	key, _ := base32.StdEncoding.DecodeString(keyStr)
	value := toBytes(epochSeconds / 30)
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	otp := strconv.Itoa(int(pwd))

	return prefix0(otp)
}

// GenerateSecret 生成一个随机 Google auth secret
func GenerateSecret(username string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	data := []byte(now + username)
	shaRes := fmt.Sprintf("%x", sha256.Sum256(data))[0:10]
	fmt.Println(shaRes)
	res := base32.StdEncoding.EncodeToString([]byte(shaRes))
	return res
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

///https://github.com/skip2/go-qrcode
