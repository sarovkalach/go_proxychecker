package proxy

import (
	requests "github.com/sarovkalach/go_requests"
	"sync"
)

// var threadCounter = func() string {
//
// }

type checker struct {
	proxyList   []string
	threadCount int
	wg          *sync.WaitGroup
}

func newChecker(proxyList *[]string, nThreads int) *checker {
	proxy := addTag(proxyList)
	return &checker{
		proxyList:   proxy,
		threadCount: nThreads,
		wg:          &sync.WaitGroup{},
	}
}

func (c *checker) Start() {

}

func addTag(proxyList *[]string) []string {
	finalProxyList := make([]string, len(*proxyList))
	for i, proxy := range *proxyList {
		finalProxyList[i] = "http://" + proxy
	}

	return finalProxyList
}

func (c *checker) splitToChunks() {
	chunkSize := len(c.proxyList) / c.threadCount
	div := len(c.proxyList) / c.threadCount

	if div != 0 {
		chunkSize++
	}

}

func processChunk(chunk []string) []string {
	results := make([]string, 0, len(chunk)/2)
	wg := &sync.WaitGroup{}

	for _, proxy := range chunk {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			response := requests.Get()
		}(wg)
	}
}
