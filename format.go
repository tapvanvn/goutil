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
