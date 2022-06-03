package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"github.com/btcsuite/btcd/btcec"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"
)

type EncryptedReqBody struct {
	Encrypted string `json:"encrypted"`
}

type EncryptedResBody struct {
	Encrypted string `json:"encrypted"`
	Iv        string `json:"iv"`
}

type ecPublicKey struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

//This type provides compatibility with the btcec package
type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

// ParseEcdsaPubKeyFromPem ...
func ParseEcdsaPubKeyFromPem(pemContent []byte) (*btcec.PublicKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, errors.New("invalid pem")
	}

	var ecp ecPublicKey
	_, err := asn1.Unmarshal(block.Bytes, &ecp)
	if err != nil {
		return nil, err
	}

	return btcec.ParsePubKey(ecp.PublicKey.RightAlign(), btcec.S256())
}

// ParseEcdsaPrivateKeyFromPem ...
func ParseEcdsaPrivateKeyFromPem(pemContent []byte) (*btcec.PrivateKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, errors.New("invalid pem")
	}

	var ecp ecPrivateKey
	_, err := asn1.Unmarshal(block.Bytes, &ecp)
	if err != nil {
		return nil, err
	}

	priKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), ecp.PrivateKey)
	return priKey, nil
}

func GenShareKey(priBCTSKey, pubGatewayKey []byte) (string, error) {
	pubKey, err := ParseEcdsaPubKeyFromPem(pubGatewayKey)
	if err != nil {
		return "", nil
	}

	priKey, err := ParseEcdsaPrivateKeyFromPem(priBCTSKey)
	if err != nil {
		return "", nil
	}

	aesKey := btcec.GenerateSharedSecret(priKey, pubKey)

	return base64.StdEncoding.EncodeToString(aesKey), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// AESEncryptStr base64加密字符串
func AESEncryptStr(src string, key, iv []byte) (encmess string, err error) {
	ciphertext, err := AESEncrypt([]byte(src), key, iv)
	if err != nil {
		return
	}

	encmess = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

// AESEncrypt 加密
func AESEncrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
	if len(iv) == 0 {
		iv = key[:16]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

func AESDecryptStr(src string, key, iv []byte) (string, error) {
	bsrc, err := base64.StdEncoding.DecodeString(src)
	bret, err := AESDecrypt(bsrc, key, iv)
	if err != nil {
		return "", err
	}
	return string(bret), nil
}

// AESDecrypt 解密
func AESDecrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
	if len(iv) == 0 {
		iv = key[:16]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(src, src)
	src = unPadding(src)
	return src, nil
}

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unPadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

func EncryptResBody(msg, aesKey string) ([]byte, []byte, error) {
	iv, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, nil, errors.New("iv generaton error")
	}

	mKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, nil, errors.New("invalid aes key")
	}

	encrypted, err := AESEncryptStr(msg, mKey, iv)
	if err != nil {
		return nil, nil, err
	}
	return []byte(encrypted), iv, nil
}

func DecryptReqBody(encrypted *EncryptedReqBody, aesKey, iv string) ([]byte, error) {
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil || len(aesIV) != 16 {
		return nil, errors.New("invalid aes iv")
	}

	mKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, errors.New("invalid aes key")
	}

	decrypted, err := AESDecryptStr(encrypted.Encrypted, mKey, aesIV)
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
}

func DecryptString(encrypted string, aesKey, iv string) ([]byte, error) {
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil || len(aesIV) != 16 {
		return nil, errors.New("invalid aes iv")
	}

	mKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, errors.New("invalid aes key")
	}

	decrypted, err := AESDecryptStr(encrypted, mKey, aesIV)
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
}

func EncryptedResponse(responseStr, aesKey string) *EncryptedResBody {
	msg, ivNew, err := EncryptResBody(responseStr, aesKey)
	if err != nil {
		logx.Errorf("crypto EncryptedResponse err:%+v", err)
		panic(err)
	}

	ivBase64 := base64.StdEncoding.EncodeToString(ivNew)

	return &EncryptedResBody{
		Encrypted: string(msg),
		Iv:        ivBase64,
	}
}

func ExportUserIDFromHeader(r *http.Request, aesKey string) (int64, error) {
	logx.Infof("[CRYPTO Header:%+v]", r.Header)
	iv := r.Header.Get("X-Encrypt-Iv")
	userId := r.Header.Get("X-User-Id")

	userIDAfterDecry, err := DecryptString(userId, aesKey, iv)
	if err != nil {
		return 0, err
	}

	userIdInt, _ := strconv.ParseInt(string(userIDAfterDecry), 10, 64)

	return userIdInt, nil
}
