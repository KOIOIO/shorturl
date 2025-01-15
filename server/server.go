package server

import (
	"crypto/md5"
	"encoding/base64"
)

var Bloom = NewBloomFilter()

func GenerateShortURLString(url string) string {
	hash := md5.Sum([]byte(url))
	shortURLBytes := hash[:]
	encoded := base64.URLEncoding.EncodeToString(shortURLBytes)
	return encoded[:8]
}
