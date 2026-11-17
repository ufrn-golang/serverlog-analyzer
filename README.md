# Server Log Analyzer

[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

This educational repository features a [Go](https://go.dev/) program that analyzes sample server requests while demonstrating some language's resources:

- The anatomy of a Go program
- Variables and zero values
- Basic and struct types
- Functions
- Control structures
- Error handling
- Slices
- Maps
- Methods
- Interfaces
- Generics

## 📊 About

The analyzer parses request lines and reports the total number of requests, average latency, counts for selected HTTP status classes, hits per path, and the number of slow requests.

Each line has four whitespace-separated fields:

```text
METHOD PATH STATUS LATENCY_MS
```

For example, `GET /home 200 12` describes a request that took 12 milliseconds. The status and latency must be integers. Lines with the wrong number of fields or invalid integers are reported and skipped. Other integer status codes are included in the total, although the report only displays `2xx`, `4xx`, and `5xx` counts. A request is considered slow when its latency exceeds `MAX_LATENCY_MS`, currently 100 milliseconds.

### 🔄 Workflow

1. Define sample log lines in `main`.
2. Parse each line, reporting and skipping invalid entries.
3. Collect request statistics and path hit counts.
4. Print the reports and count slow requests.

The sample lines are intentionally hard-coded; the program does not read a log file or prompt for input.

## 📂 Repository structure

```text
serverlog-analyzer/
├── LICENSE
├── README.md
└── log-analyzer.go
```

## 🚀 Getting started

### ✅ Prerequisites

Install Go 1.18 or later by following the instructions on the official [downloads page](https://go.dev/dl/). The program uses only the Go standard library.

### ▶️ Running

From the repository root, run:

```sh
go run log-analyzer.go
```

The included sample produces:

```text
invalid: bad latency in "GET /home 200 abc": strconv.Atoi: parsing "abc": invalid syntax
skipped: malformed line: "broken line"
7 requests, avg 70.0 ms, 2xx=4 4xx=1 5xx=1
hits: /api/users=3 /home=3 /missing=1
slow requests: 2
analysis finished
```
