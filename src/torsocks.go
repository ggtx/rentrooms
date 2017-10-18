package main

import (
	"golang.org/x/net/proxy"
	"net/http"
)

const (
	PROXY_ADDR = "127.0.0.1:9050"
)

func callBySocks5(auth *proxy.Auth, theurl string) (*http.Response, error) {
	// creative a socks5 dialer
	p, err := proxy.SOCKS5("tcp", PROXY_ADDR, auth, proxy.Direct)
	if err != nil {
		dlog.Println("make socks5 dialer fails:", err)
		return nil, err
	}
	// set our socks5 as the dialer
	httpTransport := &http.Transport{Dial: p.Dial}
	// setup a http client
	httpClient := &http.Client{Transport: httpTransport}
	// create a request
	req, err := http.NewRequest("GET", theurl, nil)

	if err != nil {
		dlog.Println("make req err:", err)
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		dlog.Println("do req fails:", err)
		return nil, err
	}

	return resp, nil
}
