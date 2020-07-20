package support

//TruncateString Actually shorten strings
func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

//AlphaCharOnly A-z
func AlphaCharOnly(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'A' && c < 'z' {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

//StripCtlAndExtFromBytes Strip all specials
func StripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 255 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}
