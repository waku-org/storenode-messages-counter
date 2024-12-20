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
	[]string{"fleetName", "storenode", "status"},
)

var totalMissingMessages = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_total_missing_messages",
		Help: "The global total number of missing messages",
	},
	[]string{"fleetName"},
)

var storenodeAvailability = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_storenode_availability",
		Help: "Indicate whether a store node is available or not",
	},
	[]string{"fleetName", "storenode"},
)

var topicLastSync = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_topic_lastsyncdate_seconds",
		Help: "Indicates the last syncdate for a pubsubtopic",
	},
	[]string{"fleetName", "pubsubtopic"},
)

var missingMessagesLastHour = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_missing_messages_last_hour",
		Help: "The number of messages missing in last hour (excluding the last 5 minutes)",
	},
	[]string{"fleetName", "storenode"},
)

var missingMessagesLastDay = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_missing_messages_last_day",
		Help: "The number of messages missing in last 24hr (with 2hr delay)",
	},
	[]string{"fleetName", "storenode"},
)

var missingMessagesLastWeek = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_missing_messages_last_week",
		Help: "The number of messages missing in last week (with 2hr delay)",
	},
	[]string{"fleetName", "storenode"},
)

var missingMessagesPreviousHour = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "msgcounter_missing_messages_prev_hour",
		Help: "The number of messages missing in the previous hour",
	},
	[]string{"fleetName", "storenode"},
)

var collectors = []prometheus.Collector{
	missingMessages,
	storenodeAvailability,
	totalMissingMessages,
	topicLastSync,
	missingMessagesPreviousHour,
	missingMessagesLastDay,
	missingMessagesLastWeek,
}

// Metrics exposes the functions required to update prometheus metrics for relay protocol
type Metrics interface {
	RecordMissingMessages(peerID peer.ID, status string, length int)
	RecordStorenodeAvailability(peerID peer.ID, available bool)
	RecordTotalMissingMessages(cnt int)
	RecordLastSyncDate(topic string, date time.Time)
	RecordMissingMessagesLastHour(peerID peer.ID, cnt int)
	RecordMissingMessagesLastDay(peerID peer.ID, cnt int)
	RecordMissingMessagesLastWeek(peerID peer.ID, cnt int)
	RecordMissingMessagesPrevHour(peerID peer.ID, cnt int)
}

type metricsImpl struct {
	log       *zap.Logger
	reg       prometheus.Registerer
	clusterId uint
	fleetName string
}

func NewMetrics(clusterId uint, fleetName string, reg prometheus.Registerer, logger *zap.Logger) Metrics {
	metricshelper.RegisterCollectors(reg, collectors...)
	return &metricsImpl{
		log:       logger,
		reg:       reg,
		clusterId: clusterId,
		fleetName: fleetName,
	}
}

func (m *metricsImpl) RecordMissingMessages(peerID peer.ID, status string, length int) {
	go func() {
		missingMessages.WithLabelValues(m.fleetName, peerID.String(), status).Set(float64(length))
	}()
}

func (m *metricsImpl) RecordStorenodeAvailability(peerID peer.ID, available bool) {
	go func() {
		gaugeValue := float64(1)
		if !available {
			gaugeValue = 0
		}
		storenodeAvailability.WithLabelValues(m.fleetName, peerID.String()).Set(gaugeValue)
	}()
}

func (m *metricsImpl) RecordTotalMissingMessages(cnt int) {
	go func() {
		totalMissingMessages.WithLabelValues(m.fleetName).Set(float64(cnt))
	}()
}

func (m *metricsImpl) RecordLastSyncDate(topic string, date time.Time) {
	go func() {
		topicLastSync.WithLabelValues(m.fleetName, topic).Set(float64(date.Unix()))
	}()
}

func (m *metricsImpl) RecordMissingMessagesLastHour(peerID peer.ID, cnt int) {
	go func() {
		missingMessagesLastHour.WithLabelValues(m.fleetName, peerID.String()).Set(float64(cnt))
	}()
}

func (m *metricsImpl) RecordMissingMessagesLastDay(peerID peer.ID, cnt int) {
	go func() {
		missingMessagesLastDay.WithLabelValues(m.fleetName, peerID.String()).Set(float64(cnt))
	}()
}

func (m *metricsImpl) RecordMissingMessagesLastWeek(peerID peer.ID, cnt int) {
	go func() {
		missingMessagesLastWeek.WithLabelValues(m.fleetName, peerID.String()).Set(float64(cnt))
	}()
}

func (m *metricsImpl) RecordMissingMessagesPrevHour(peerID peer.ID, cnt int) {
	go func() {
		missingMessagesPreviousHour.WithLabelValues(m.fleetName, peerID.String()).Set(float64(cnt))
	}()
}
