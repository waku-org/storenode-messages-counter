package metrics

import (
	"github.com/libp2p/go-libp2p/p2p/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var missingMessages = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "missing_messages",
		Help: "The messages identified as missing and the reason why they're missing",
	},
	[]string{"storenode", "status"},
)

var storenodeUnavailable = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "storenode_unavailable",
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
	RecordMissingMessage(storenode string, status string)
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

func (m *metricsImpl) RecordMissingMessage(storenode string, status string) {
	go func() {
		missingMessages.WithLabelValues(storenode, status).Inc()
	}()
}

func (m *metricsImpl) RecordStorenodeUnavailable(storenode string) {
	go func() {
		storenodeUnavailable.WithLabelValues(storenode).Inc()
	}()
}
