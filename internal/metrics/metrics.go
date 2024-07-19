package metrics

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var missingMessages = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_missing_messages",
		Help: "The messages identified as missing and the reason why they're missing",
	},
	[]string{"storenode", "status"},
)

var totalMissingMessages = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "msgcounter_total_missing_messages",
		Help: "The global total number of missing messages",
	},
)

var storenodeAvailability = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_storenode_availability",
		Help: "Indicate whether a store node is available or not",
	},
	[]string{"storenode"},
)

var topicLastSync = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_topic_lastsyncdate_seconds",
		Help: "Indicates the last syncdate for a pubsubtopic",
	},
	[]string{"pubsubtopic"},
)

var collectors = []prometheus.Collector{
	missingMessages,
	storenodeAvailability,
	totalMissingMessages,
	topicLastSync,
}

// Metrics exposes the functions required to update prometheus metrics for relay protocol
type Metrics interface {
	RecordMissingMessages(peerID peer.ID, status string, length int)
	RecordStorenodeAvailability(peerID peer.ID, available bool)
	RecordTotalMissingMessages(cnt int)
	RecordLastSyncDate(topic string, date time.Time)
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

func (m *metricsImpl) RecordMissingMessages(peerID peer.ID, status string, length int) {
	go func() {
		missingMessages.WithLabelValues(peerID.String(), status).Set(float64(length))
	}()
}

func (m *metricsImpl) RecordStorenodeAvailability(peerID peer.ID, available bool) {
	go func() {
		gaugeValue := float64(1)
		if !available {
			gaugeValue = 0
		}
		storenodeAvailability.WithLabelValues(peerID.String()).Set(gaugeValue)
	}()
}

func (m *metricsImpl) RecordTotalMissingMessages(cnt int) {
	go func() {
		totalMissingMessages.Set(float64(cnt))
	}()
}

func (m *metricsImpl) RecordLastSyncDate(topic string, date time.Time) {
	go func() {
		topicLastSync.WithLabelValues(topic).Set(float64(date.Unix()))
	}()
}
