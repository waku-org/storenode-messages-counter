package main

import (
	"time"

	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/waku-org/go-waku/waku/cliutils"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
)

var cliFlags = []cli.Flag{
	altsrc.NewIntFlag(&cli.IntFlag{
		Name:        "tcp-port",
		Aliases:     []string{"port", "p"},
		Value:       0,
		Usage:       "Libp2p TCP listening port (0 for random)",
		Destination: &options.Port,
		EnvVars:     []string{"STORE_MSG_CTR_TCP_PORT"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "address",
		Aliases:     []string{"host", "listen-address"},
		Value:       "0.0.0.0",
		Usage:       "Listening address",
		Destination: &options.Address,
		EnvVars:     []string{"STORE_MSG_CTR_ADDRESS"},
	}),
	&cli.StringFlag{Name: "config-file", Usage: "loads configuration from a TOML file (cmd-line parameters take precedence)"},
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "fleet-name",
		Usage:       "Fleet name",
		Destination: &options.FleetName,
		EnvVars:     []string{"STORE_MSG_CTR_FLEET_NAME"},
	}),
	cliutils.NewGenericFlagMultiValue(&cli.GenericFlag{
		Name:  "storenode",
		Usage: "Multiaddr of peers that supports storeV3 protocol. Option may be repeated",
		Value: &cliutils.MultiaddrSlice{
			Values: &options.StoreNodes,
		},
		EnvVars: []string{"STORE_MSG_CTR_STORENODE"},
	}),
	altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:        "dns-discovery-url",
		Usage:       "URL for DNS node list in format 'enrtree://<key>@<fqdn>'. Option may be repeated",
		Destination: &options.DNSDiscoveryURLs,
		EnvVars:     []string{"STORE_MSG_CTR_DNS_DISC_URL"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "dns-discovery-name-server",
		Usage:       "DNS nameserver IP to query (empty to use system's default)",
		Destination: &options.DNSDiscoveryNameserver,
		EnvVars:     []string{"STORE_MSG_CTR_DNS_DISC_NAMESERVER"},
	}),
	altsrc.NewUintFlag(&cli.UintFlag{
		Name:        "cluster-id",
		Usage:       "ClusterID to use",
		Destination: &options.ClusterID,
		Value:       0,
		EnvVars:     []string{"STORE_MSG_CTR_CLUSTER_ID"},
	}),
	altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:        "pubsub-topic",
		Required:    true,
		Usage:       "Pubsub topic used for the query. Argument may be repeated.",
		Value:       cli.NewStringSlice(relay.DefaultWakuTopic),
		Destination: &options.PubSubTopics,
		EnvVars:     []string{"STORE_MSG_CTR_PUBSUB_TOPICS"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "db-url",
		Usage:       "The database connection URL for persistent storage.",
		Value:       "",
		Destination: &options.DatabaseURL,
		EnvVars:     []string{"MSG_VERIF_DB_URL"},
	}),
	altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:        "retention-policy",
		Usage:       "Retention policy. ",
		Destination: &options.RetentionPolicy,
		Value:       8 * 24 * time.Hour,
		EnvVars:     []string{"STORE_MSG_CTR_RETENTION_POLICY"},
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Value: &cliutils.ChoiceValue{
			Choices: []string{"DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", "FATAL"},
			Value:   &options.LogLevel,
		},
		Usage:   "Define the logging level (allowed values: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL)",
		EnvVars: []string{"STORE_MSG_CTR_LOG_LEVEL"},
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:  "log-encoding",
		Usage: "Define the encoding used for the logs (allowed values: console, nocolor, json)",
		Value: &cliutils.ChoiceValue{
			Choices: []string{"console", "nocolor", "json"},
			Value:   &options.LogEncoding,
		},
		EnvVars: []string{"STORE_MSG_CTR_LOG_ENCODING"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "log-output",
		Value:       "stdout",
		Usage:       "specifies where logging output should be written  (stdout, file, file:./filename.log)",
		Destination: &options.LogOutput,
		EnvVars:     []string{"STORE_MSG_CTR_LOG_OUTPUT"},
	}),
	altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:        "metrics-server",
		Aliases:     []string{"metrics"},
		Usage:       "Enable the metrics server",
		Destination: &options.EnableMetrics,
		EnvVars:     []string{"STORE_MSG_CTR_METRICS_SERVER"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "metrics-server-address",
		Aliases:     []string{"metrics-address"},
		Value:       "127.0.0.1",
		Usage:       "Listening address of the metrics server",
		Destination: &options.MetricsAddress,
		EnvVars:     []string{"STORE_MSG_CTR_METRICS_SERVER_ADDRESS"},
	}),
	altsrc.NewIntFlag(&cli.IntFlag{
		Name:        "metrics-server-port",
		Aliases:     []string{"metrics-port"},
		Value:       8008,
		Usage:       "Listening HTTP port of the metrics server",
		Destination: &options.MetricsPort,
		EnvVars:     []string{"STORE_MSG_CTR_METRICS_SERVER_PORT"},
	}),
}
