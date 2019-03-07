package main

import (
	"benchmark_cockroachdb/cli"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("---Start Benchmark CLI---")
	cli.Initialize()
	log.Info("---Shutdown Benchmark CLI---")
}
