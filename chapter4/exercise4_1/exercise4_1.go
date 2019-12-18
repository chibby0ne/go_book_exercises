package main

import (
	"crypto/sha256"
	"fmt"
)

func getNumberSetBits(b byte) uint64 {
	var count uint64
	for ; b != 0; count++ {
		b = b & (b - 1)
	}
	return count
}

func getNumDifferentBits(h1, h2 [32]byte) uint64 {
	var diffBits uint64
	for i := range h1 {
		diffBits += getNumberSetBits(h1[i] ^ h2[i])
	}
	return diffBits
}

func main() {
	hash1 := sha256.Sum256([]byte("x"))
	hash2 := sha256.Sum256([]byte("X"))
	diffBits := getNumDifferentBits(hash1, hash2)
	fmt.Printf("hash1: %x\nhash2: %x\ndiff bits: %v\n", hash1, hash2, diffBits)
}
