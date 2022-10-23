package proxy

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

type Pool struct {
	Proxies  []string
	Protocol string
	Index    int
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func PoolFromFile(path string, protocol string) (*Pool, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, err
	}

	return NewPool(lines, protocol), nil
}

func NewPool(proxies []string, protocol string) *Pool {
	return &Pool{
		Proxies:  proxies,
		Protocol: protocol,
		Index:    0,
	}
}

func (p *Pool) Shuffle() {
	for i := range p.Proxies {
		j := rand.Intn(i + 1)
		p.Proxies[i], p.Proxies[j] = p.Proxies[j], p.Proxies[i]
	}
}

func (p *Pool) Next() {
	p.Index += 1
	if p.Index >= len(p.Proxies) {
		p.Index = 0
	}
}

func (p *Pool) Get() string {
	proxyStr := p.Proxies[p.Index]
	p.Next()

	return proxyStr
}

func (p *Pool) GetProxy() (*Proxy, error) {
	proxy, err := FromString(p.Get(), p.Protocol)
	if err != nil {
		return nil, err
	}

	return proxy, nil
}

func (p *Pool) GetTransport() *http.Transport {
	proxy, err := p.GetProxy()
	if err != nil {
		fmt.Println("[Proxy Pool] Error creating proxy:", err)
		return nil
	}

	transport, err := proxy.CreateTransport()
	if err != nil {
		fmt.Println("[Proxy Pool] Error creating transport:", err)
		return nil
	}

	return transport
}
