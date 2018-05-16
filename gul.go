package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"plugin"
	"strings"
	"time"
)

var num = flag.Int("n", 1, "number of requests")
var timeLimit = flag.Duration("t", 0, "duration of how long to test")
var requestBody func() string

func init() {
	flag.Parse()
	buildPlugin()
	loadPlugin()
}

func main() {
	//flag.PrintDefaults()
	// flag.Parse()

	fmt.Println("num :", *num)
	fmt.Println("timeLimit : ", *timeLimit)
	fmt.Println("args : ", flag.Args())
	fmt.Println("plugin request generator : ", requestBody())

	if *timeLimit > 0 {
		*num = 0
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

func buildPlugin() {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "request_plg.go")
	//defer os.Remove("request_plg.so")
	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("err : %q\n", errout.String())
		panic(err)
	}
}

func loadPlugin() {
	p, err := plugin.Open("request_plg.so")
	if err != nil {
		panic(err)
	}

	f, err := p.Lookup("NewRequestBody")
	if err != nil {
		panic(err)
	}
	requestBody = f.(func() string)
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
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

func checkExit(iteration int, totalIterations int, timeLimit time.Duration, startTime time.Time) bool {

	if totalIterations > 0 && iteration >= totalIterations {
		return true
	}
	if timeLimit > 0 && time.Since(startTime) >= timeLimit {
		fmt.Println("time limit hit, break")
		return true
	}
	return false
}

func runUrl(urlRaw string, it int, timeLimit time.Duration, client HttpClient) runResult {
	// TODO LH all the dynamic stuff and then collection of the results
	var result runResult
	result.responses = make([]response, 0)
	result.startTime = time.Now()

	exitLoop := false
	for !exitLoop {
		result.totalRequests++
		start := time.Now()
		request, _ := http.NewRequest("POST", urlRaw, strings.NewReader(requestBody()))
		resp, err := client.Do(request) //client.Get(urlRaw)
		elapsed := time.Since(start)
		if err != nil {
			result.errorCount++
		} else {
			result.successCount++
			response := response{
				status:   resp.StatusCode,
				duration: elapsed,
			}
			result.responses = append(result.responses, response)
			resp.Body.Close()
		}

		exitLoop = checkExit(result.totalRequests, it, timeLimit, result.startTime)
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
