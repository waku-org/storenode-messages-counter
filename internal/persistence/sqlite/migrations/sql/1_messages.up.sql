CREATE TABLE IF NOT EXISTS syncTopicStatus (
	clusterId INTEGER NOT NULL,
	pubsubTopic VARCHAR NOT NULL,
	lastSyncTimestamp INTEGER NOT NULL,
	PRIMARY KEY (clusterId, pubsubTopic)
) WITHOUT ROWID;

