package main

import (
	"flag"
	"fmt"
	"github.com/ii2day/moon/requester"
	"math"
	"net/http"
	"strings"
	"time"
)

var (
	c = flag.Int("c", 50, "")
	d = flag.Duration("d", 10, "")
	q = flag.Int("q", 2000, "")
	m = flag.String("m", "GET", "")
)

func main() {
	flag.Parse()
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
