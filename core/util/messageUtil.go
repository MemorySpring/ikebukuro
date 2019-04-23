package util

import (
	"bytes"
	"crypto/md5"
	"errors"
	"math/rand"
	"net"
	"time"
)

const (
	MAX_INT = 256*256*256*256 - 1
)

func WrapControlMessage(control byte, args ...interface{}) ([]byte, error) {
	data := []byte{control}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			b, err := StringTo32Bytes(arg.(string))
			if err != nil {
				return data, err
			}
			data = MergeBytes(data, b)
		case int:

		}
	}
	return data, nil
}

func StringTo32Bytes(str string) ([]byte, error) {
	b := []byte(str)
	if len(b) > 32 {
		return b, errors.New("字符串过长")
	}
	data := make([]byte, 32)
	copy(data, b)
	return data, nil
}

func GetMD5(str string) [16]byte {
	data := []byte(str)
	hash := md5.Sum(data)
	return hash
}

func MergeBytes(bs ...[]byte) []byte {
	return bytes.Join(bs, []byte(""))
}

func GetString(conn net.Conn) (string, error) {
	data := make([]byte, 32)
	_, err := conn.Read(data)
	return string(GetValidByte(data)), err
}

func GetInt(conn net.Conn) (int, error) {
	data := make([]byte, 4)
	_, err := conn.Read(data)
	return ByteToInt(data), err
}

func ByteToInt(bs []byte) int {
	i := 0
	base := 1
	for _, b := range bs {
		i += int(b) * base
		base *= 256
	}
	return i
}

func IntTo4Bytes(i int) ([]byte, error) {
	data := make([]byte, 4)
	if i > MAX_INT {
		return data, errors.New("超过最大值")
	}
	base := 256 * 256 * 256
	data[3] = byte(i / base)
	i = i % base
	base /= 256
	data[2] = byte(i / base)
	i = i % base
	base /= 256
	data[1] = byte(i / base)
	i = i % base
	base /= 256
	data[0] = byte(i / base)
	return data, nil
}

func GetValidByte(src []byte) []byte {
	var str_buf []byte
	for _, v := range src {
		if v != 0 {
			str_buf = append(str_buf, v)
		} else {
			break
		}
	}
	return str_buf
}
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
