package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Proxy struct {
	Host     string
	Port     string
	User     string
	Password string
	Protocol string
}

func FromString(proxyStr string, protocol string) (*Proxy, error) {
	splittedProxy := strings.Split(strings.TrimSpace(proxyStr), ":")

	if len(splittedProxy) < 2 {
		return nil, errors.New("invalid format for proxy")
	}

	proxy := &Proxy{
		Protocol: protocol,
		Host:     splittedProxy[0],
		Port:     splittedProxy[1],
	}

	if len(splittedProxy) == 4 {
		proxy.User = splittedProxy[2]
		proxy.Password = splittedProxy[3]
	}

	return proxy, nil
}

func (p *Proxy) GetURL() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}

func (p *Proxy) CreateTransport() (*http.Transport, error) {
	proxyAuth := ""
	if p.User != "" && p.Password != "" {
		proxyAuth = fmt.Sprintf("%s:%s@", p.User, p.Password)
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s://%s%s:%s", p.Protocol, proxyAuth, p.Host, p.Port))
	if err != nil {
		return nil, err
	}

	return &http.Transport{Proxy: http.ProxyURL(parsedURL)}, nil
}
