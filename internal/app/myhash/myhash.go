package myhash

import (
	"math/rand"
)

func Base62() string {
	s := "012345689abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	num := rand.Int63()
	if num == 0 {
		num++
	}
	var hashStr string
	for num > 0 {
		hashStr = string(s[num%62]) + hashStr
		num /= 62
	}
	return hashStr
}
