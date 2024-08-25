package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	url := flag.String("url", "", "🌐 URL of the service to be tested")
	requests := flag.Int("requests", 100, "📊 Total number of requests")
	concurrency := flag.Int("concurrency", 10, "🚀 Number of simultaneous calls")

	flag.Parse()

	if *url == "" {
		color.Red("❌ The service URL is required. Use the --url flag to specify it.")
		return
	}

	color.Cyan("🏁 Starting the load test for %s...", *url)
	runLoadTest(*url, *requests, *concurrency)
}

func runLoadTest(url string, totalRequests int, concurrencyLevel int) {
	var wg sync.WaitGroup
	requestsPerWorker := totalRequests / concurrencyLevel
	extraRequests := totalRequests % concurrencyLevel

	results := make(chan int, totalRequests)
	statusCodeCount := make(map[int]int)
	networkErrorCount := 0
	startTime := time.Now()

	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go func(requests int) {
			defer wg.Done()
			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			for j := 0; j < requests; j++ {
				resp, err := client.Get(url)
				if err != nil {
					color.Red("❌ Network error: %v", err)
					results <- -1
					continue
				}
				results <- resp.StatusCode
				resp.Body.Close()
			}
		}(requestsPerWorker + boolToInt(i < extraRequests))
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for statusCode := range results {
		if statusCode == -1 {
			networkErrorCount++
		} else {
			statusCodeCount[statusCode]++
		}
	}

	totalTime := time.Since(startTime)

	generateReport(totalTime, totalRequests, statusCodeCount, networkErrorCount)
}

func generateReport(totalTime time.Duration, totalRequests int, statusCodeCount map[int]int, networkErrorCount int) {
	color.Green("\n===== 📝 Load Test Report =====")
	fmt.Printf("⏳ Total time: %v\n", totalTime)
	fmt.Printf("📊 Total requests: %d\n", totalRequests)
	color.Cyan("✅ Successful requests (HTTP 200): %d\n", statusCodeCount[200])

	delete(statusCodeCount, 200)

	if len(statusCodeCount) > 0 {
		color.Yellow("\n📉 Distribution of other HTTP status codes:")
		for status, count := range statusCodeCount {
			if status >= 400 {
				color.Red("  ❌ Failed requests (HTTP %d): %d", status, count)
			} else {
				fmt.Printf("  - HTTP %d: %d\n", status, count)
			}
		}
	}

	if networkErrorCount > 0 {
		color.Red("\n❌ Network errors: %d", networkErrorCount)
	}

	color.Magenta("\n⚡ Requests per second: %.2f\n", float64(totalRequests)/totalTime.Seconds())
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
