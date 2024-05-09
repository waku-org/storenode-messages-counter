package main

import (
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
)

type Options struct {
	LogLevel     string
	LogEncoding  string
	LogOutput    string
	ClusterID    uint
	PubSubTopics cli.StringSlice
	DatabaseURL  string
	StoreNodes   []multiaddr.Multiaddr
}
