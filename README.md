## Describe 
An HTTP benchmarking tool based on [Hey](https://github.com/rakyll/hey) and [go-wrk](https://github.com/tsliwowicz/go-wrk) implementations has smaller memory consumption, higher performance than Hey, and a more comprehensive metric output than go-wrk

## Usage

```shell
~# moon -q 200 -c 50 -d 10s -m GET http://www.baidu.com
```

## Result
```shell
~# moon -q 200 -c 50 -d 10s http://www.baidu.com
Requests: 183
Success: 183
Duration: 32.297199209s
reuqests/sec: 5.666125994882085
Latencies:
  P50: 1447ms
  P90: 8842ms
  P95: 12718.5ms
  P99: 22612.5ms
  Max: 32044ms
  Min: 93ms
  Mean: 3571.0764ms
status code: map[200:183]
errors: map[]
```