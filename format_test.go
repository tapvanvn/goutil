package goutil_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/goutil"
)

func TestEscapeMarkdown(t *testing.T) {

	str := "*abcd_ef\\gh"
	str2 := goutil.MarkdownEscape(str)
	fmt.Println(str2)
}
