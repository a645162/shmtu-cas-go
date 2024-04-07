package str

// SplitLengthN 将字符串每n个字符分割成一个切片
func SplitLengthN(s string, n int) []string {
	var parts []string
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		parts = append(parts, s[i:end])
	}
	return parts
}
