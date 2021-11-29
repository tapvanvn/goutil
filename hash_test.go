package goutil_test

import (
	"testing"

	"github.com/tapvanvn/goutil"
)

func TestMD5Message(t *testing.T) {

	if goutil.MD5Message("hello") != "68656c6c6fd41d8cd98f00b204e9800998ecf8427e" {
		t.Error()
	}
}
