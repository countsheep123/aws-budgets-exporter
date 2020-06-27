package exporter

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"go.uber.org/zap"
)

type config struct {
	ListenAddr      string `env:"LISTEN_ADDR" envDefault:":8080"`
	MetricsEndpoint string `env:"METRICS_ENDPOINT" envDefault:"/metrics"`
	Timeout         string `env:"TIMEOUT" envDefault:"5s"`
}

func Serve(m Metrics) error {
	// env
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return err
	}

	// collector
	zap.S().Infof("Starting %s %s", m.GetName(), version.Info())
	zap.S().Infof("Build context: %s", version.BuildContext())

	prometheus.MustRegister(version.NewCollector(m.GetName()))

	collector, err := new(m)
	if err != nil {
		return err
	}
	prometheus.MustRegister(collector)

	// http server
	mux := http.NewServeMux()
	mux.Handle(cfg.MetricsEndpoint, promhttp.Handler())

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	d, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(ctx); err != nil {
			zap.S().Fatal("HTTP server shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed

	return nil
}
