package httpc

import (
	"fmt"
	"testing"
)

func TestFromGbk(t *testing.T) {
	runes, err := From("gbk", []byte{
		'g', 'b', 'k', 0xb1, 0xe0, 0xc2, 0xeb, 0xb5,
		0xc4, 0xd7, 0xd6, 0xb7, 0xfb, 0xb4, 0xae,
	})
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%s\n", string(runes))
	if string(runes) != "gbk编码的字符串" {
		t.Fail()
	}
}
