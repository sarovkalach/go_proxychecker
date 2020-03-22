package main

import (
	"bufio"
	"flag"
	"os"

	proxy "github.com/sarovkalach/go_proxychecker"
	log "github.com/sirupsen/logrus"
)

func parseFlags() map[string]string {
	file := flag.String("file", "proxy.txt", "file to process")
	nThreads := flag.String("n", "128", "N threads")
	timeout := flag.String("t", "5", "request timeout")

	flag.Parse()

	return map[string]string{
		"file":     *file,
		"nThreads": *nThreads,
		"timeout":  *timeout,
	}
}

func main() {
	file, _ := os.OpenFile("../clean_proxy.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	writer := bufio.NewWriter(file)

	cfg := parseFlags()
	checker := proxy.NewChecker(cfg)

	log.SetLevel(log.DebugLevel)
	go func(writer *bufio.Writer) {
		for proxy := range checker.ResCh {
			log.Info(proxy)
			writer.WriteString(proxy + "\n")
		}
	}(writer)
	defer writer.Flush()

	checker.Start()
}
