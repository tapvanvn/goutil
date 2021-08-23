package goutil

import "regexp"

var re = regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")

func TripJSONComment(jsonc []byte) []byte {
	return re.ReplaceAll(jsonc, nil)
}
