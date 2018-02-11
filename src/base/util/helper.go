package util

import (
    "crypto/rand"
)

const (
    letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

//随机产生一个字符串，长度自定义
func GenerateRandomString(n int) (string, error) {
    bytes, err := generateRandomBytes(n)
    if err != nil {
        return "", err
    }
    length := byte(len(letters))
    for i, b := range bytes {
        bytes[i] = letters[b%length]
    }
    return string(bytes), nil
}

//=======================================================
func generateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}