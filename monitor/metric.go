package monitor

import (
	"encoding/json"
	"sync/atomic"
)

var instance *Metric

// Metric holds metrics information
type Metric struct {
	HTTPRequests int64 `json:"http.requests"`
	HTTPQueries  int64 `json:"http.queries"`
	HTTPUploads  int64 `json:"http.uploads"`
}

// IncHTTPRequests increment http requests count
func IncHTTPRequests() {
	go atomic.AddInt64(&instance.HTTPRequests, 1)
}

// IncHTTPQueries increment queries count
func IncHTTPQueries() {
	go atomic.AddInt64(&instance.HTTPQueries, 1)
}

// IncHTTPUploads increment uploads count
func IncHTTPUploads() {
	go atomic.AddInt64(&instance.HTTPUploads, 1)
}

// MetricToJSON marshal metric to json
func MetricToJSON() []byte {
	b, _ := json.Marshal(&instance)
	return b
}

// StartMonitor starts monitoring
func StartMonitor() {
	instance = &Metric{}
}
