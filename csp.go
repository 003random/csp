package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

/*
	Given a list of hosts, this small utility fetches all whitelisted
	domains from the hosts' CSPs. I use this for reconnaissance
	purposes while bug bounty hunting.

	$ cat hosts.txt
	http://example.com/
	$ cat hosts.txt | csp
	example.com
	subdomain.example.com
	...
*/

// Forgive me, Father, for I have sinned.
// Father: You have been forgiven, son
var r = regexp.MustCompile("(([a-zA-Z](-?[a-zA-Z0-9])*)\\.)+[a-zA-Z]{2,}")

// Consider using https://github.com/tike/csp in future.

// requestCSP requests and extracts CSP for a given URL.
func requestCSP(client *http.Client, url string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %v", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response, cannot analyze")
	}
	defer resp.Body.Close()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not discard response body: %v", err)
	}

	csp := resp.Header.Get("Content-Security-Policy")
	results := r.FindAllString(csp, -1)

	return results, err
}

func main() {
	concurrency := flag.Int("c", 20, "set the concurrency level")
	flag.Parse()

	urlsChannel := make(chan string)

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()

			// Stolen from my mentor, TomNomNom! 🐑
			var tr = &http.Transport{
				MaxIdleConns:      30,
				IdleConnTimeout:   time.Second,
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			}

			client := &http.Client{
				Transport: tr,
			}

			for url := range urlsChannel {
				res, err := requestCSP(client, url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not get CSP: %v\n", err)
					return
				}

				// Ensure we do not print out an empty string.
				if len(res) > 0 {
					fmt.Println(strings.Join(res, "\n"))
				}
			}
		}()
	}

	/*
		Read input from stdin.
		$ cat hosts.txt | csp
	*/
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		host := strings.ToLower(sc.Text())
		urlsChannel <- host
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
	}

	close(urlsChannel)
	wg.Wait()
}
