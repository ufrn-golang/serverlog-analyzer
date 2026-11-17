package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const MAX_LATENCY_MS = 100

// ErrMalformed is returned by parseLine when the line is
// not in the expected format.
var ErrMalformed = errors.New("malformed line")

// Request represents one log line.
type Request struct {
	method    string
	path      string
	status    int
	latencyMs int
}

// parseLine reads "METHOD PATH STATUS LATENCY_MS".
// It returns ErrMalformed if the line is not in that format.
func parseLine(line string) (Request, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Request{}, fmt.Errorf("%w: %q", ErrMalformed, line)
	}
	status, err := strconv.Atoi(fields[2])
	if err != nil {
		return Request{}, fmt.Errorf("bad status in %q: %w", line, err)
	}
	latencyMs, err := strconv.Atoi(fields[3])
	if err != nil {
		return Request{}, fmt.Errorf("bad latency in %q: %w", line, err)
	}
	return Request{
		method:    fields[0],
		path:      fields[1],
		status:    status,
		latencyMs: latencyMs,
	}, nil
}

// Stats collects overall request statistics.
type Stats struct {
	count   int         // number of requests
	totalMs int         // total latency in milliseconds
	byClass map[int]int // counts of 2xx, 4xx, 5xx
}

// Add adds one request to the stats.
func (s *Stats) Add(r Request) {
	if s.byClass == nil {
		s.byClass = make(map[int]int)
	}
	s.byClass[r.status/100]++
	s.count++
	s.totalMs += r.latencyMs
}

// AvgMs returns the average latency in milliseconds, or 0 if there are no requests.
func (s Stats) AvgMs() float64 {
	if s.count == 0 {
		return 0
	}
	return float64(s.totalMs) / float64(s.count)
}

// Hits counts requests per path.
type Hits map[string]int

// Reporter is satisfied by any type with a Report method.
type Reporter interface {
	Report() string
}

// Report returns a summary of the stats.
func (s Stats) Report() string {
	return fmt.Sprintf("%d requests, avg %.1f ms, 2xx=%d 4xx=%d 5xx=%d",
		s.count, s.AvgMs(), s.byClass[2], s.byClass[4], s.byClass[5])
}

// Report returns a summary of the hits.
func (h Hits) Report() string {
	paths := make([]string, 0, len(h))
	for path := range h {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		parts = append(parts, fmt.Sprintf("%s=%d", path, h[path]))
	}
	return "hits: " + strings.Join(parts, " ")
}

// Adds the request to the stats and hits, and recovers from any panic.
func process(s *Stats, h Hits, r Request) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("process %s: %v", r.path, p)
		}
	}()
	s.Add(r)
	h[r.path]++
	return nil
}

// Filter returns a new slice containing only the items for which
// keep returns true.
func Filter[T any](items []T, keep func(T) bool) []T {
	var out []T
	for _, item := range items {
		if keep(item) {
			out = append(out, item)
		}
	}
	return out
}

func parseLines(lines []string) []Request {
	var requests []Request
	for _, line := range lines {
		request, err := parseLine(line)
		if errors.Is(err, ErrMalformed) {
			fmt.Println("skipped:", err)
			continue
		}
		if err != nil {
			fmt.Println("invalid:", err)
			continue
		}
		requests = append(requests, request)
	}
	return requests
}

func analyzeRequests(requests []Request) (Stats, Hits) {
	var stats Stats
	hits := Hits{}
	for _, request := range requests {
		if err := process(&stats, hits, request); err != nil {
			fmt.Println("rejected:", err)
		}
	}
	return stats, hits
}

func main() {
	defer fmt.Println("analysis finished")

	lines := []string{
		"GET /home 200 12",
		"GET /api/users 200 134",
		"POST /api/users 201 87",
		"GET /missing 404 5",
		"GET /api/users 500 240",
		"GET /home 200 abc",
		"broken line",
		"GET /home 999 3",
		"GET /home 200 9",
	}

	requests := parseLines(lines)
	stats, hits := analyzeRequests(requests)

	for _, reporter := range []Reporter{stats, hits} {
		fmt.Println(reporter.Report())
	}

	// Filter slow requests
	slow := Filter(requests, func(r Request) bool {
		return r.latencyMs > MAX_LATENCY_MS
	})
	fmt.Println("slow requests:", len(slow))
}
