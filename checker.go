package proxy

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	req "github.com/sarovkalach/go_requests"
	log "github.com/sirupsen/logrus"
)

const buffSize = 500000
const (
	testURL = "https://example.com"
)

type checker struct {
	proxyList   chan string
	threadCount int
	timeout     int
	wg          *sync.WaitGroup
	ResCh       chan string
}

func NewChecker(cfg map[string]string) *checker {
	nThreads, _ := strconv.Atoi(cfg["nThreads"])
	requestTimeout, _ := strconv.Atoi(cfg["timeout"])
	c := &checker{
		proxyList:   make(chan string, buffSize),
		threadCount: nThreads,
		timeout:     requestTimeout,
		wg:          &sync.WaitGroup{},
		ResCh:       make(chan string),
	}
	c.readFile(cfg["file"])

	return c
}

func (c *checker) readFile(filename string) {
	readFile, _ := os.Open(filename)
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		c.proxyList <- "http://" + fileScanner.Text()
	}

	log.Debug(fmt.Sprintf("Len of proxylist: %d", len(c.proxyList)))
}

func (c *checker) Start() {
	chunkSize := len(c.proxyList) / c.threadCount
	mod := len(c.proxyList) / c.threadCount

	if mod != 0 {
		chunkSize++
	}
	nChunks := len(c.proxyList) / chunkSize
	mod = len(c.proxyList) % chunkSize

	log.Debug(fmt.Sprintf("nChunks: %d ChunkSize: %d Timeout: %d", nChunks, chunkSize, c.timeout))

	for i := 0; i < nChunks; i++ {
		chunk := make([]string, 0, chunkSize)

		switch i {
		case nChunks - 1:
			for j := 0; j < mod; j++ {
				chunk = append(chunk, <-c.proxyList)
			}
		default:
			for j := 0; j < chunkSize; j++ {
				chunk = append(chunk, <-c.proxyList)
			}
		}

		c.wg.Add(1)
		go c.processChunk(chunk)
	}

	c.wg.Wait()
	close(c.ResCh)
}

func (c *checker) processChunk(chunk []string) {
	for _, proxy := range chunk {
		request := req.NewRequest()
		resp, _ := request.Get(testURL, proxy, c.timeout)
		if resp.StatusCode == http.StatusOK {
			c.ResCh <- proxy[7:]
		}
		fmt.Println(resp.StatusCode)
	}
	c.wg.Done()
}
