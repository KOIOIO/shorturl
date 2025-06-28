package repository

// base62Chars 定义了Base62编码使用的字符集
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// base62Encode 将一个整数转换为Base62编码的字符串
// 这个函数用于生成短URL的字符串表示
func Base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte
	for num > 0 {
		remainder := num % 62
		num /= 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
	}
	return string(encoded)
}
