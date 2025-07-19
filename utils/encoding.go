package utils

import (
	"bytes"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// 常见的编码类型
var encodings = []struct {
	name string
	enc  encoding.Encoding
}{
	{"GBK", simplifiedchinese.GBK},
	{"GB18030", simplifiedchinese.GB18030},
	{"Big5", traditionalchinese.Big5},
	{"Shift_JIS", japanese.ShiftJIS},
	{"EUC-JP", japanese.EUCJP},
	{"EUC-KR", korean.EUCKR},
}

// DecodeString 尝试使用多种编码解码字符串
func DecodeString(data []byte) string {
	if utf8.Valid(data) {
		return string(data)
	}

	for _, enc := range encodings {
		if decoded, err := decodeWithEncoding(data, enc.enc); err == nil {
			return decoded
		}
	}

	return string(data)
}

// decodeWithEncoding 使用指定编码解码
func decodeWithEncoding(data []byte, enc encoding.Encoding) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// DecodeFileName 专门用于解码文件名
func DecodeFileName(name string) string {
	// 如果文件名看起来已经是正确的，直接返回
	if !containsInvalidUTF8(name) {
		return name
	}

	// 尝试解码
	decoded := DecodeString([]byte(name))

	return decoded
}

// containsInvalidUTF8 检查字符串是否包含无效的 UTF-8 字符
func containsInvalidUTF8(s string) bool {
	for _, r := range s {
		if r == 0xFFFD { // Unicode replacement character
			return true
		}
	}
	return false
}
