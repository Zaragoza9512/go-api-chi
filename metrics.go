package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ====================================================================
// MÉTRICAS DE PROMETHEUS
// ====================================================================

// 1. Counter: Total de requests HTTP
var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de peticiones HTTP recibidas",
	},
	[]string{"method", "endpoint", "status"}, // Labels
)

// 2. Histogram: Duración de requests
var httpRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duración de las peticiones HTTP en segundos",
		Buckets: prometheus.DefBuckets, // Buckets por defecto
	},
	[]string{"method", "endpoint"},
)

// 3. Gauge: Requests en vuelo
var httpRequestsInFlight = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Número de peticiones HTTP siendo procesadas actualmente",
	},
)

// ====================================================================
// RESPONSE WRITER PERSONALIZADO
// ====================================================================

// ResponseWriter personalizado que captura el código de estado
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ====================================================================
// MIDDLEWARE DE MÉTRICAS
// ====================================================================

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Incrementar requests en vuelo
		httpRequestsInFlight.Inc()

		// 2. Asegurar que SIEMPRE se decremente (usa defer)
		defer httpRequestsInFlight.Dec()

		// 3. Capturar tiempo de inicio
		start := time.Now()

		// 4. Crear responseWriter personalizado para capturar status
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     200, // Status por defecto
		}

		// 5. Ejecutar el siguiente handler
		next.ServeHTTP(rw, r)

		// 6. Calcular duración
		duration := time.Since(start).Seconds()

		// 7. Obtener datos del request
		method := r.Method
		endpoint := r.URL.Path
		status := strconv.Itoa(rw.statusCode)

		// 8. Registrar métricas
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	})
}