package main

import (
	"example.com/instrumentation/internal/without/metrics"
	"example.com/instrumentation/internal/without/utils"
)

func main() {
	//metrics.InitMetrics()
	utils.InitMetrics()

	address := "0.0.0.0"
	port := 9090
	endpoint := "/metrics"
	metrics.CreateAndStartServer(address, port, endpoint)

	utils.Count(100)
}
