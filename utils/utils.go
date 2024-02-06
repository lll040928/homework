package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"homework/api/global"
)

// Encrypt 加密
func Encrypt(plainText string) (string, error) {
	key := global.C.CryptInfo.Key
	iv := global.C.CryptInfo.Iv
	data, err := aesCBCEncrypt([]byte(plainText), []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt 解密
func Decrypt(cipherText string) (string, error) {
	key := global.C.CryptInfo.Key
	iv := global.C.CryptInfo.Iv
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	dnData, err := aesCBCDecrypt(data, []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}

	return string(dnData), nil
}

// aesCBCEncrypt AES/CBC/PKCS7Padding 加密
func aesCBCEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	// AES
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// PKCS7 填充
	plaintext = paddingPKCS7(plaintext, aes.BlockSize)

	// CBC 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(plaintext, plaintext)

	return plaintext, nil
}

// aesCBCDecrypt AES/CBC/PKCS7Padding 解密
func aesCBCDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	// AES
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	// CBC 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// PKCS7 反填充
	result := unPaddingPKCS7(ciphertext)
	return result, nil
}

// PKCS7 填充
func paddingPKCS7(plaintext []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plaintext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(plaintext, paddingText...)
}

// PKCS7 反填充
func unPaddingPKCS7(s []byte) []byte {
	length := len(s)
	if length == 0 {
		return s
	}
	unPadding := int(s[length-1])
	return s[:(length - unPadding)]
}
