package main

import (
	"flag"

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
	cfg := parseFlags()
	checker := proxy.NewChecker(cfg)

	log.SetLevel(log.DebugLevel)
	go func(ch chan string) {
		for proxy := range ch {
			log.Info(proxy)
		}
	}(checker.ResCh)

	checker.Start()
}
