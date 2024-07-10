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

var storenodeUnavailable = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "counter_storenode_unavailable",
		Help: "Number of PubSub Topics node is subscribed to",
	},
	[]string{"storenode"},
)

var collectors = []prometheus.Collector{
	missingMessages,
	storenodeUnavailable,
}

// Metrics exposes the functions required to update prometheus metrics for relay protocol
type Metrics interface {
	RecordMissingMessages(storenode string, status string, length int)
	RecordStorenodeUnavailable(storenode string)
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

func (m *metricsImpl) RecordStorenodeUnavailable(storenode string) {
	go func() {
		storenodeUnavailable.WithLabelValues(storenode).Set(1)
	}()
}
