package main

import (
	"runtime"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var queue = make(chan *write.Point, 600)

func startQueue(endpoint api.WriteAPI) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(endpoint)
	}
}

func worker(endpoint api.WriteAPI) {
	for p := range queue {
		endpoint.WritePoint(p)
	}
}
