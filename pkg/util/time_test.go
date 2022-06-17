package util

import "testing"

func TestGetCurrentStr(t *testing.T) {
	res := GetCurrentStr()
	if res == "" {
		t.Errorf("TestEncodeMD5 failed, res:%s", res)
	}
}
