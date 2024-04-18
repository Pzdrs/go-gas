package gasstation

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const namespace = "gogas"

type metrics struct {
	carsProcessedTotal prometheus.Counter
	simulationTime     prometheus.Gauge
	fuelTime           *prometheus.HistogramVec
	lineQueueTime      prometheus.Histogram

	registerTime      *prometheus.HistogramVec
	registerQueueTime prometheus.Histogram
}

func registerMetrics() *metrics {
	return &metrics{
		carsProcessedTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "cars_processed_total",
			Help:      "The total number of cars processed",
		}),
		simulationTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "simulation_time",
			Help:      "Time spent simulating the gas station",
		}),

		fuelTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "fuel_time",
			Help:      "Time spent at the pump",
			Buckets:   []float64{1, 2.5, 5, 10, 20, 50},
		}, []string{"fuel_type"}),
		lineQueueTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "line_queue_time",
			Help:      "Time spent waiting for a free pump",
			Buckets:   []float64{1, 2.5, 5, 10, 20, 50},
		}),

		registerTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "register_time",
			Help:      "Time spent at the cash register",
			Buckets:   []float64{1, 2.5, 5, 10, 20, 50},
		}, []string{"register"}),
		registerQueueTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "register_queue_time",
			Help:      "Time spent waiting for a free register",
			Buckets:   []float64{10, 20, 50, 100, 500, 1000},
		}),
	}
}

func (s *GasStation) pushMetrics() {
	s.MetricsWg.Add(1)
	go func() {
		defer s.MetricsWg.Done()
		if err := push.New(configuration.PrometheusConfig.PushGateway, "go_gas").
			Collector(s.Metrics.carsProcessedTotal).
			Collector(s.Metrics.simulationTime).
			Collector(s.Metrics.fuelTime).
			Collector(s.Metrics.lineQueueTime).
			Collector(s.Metrics.registerTime).
			Collector(s.Metrics.registerQueueTime).
			Grouping("simulation", s.SimulationID.String()).
			Push(); err != nil {
			fmt.Println("Could not push metrics to Pushgateway:", err)
		}
	}()
}

func (s *GasStation) logMetric(callback func()) {
	s.MetricsWg.Add(1)
	go func() {
		defer s.MetricsWg.Done()
		callback()
	}()
}
