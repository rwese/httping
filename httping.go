// httping is a simple tool to verify a http server is available
//
// Usage:
//
// httping -url http://server -code 200 -code 202 -contains musthavestring
//
// Flags:
// 	-code [INT], can appear multiple times to allow multiple http-codes
//  -url [STRING], url to call
//  -timeout [INT], timeout for the response
//  -contains [STRING], check for string in the return
//
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type returnCodes []int

func (c *returnCodes) Set(value string) error {
	var valueInt, err = strconv.Atoi(value)
	if err != nil {
		return errors.New("code must be type int")
	}
	*c = append(*c, valueInt)

	return nil
}

func (c *returnCodes) String() string {
	return fmt.Sprint(*c)
}
func usage() {
	fmt.Print(`Usage: httping [options]

-url		site to call
-code		(multiple) response code must be
-contain	response must contain
-timeout	http-timeout
`)
}
func main() {
	var url string
	var validContains string
	var validCodes returnCodes
	var timeout int

	flag.IntVar(&timeout, "timeout", 5, "timeout in seconds")
	flag.StringVar(&url, "url", "", "target url to check for return code")
	flag.StringVar(&validContains, "contain", "", "body must contain")
	flag.Var(&validCodes, "code", "valid return-codes, use multiple times")

	flag.Parse()
	if len(validCodes) == 0 {
		validCodes.Set("200")
	}

	if url == "" {
		log.Print("url is required")
		usage()
		os.Exit(1)
	}

	runInfo := fmt.Sprintf(`

Testing URL: %s
Valid Return-Codes: %s
Must Contain: %s
Timeout: %d
`, url, fmt.Sprintf("%d", validCodes), validContains, timeout)

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err, runInfo)
	}

	var isValidCode bool
	for _, c := range validCodes {
		if resp.StatusCode == c {
			isValidCode = true
			break
		}
	}
	if !isValidCode {
		log.Fatalf("Invalid return-code [%d] %s", resp.StatusCode, runInfo)
		os.Exit(1)
	}

	if validContains != "" {
		var isValidContains bool
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err, runInfo)
		}

		if strings.Contains(string(body), validContains) {
			isValidContains = true
		}

		if !isValidContains {
			log.Fatalf("Return doesn't contain [%s] %s", validContains, runInfo)
			os.Exit(1)
		}
	}

}
