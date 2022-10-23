package utils

import (
	"bufio"
	"os"
)

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func ListToChannel(list []string) chan string {
	channel := make(chan string, len(list))
	for _, item := range list {
		channel <- item
	}
	close(channel)

	return channel
}

func PathToChannel(path string) (chan string, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, err
	}

	return ListToChannel(lines), nil
}
