package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("Starting metrics server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}
