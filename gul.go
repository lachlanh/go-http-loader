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

	runUrl(urlRaw, *num)

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

func runUrl(urlRaw string, it int) runResult {
	// TODO LH all the dynamic stuff and then collection of the results
	var result runResult
	result.responses = make([]response, it)
	for i := 0; i < it; i++ {
		start := time.Now()
		resp, err := http.Get(urlRaw)
		elapsed := time.Since(start)
		if err != nil {
			result.errorCount++
			continue // TODO LH not sure how I feel about this
		}
		result.successCount++
		response := response{
			status:   resp.StatusCode,
			duration: elapsed,
		}
		result.responses = append(result.responses, response)
		fmt.Println("err : ", err)
		fmt.Println("resp : ", resp)
		fmt.Println("elapsed : ", elapsed)
	}

	return result
}
