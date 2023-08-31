package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 加密
func Encrypt(data string, key []byte) string {
	fillKey := FillKey(key)
	encryptCBC := AesEncryptCBC([]byte(data), fillKey)   // AES CBC加密
	return base64.URLEncoding.EncodeToString(encryptCBC) // 在通过base64加密
}

// 解密
func Decryption(data string, key []byte) string {
	fillKey := FillKey(key)
	decodeString, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	decryptCBC := AesDecryptCBC(decodeString, fillKey)
	return string(decryptCBC)
}

// ===================AES  CBC  模式======================
// Aes 加密
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	block, _ := aes.NewCipher(key)                              // NewCipher该函数限制了输入k的长度必须为16, 24或者32
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

// Aes 解密
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length-unpadding <= 0 {
		return []byte("")
	}
	return origData[:(length - unpadding)]
}

// 检查加密key的长度,当长度不够时,使用特定字符串来填充
func FillKey(key []byte) []byte {
	k := len(key)
	// 用于填充的字符串
	fillStr := "cPaw45cBxur1I9OwudEf7PYviFAjaMWG"
	if k < 16 {
		/* 填充到16位*/
		fillLen := 16 - k
		useStr := fillStr[:fillLen]
		useByte := []byte(useStr)
		key = append(key, useByte...)
		//common.Logger.Infof("16位填充之后的字符串 %s", string(key))
		return key
	} else if k > 16 && k < 24 {
		/* 填充到24位*/
		fillLen := 24 - k
		useStr := fillStr[:fillLen]
		useByte := []byte(useStr)
		key = append(key, useByte...)
		//common.Logger.Infof("24位填充之后的字符串 %s", string(key))
		return key
	} else if k > 24 && k < 32 {
		/* 填充到32位*/
		fillLen := 32 - k
		useStr := fillStr[:fillLen]
		useByte := []byte(useStr)
		key = append(key, useByte...)
		//common.Logger.Infof("32位填充之后的字符串 %s", string(key))
		return key
	} else if k > 32 {
		/* 切割到32位*/
		key = key[:32]
		//common.Logger.Infof("超长之后切割字符串 %s", string(key))
		return key
	} else {
		/*不出来返回*/
		return key
	}
}
