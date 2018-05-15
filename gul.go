package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var num = flag.Int("n", 1, "number of requests")

func main() {
	//flag.PrintDefaults()
	flag.Parse()

	fmt.Println("num :", *num)
	fmt.Println("args : ", flag.Args())

	if len(flag.Args()) == 0 {
		fmt.Println("specify a url")
		os.Exit(0)
	}

	urlRaw := flag.Args()[0]
	fmt.Println("urlRaw : ", urlRaw)
	client := &http.Client{}
	result := runUrl(urlRaw, *num, client)
	reportResult(result)

}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type response struct {
	status   int
	duration time.Duration
}

type runResult struct {
	successCount int
	errorCount   int
	responses    []response
}

func runUrl(urlRaw string, it int, client HttpClient) runResult {
	// TODO LH all the dynamic stuff and then collection of the results
	var result runResult
	result.responses = make([]response, 0)

	for i := 0; i < it; i++ {
		start := time.Now()
		resp, err := client.Get(urlRaw)
		elapsed := time.Since(start)
		if err != nil {
			result.errorCount++
			fmt.Println("err : ", err)
			continue // TODO LH not sure how I feel about this
		}
		result.successCount++
		response := response{
			status:   resp.StatusCode,
			duration: elapsed,
		}
		result.responses = append(result.responses, response)

	}

	return result
}

func reportResult(result runResult) {
	fmt.Println("Success Count : ", result.successCount)
	fmt.Println("Error Count : ", result.errorCount)
	for _, response := range result.responses {
		fmt.Println("resp status : ", response.status)
		fmt.Println("elapsed : ", response.duration)
	}

}
