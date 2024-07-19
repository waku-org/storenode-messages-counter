package metrics

import (
	"github.com/libp2p/go-libp2p/p2p/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var missingMessages = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "counter_missing_messages",
		Help: "The messages identified as missing and the reason why they're missing",
	},
	[]string{"storenode", "status"},
)

var storenodeAvailability = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "counter_storenode_availability",
		Help: "Indicate whether a store node is available or not",
	},
	[]string{"storenode"},
)

var collectors = []prometheus.Collector{
	missingMessages,
	storenodeAvailability,
}

// Metrics exposes the functions required to update prometheus metrics for relay protocol
type Metrics interface {
	RecordMissingMessages(storenode string, status string, length int)
	RecordStorenodeAvailability(storenode string, available bool)
}

type metricsImpl struct {
	log *zap.Logger
	reg prometheus.Registerer
}

func NewMetrics(reg prometheus.Registerer, logger *zap.Logger) Metrics {
	metricshelper.RegisterCollectors(reg, collectors...)
	return &metricsImpl{
		log: logger,
		reg: reg,
	}
}

func (m *metricsImpl) RecordMissingMessages(storenode string, status string, length int) {
	go func() {
		missingMessages.WithLabelValues(storenode, status).Set(float64(length))
	}()
}

func (m *metricsImpl) RecordStorenodeAvailability(storenode string, available bool) {
	go func() {
		gaugeValue := float64(1)
		if !available {
			gaugeValue = 0
		}
		storenodeAvailability.WithLabelValues(storenode).Set(gaugeValue)
	}()
}
