package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var secretKey string

func init() {
	secretKey = os.Getenv("TOKEN_SECRET")
	if secretKey == "" {
		key := make([]byte, 32)
		rand.Read(key)
		secretKey = base64.StdEncoding.EncodeToString(key)
	}
}

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TokenPayload Token 载荷结构
type TokenPayload struct {
	Type string `json:"type"`
	Exp  int64  `json:"exp"`
}

// GenerateToken 生成 Token
func GenerateToken(tokenType string, expireHours int) string {
	payload := TokenPayload{
		Type: tokenType,
		Exp:  time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadBase64 := base64.URLEncoding.EncodeToString(payloadJSON)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payloadBase64))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return payloadBase64 + "." + signature
}

// ValidateToken 验证 Token
func ValidateToken(token string, expectedType string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	// 验证签名
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(parts[0]))
	expectedSig := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if parts[1] != expectedSig {
		return false
	}

	// 解析 Payload
	payloadJSON, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	var payload TokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return false
	}

	// 验证类型
	if payload.Type != expectedType {
		return false
	}

	// 验证过期时间
	if time.Now().Unix() > payload.Exp {
		return false
	}

	return true
}
