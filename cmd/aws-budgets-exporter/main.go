package main

import (
	"log"

	"github.com/countsheep123/aws-budgets-exporter/metrics"
	"github.com/countsheep123/aws-budgets-exporter/exporter"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	m, err := metrics.New()
	if err != nil {
		zap.S().Fatal(err)
	}

	if err := exporter.Serve(m); err != nil {
		zap.S().Fatal(err)
	}
}
