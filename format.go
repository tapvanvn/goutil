package goutil

import "regexp"

var __re = regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")

func TripJSONComment(jsonc []byte) []byte {
	return __re.ReplaceAll(jsonc, nil)
}

var __markdown_special = []rune{'\\', '`', '*', '_', '{', '}', '[', ']', '(', ')', '#', '+', '-', '.', '!'}

func MarkdownEscape(str string) string {
	rs := []rune{}
	for _, r := range str {
		found := false
		for _, m := range __markdown_special {
			if r == m {
				found = true
				break
			}
		}
		if found {
			rs = append(rs, '\\')
		}
		rs = append(rs, r)
	}
	return string(rs)
}
func Claim(value int, min int, max int) int {

	if value < min {

		return min
	}
	if value > max {

		return max
	}
	return value
}
func Ratio(value int, min int, max int) float32 {
	if min == max {
		panic("min cannot equal to max")
	}

	return float32(Claim(value, min, max)) / float32((max - min))
}

func ClaimFloat32(value float32, min float32, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func RatioFloat32(value float32, min float32, max float32) float32 {
	if min == max {
		panic("min cannot equal to max")
	}

	return ClaimFloat32(value, min, max) / (max - min)
}
