package server

import (
	"crypto/md5"
	"crypto/sha1"
)

// bloomFilterSize 定义了布隆过滤器的大小，即位数组的长度。
const bloomFilterSize = 1000000

// numHashFunctions 定义了使用的哈希函数的数量。
const numHashFunctions = 7

// BloomFilter 是一个布隆过滤器结构体，包含一个位数组。
type BloomFilter struct {
	bitArray []bool
}

// NewBloomFilter 创建并返回一个新的布隆过滤器实例。
func NewBloomFilter() *BloomFilter {
	return &BloomFilter{
		bitArray: make([]bool, bloomFilterSize),
	}
}

// Add 方法将给定的URL添加到布隆过滤器中。
// 它通过计算多个哈希值并将相应的位设置为true来实现。
func (bf *BloomFilter) Add(url string) {
	hashes := bf.getHashIndices(url)
	for _, index := range hashes {
		bf.bitArray[index] = true
	}
}

// MightContain 方法检查布隆过滤器是否可能包含给定的URL。
// 如果所有计算出的哈希位都为true，则返回true，表示可能包含。
// 请注意，这个方法可能会产生假阳性，但不会产生假阴性。
func (bf *BloomFilter) MightContain(url string) bool {
	hashes := bf.getHashIndices(url)
	for _, index := range hashes {
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

// getHashIndices 方法根据给定的URL计算出多个哈希索引。
// 它使用MD5和SHA1哈希函数组合来生成这些索引，并确保它们在位数组的范围内。
func (bf *BloomFilter) getHashIndices(url string) []int {
	hashes := make([]int, numHashFunctions)
	hash1 := md5.Sum([]byte(url))
	hash2 := sha1.Sum([]byte(url))
	for i := 0; i < numHashFunctions; i++ {
		combinedHash := hash1[i%len(hash1)] + hash2[(i*2)%len(hash2)]
		index := int(combinedHash) % bloomFilterSize
		if index < 0 {
			index = -index
		}
		hashes[i] = index
	}
	return hashes
}
