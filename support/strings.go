package support

//TruncateString Actually shorten strings
func TruncateString(str string, num int) (string, bool) {
	output := str

	if len(str) > num {
		output = str[0:num]
		return output, true
	}
	return output, false
}

//AlphaOnly A-z
func AlphaOnly(str string) string {
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

//All characters except A-z and control
func NonAlpha(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if (c >= ' ' && c < 'a') || (c > 'z' && c < 255) {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

//0-9 only
func NumericOnly(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= '0' && c <= '9' {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

//No ASCII control characters
func StripControl(str string) string {
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
