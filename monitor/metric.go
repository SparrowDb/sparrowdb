package monitor

import (
	"encoding/json"
	"sync/atomic"
)

var instance *Metric

type Metric struct {
	Http_requests int64 `json:"http.requests"`
	Http_queries  int64 `json:"http.queries"`
	Http_uploads  int64 `json:"http.uploads"`
}

func IncHTTPRequests() {
	go atomic.AddInt64(&instance.Http_requests, 1)
}

func IncHTTPQueries() {
	go atomic.AddInt64(&instance.Http_queries, 1)
}

func IncHTTPUploads() {
	go atomic.AddInt64(&instance.Http_uploads, 1)
}

func MetricToJSON() []byte {
	b, _ := json.Marshal(&instance)
	return b
}

func StartMonitor() {
	instance = &Metric{}
}
