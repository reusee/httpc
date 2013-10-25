package httpc

import (
	"fmt"
	"github.com/reusee/goquery"
	"testing"
)

func TestHttpc(t *testing.T) {
	client, err := NewClient(&Client{
		Encoding: "gb2312",
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
	res, err := client.GetFind("http://qq.com", "a")
	if err != nil {
		t.Fatalf("%v", err)
	}
	res.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		fmt.Printf("%s %s\n", s.Text(), href)
	})
}

func TestGetBytes(t *testing.T) {
	client, err := NewClient(&Client{
		Encoding: "gb2312",
	})
	if err != nil {
		t.Fatal(err)
	}
	content, err := client.GetBytes("http://qq.com")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("%d %s\n", len(content), content[1024:2048])
}
