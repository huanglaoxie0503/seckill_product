package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 16 24 32 位字符，分别对应AES-128,AES-192 AES-256 加密方法
var PwdKey = []byte("DIS**#AAAGGGDDD")

// PKCS7 填充模式
func PKCS7Padding(cipher []byte, blockSize int) []byte {
	padding := blockSize - len(cipher)%blockSize
	// Repeat() :函数功能是把切片[]byte{byte(padding)}复制 padding 个，然后合并成新的字节切片返回
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipher, padText...)
}

// 填充反向操作，删除填充字符
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	// 获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		// 获取填充字符串长度
		unPadding := int(origData[length-1])
		return origData[:(length - unPadding)], nil
	}
}

// 实现AES加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// blockSize 大小
	blockSize := block.BlockSize()
	// 对数据进行填充，让数据长度满足要求
	origData = PKCS7Padding(origData, blockSize)
	// 采用AES中的CBC加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	cryptoSlice := make([]byte, len(origData))
	// 执行加密
	blockMode.CryptBlocks(cryptoSlice, origData)
	return cryptoSlice, nil
}

// 实现AES解密
func AesDeEncrypt(cryptoSlice []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块大小
	blockSize := block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryptoSlice))
	// 执行函数
	blockMode.CryptBlocks(origData, cryptoSlice)
	// 去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// base64 加密 切片转 string
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// base64字符串解密
func DePwdCode(pwd string) ([]byte, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesEncrypt(pwdByte, PwdKey)
}
