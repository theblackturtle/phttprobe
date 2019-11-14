package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	USER_AGENT  = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	ACCEPT      = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	ACCEPT_LANG = "en-US,en;q=0.8"
)

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")
	var to int
	flag.IntVar(&to, "t", 10, "timeout (second)")
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "output errors to stderr")
	flag.Parse()

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	var client = &http.Client{
		Timeout:       time.Duration(to) * time.Second,
		CheckRedirect: re,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
				DualStack: true,
			}).DialContext,
			DisableKeepAlives:     true,
			MaxIdleConns:          200,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   20 * time.Second,
			ExpectContinueTimeout: 20 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	urls := make(chan string)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for url := range urls {
				if isListening(client, url) {
					fmt.Println(url)
					continue
				}

				if verbose {
					fmt.Fprintf(os.Stderr, "failed: %s\n", url)
				}
			}
			wg.Done()
		}()
	}

	// accept domains on stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.ToLower(sc.Text())

		lineArgs := strings.Split(line, ",")
		if len(lineArgs) < 3 {
			continue
		}
		domain := lineArgs[0]
		ports := lineArgs[1:]
		for _, p := range ports {
			urls <- fmt.Sprintf("http://%s:%s", domain, p)
			urls <- fmt.Sprintf("https://%s:%s", domain, p)
		}
	}
	close(urls)
	wg.Wait()
}

func isListening(client *http.Client, url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Accept", ACCEPT)
	req.Header.Add("Accept-Language", ACCEPT_LANG)
	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req)
	if resp != nil {
		// Discard response body
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		return false
	}
	// Pass if response status code > 500 (Server Error)
	if resp.StatusCode >= 500 {
		return false
	}
	return true
}
