package main

import (
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
)

type Options struct {
	Port                   int
	Address                string
	LogLevel               string
	LogEncoding            string
	LogOutput              string
	FleetName              string
	ClusterID              uint
	PubSubTopics           cli.StringSlice
	DatabaseURL            string
	RetentionPolicy        time.Duration
	StoreNodes             []multiaddr.Multiaddr
	DNSDiscoveryNameserver string
	DNSDiscoveryURLs       cli.StringSlice
	EnableMetrics          bool
	MetricsAddress         string
	MetricsPort            int
}
