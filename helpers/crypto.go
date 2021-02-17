package helpers

import (
	"encoding/base64"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const time_diff_limit = 480

func tsDiff(ts string) bool {
	_ts, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return false
	}
	return math.Abs(float64(_ts-time.Now().Unix())) <= time_diff_limit
}

func doubleEncrypt(str string, cid string, sck string) string {
	arr := []byte(str)
	result := encrypt(arr, cid)
	result = encrypt(result, sck)
	return strings.Replace(strings.Replace(strings.TrimRight(base64.StdEncoding.EncodeToString(result), "="), "+", "-", -1), "/", "_", -1)
}

func encrypt(str []byte, k string) []byte {
	var result []byte
	strls := len(str)
	strlk := len(k)
	for i := 0; i < strls; i++ {
		char := str[i]
		keychar := k[(i+strlk-1)%strlk]
		char = byte((int(char) + int(keychar)) % 128)
		result = append(result, char)
	}
	return result
}

func doubleDecrypt(str string, cid string, sck string) string {
	if i := len(str) % 4; i != 0 {
		str += strings.Repeat("=", 4-i)
	}
	result, err := base64.StdEncoding.DecodeString(strings.Replace(strings.Replace(str, "-", "+", -1), "_", "/", -1))
	if err != nil {
		return ""
	}
	result = decrypt(result, cid)
	result = decrypt(result, sck)
	return string(result[:])
}

func decrypt(str []byte, k string) []byte {
	var result []byte
	strls := len(str)
	strlk := len(k)
	for i := 0; i < strls; i++ {
		char := str[i]
		keychar := k[(i+strlk-1)%strlk]
		char = byte(((int(char) - int(keychar)) + 256) % 128)
		result = append(result, char)
	}
	return result
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyHash(hashed, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	r := string(b)
	return r
}
