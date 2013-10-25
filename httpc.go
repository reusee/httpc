package httpc

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/proxy"
	"github.com/reusee/goquery"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	Socks5Proxy string
	Encoding    string
	DialTimeout time.Duration
	ReadTimeout time.Duration
	Retry       int

	Client *http.Client
}

func NewClient(client *Client) (*Client, error) {
	if client == nil {
		client = new(Client)
	}
	if client.DialTimeout == 0 {
		client.DialTimeout = time.Second * 10
	}
	if client.ReadTimeout == 0 {
		client.ReadTimeout = time.Second * 10
	}
	if client.Encoding == "" {
		client.Encoding = "utf-8"
	}
	if client.Retry == 0 {
		client.Retry = 3
	}
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, client.DialTimeout)
		},
		Proxy: http.ProxyFromEnvironment,
		ResponseHeaderTimeout: client.ReadTimeout,
	}
	if client.Socks5Proxy != "" {
		p, err := proxy.SOCKS5("tcp", client.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		transport.Dial = p.Dial
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Client = &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	return client, nil
}

func (self *Client) Get(url string) (*http.Response, error) {
	retry := self.Retry
	for {
		resp, err := self.Client.Get(url)
		if err != nil {
			if retry == 0 {
				return nil, err
			}
			retry -= 1
			continue
		}
		return resp, nil
	}
}

func (self *Client) GetBytes(url string) ([]byte, error) {
	resp, err := self.Get(url)
	if err != nil {
		return nil, err
	}
	var reader io.Reader = resp.Body
	if self.Encoding != "utf-8" {
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)
		runes, err := From(self.Encoding, buf.Bytes())
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader([]byte(string(runes)))
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return buf.Bytes(), nil
}

func (self *Client) GetDoc(url string) (*goquery.Document, error) {
	resp, err := self.Get(url)
	if err != nil {
		return nil, err
	}
	var reader io.Reader = resp.Body
	if self.Encoding != "utf-8" {
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)
		runes, err := From(self.Encoding, buf.Bytes())
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader([]byte(string(runes)))
	}
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return goquery.NewDocumentFromNode(node), nil
}

func (self *Client) GetFind(url string, selector string) (*goquery.Selection, error) {
	doc, err := self.GetDoc(url)
	if err != nil {
		return nil, err
	}
	return doc.Find(selector), nil
}
