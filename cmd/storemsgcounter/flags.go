package main

import (
	"time"

	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/waku-org/go-waku/waku/cliutils"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
)

var cliFlags = []cli.Flag{
	&cli.StringFlag{Name: "config-file", Usage: "loads configuration from a TOML file (cmd-line parameters take precedence)"},
	cliutils.NewGenericFlagMultiValue(&cli.GenericFlag{
		Name:     "storenode",
		Required: true,
		Usage:    "Multiaddr of peers that supports storeV3 protocol. Option may be repeated",
		Value: &cliutils.MultiaddrSlice{
			Values: &options.StoreNodes,
		},
		EnvVars: []string{"MSGVERIF_STORENODE"},
	}),
	altsrc.NewUintFlag(&cli.UintFlag{
		Name:        "cluster-id",
		Usage:       "ClusterID to use",
		Destination: &options.ClusterID,
		Value:       0,
		EnvVars:     []string{"MSGVERIF_CLUSTER_ID"},
	}),
	altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:        "pubsub-topic",
		Required:    true,
		Usage:       "Pubsub topic used for the query. Argument may be repeated.",
		Value:       cli.NewStringSlice(relay.DefaultWakuTopic),
		Destination: &options.PubSubTopics,
		EnvVars:     []string{"MSGVERIF_PUBSUB_TOPICS"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "db-url",
		Usage:       "The database connection URL for persistent storage.",
		Value:       "sqlite3://storage.db",
		Destination: &options.DatabaseURL,
		EnvVars:     []string{"MSG_VERIF_DB_URL"},
	}),
	altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:        "retention-policy",
		Usage:       "Retention policy. ",
		Destination: &options.RetentionPolicy,
		Value:       15 * 24 * time.Hour,
		EnvVars:     []string{"MSGVERIF_RETENTION_POLICY"},
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Value: &cliutils.ChoiceValue{
			Choices: []string{"DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", "FATAL"},
			Value:   &options.LogLevel,
		},
		Usage:   "Define the logging level (allowed values: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL)",
		EnvVars: []string{"MSGVERIF_LOG_LEVEL"},
	}),
	cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
		Name:  "log-encoding",
		Usage: "Define the encoding used for the logs (allowed values: console, nocolor, json)",
		Value: &cliutils.ChoiceValue{
			Choices: []string{"console", "nocolor", "json"},
			Value:   &options.LogEncoding,
		},
		EnvVars: []string{"MSGVERIF_LOG_ENCODING"},
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "log-output",
		Value:       "stdout",
		Usage:       "specifies where logging output should be written  (stdout, file, file:./filename.log)",
		Destination: &options.LogOutput,
		EnvVars:     []string{"MSGVERIF_LOG_OUTPUT"},
	}),
}
