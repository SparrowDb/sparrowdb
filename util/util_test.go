package util

import (
	"log"
	"testing"
)

func Test_BloomFilter(T *testing.T) {
	bf := NewBloomFilter(500, 0.01)
	bf.Add("www.github.com")
	bf.Add("www.google.com")
	bf.Add("www.yahoo.com")
	bf.Add("www.bing.com")

	slog.Infof(("%v", bf.Contains("www.github.com"))
	slog.Infof(("%v", bf.Contains("www.bing.com"))
	slog.Infof(("%v", bf.Contains("www.ebay.com"))
	slog.Infof(("%v", bf.Contains("www.google.com"))
}
