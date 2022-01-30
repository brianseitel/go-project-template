package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brianseitel/go-project-template/internal/metrics"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TimingMiddleware is a set of middleware that helps track how long each request
// takes. First, it identifies the current timestamp when the request begins, then
// it executes the request, then it calculates the elapsed time. Finally, it takes
// the total elapsed time and adds it to the X-Timing header.
func TimingMiddleware(next http.Handler) http.Handler {
	logger, _ := zap.NewDevelopment()
	cfg := metrics.Config{
		Host:        viper.GetString("grafana_url"),
		Source:      viper.GetString("app_name"),
		User:        viper.GetString("grafana_user"),
		APIKey:      viper.GetString("grafana_apikey"),
		Environment: viper.GetString("environment"),
	}

	_, err := metrics.New(cfg, logger)
	if err != nil {
		logger.Sugar().Errorf("failed to instantiate metrics middleware: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Current time
		start := time.Now()
		// Continue executing the request
		next.ServeHTTP(w, r)
		// This gets current time again and subtracts the start time from it,
		// which should give us the total elapsed time of the above request.
		elapsed := time.Since(start)

		// Build the metric name, which is: {environment}.{app_name}.{endpoint}.{method}.{name}
		// For example: development.ziggy.v1_hello.get.elapsed
		endpoint := intoMetricName(r.URL.Path)
		metricName := fmt.Sprintf("%s.%s.%s.%s.%s", cfg.Environment, cfg.Source, endpoint, strings.ToLower(r.Method), "elapsed")

		// Send elapsed time to Grafana
		metrics.SendMetric(metricName, int(elapsed))
	})
}

// Replaces /, @, or . with _
var replacer = strings.NewReplacer("/", "_", "@", "_", ".", "_")

func intoMetricName(path string) string {
	path = replacer.Replace(strings.TrimRight(path, "/"))

	return path
}
