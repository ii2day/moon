package main

import (
	"flag"
	"fmt"
	"github.com/ii2day/moon/requester"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	c = flag.Int("c", 50, "")
	d = flag.Duration("d", 10, "")
	q = flag.Int("q", 2000, "")
	m = flag.String("m", "GET", "")
)

var usage = `Usage: hey [options...] <url>

Options:
  -c  Number of workers to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -q  Rate limit, in queries per second (QPS) per worker. Default is no rate limit.
  -d  Duration of application to send requests. When duration is reached,
      application stops and exits. If duration is specified, n is ignored.
      Examples: -d 10s -d 3m.
  -m  HTTP method, one of GET, POST, PUT, DELETE, HEAD, OPTIONS.
  -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
      For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
  -t  Timeout for each request in seconds. Default is 20, use 0 for infinite.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	if flag.NArg() < 1 {
		usageAndExit("")
	}
	dur := *d
	url := flag.Args()[0]
	method := strings.ToUpper(*m)
	req, _ := http.NewRequest(method, url, nil)
	w := &requester.Work{
		Request:       req,
		NumberRequest: math.MaxInt32,
		QPS:           *q,
		Concurrency:   *c,
	}
	w.Init()
	if dur > 0 {
		go func() {
			time.Sleep(dur)
			w.Stop()
		}()
	}
	w.Run()

	metrics := w.AggregateMetric()
	fmt.Printf("Requests: %v \n", metrics.Requests)
	fmt.Printf("Success: %v \n", metrics.Success)
	fmt.Printf("Duration: %v \n", metrics.Duration)
	fmt.Printf("reuqests/sec: %v \n", metrics.Rate)
	latencies := `Latencies:
  P50: %vms 
  P90: %vms 
  P95: %vms 
  P99: %vms 
  Max: %vms 
  Min: %vms 
  Mean: %vms`
	fmt.Println(fmt.Sprintf(
		latencies,
		metrics.Latencies.P50,
		metrics.Latencies.P90,
		metrics.Latencies.P95,
		metrics.Latencies.P99,
		metrics.Latencies.Max,
		metrics.Latencies.Min,
		metrics.Latencies.Mean))
	fmt.Printf("status code: %v \n", metrics.StatusCodes)
	fmt.Printf("errors: %v \n", metrics.Errors)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(0)
}
