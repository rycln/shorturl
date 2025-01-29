package hash

func Base62(num int64) string {

	s := "012345689abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var hashStr string
	for num > 0 {
		hashStr = string(s[num%62]) + hashStr
		num /= 62
	}
	return hashStr
}
