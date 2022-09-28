package util

import "testing"

func TestEncodeMD5(t *testing.T) {
	res := EncodeMD5("for test")
	if res != "a733a5c2c674b0728d3777f9443d9248" {
		t.Errorf("TestEncodeMD5 failed, res:%s", res)
	}
}
