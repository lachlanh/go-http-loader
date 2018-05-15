package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var num = flag.Int("n", 1, "number of requests")
var timeLimit = flag.Duration("t", 0, "duration of how long to test")

func main() {
	//flag.PrintDefaults()
	flag.Parse()

	fmt.Println("num :", *num)
	fmt.Println("timeLimit : ", *timeLimit)
	fmt.Println("args : ", flag.Args())

	if timeLimit != nil {
		*num = 50000 //similar to ab
	}

	if len(flag.Args()) == 0 {
		fmt.Println("specify a url")
		os.Exit(0)
	}

	urlRaw := flag.Args()[0]
	fmt.Println("urlRaw : ", urlRaw)
	client := &http.Client{}
	result := runUrl(urlRaw, *num, *timeLimit, client)
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
	totalRequests int
	successCount  int
	errorCount    int
	startTime     time.Time
	endTime       time.Time
	responses     []response
}

func runUrl(urlRaw string, it int, timeLimit time.Duration, client HttpClient) runResult {
	// TODO LH all the dynamic stuff and then collection of the results
	var result runResult
	result.responses = make([]response, 0)
	result.startTime = time.Now()
	for i := 0; i < it; i++ {
		result.totalRequests++
		start := time.Now()
		resp, err := client.Get(urlRaw)
		elapsed := time.Since(start)
		if err != nil {
			result.errorCount++
		} else {
			result.successCount++
			response := response{
				status:   resp.StatusCode,
				duration: elapsed,
			}
			//resp.Close()
			result.responses = append(result.responses, response)
		}
		if time.Since(result.startTime) >= timeLimit {
			fmt.Println("time limit hit, break")
			break
		}

	}
	result.endTime = time.Now()
	return result
}

func reportResult(result runResult) {
	fmt.Println("Total Requests : ", result.totalRequests)
	fmt.Println("Success Count : ", result.successCount)
	fmt.Println("Error Count : ", result.errorCount)
	elapsed := result.endTime.Sub(result.startTime).Seconds()
	fmt.Println("Elapsed : ", elapsed)
	rps := float64(result.totalRequests) / float64(elapsed)
	fmt.Println("RPS : ", rps)
	// for _, response := range result.responses {
	// 	fmt.Println("resp status : ", response.status)
	// 	fmt.Println("elapsed : ", response.duration)
	// }
}
