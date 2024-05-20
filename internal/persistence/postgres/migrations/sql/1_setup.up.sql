CREATE TABLE IF NOT EXISTS syncTopicStatus (
	clusterId INTEGER NOT NULL,
	pubsubTopic TEXT NOT NULL,
	lastSyncTimestamp BIGINT NOT NULL,
	PRIMARY KEY (clusterId, pubsubTopic)
);

CREATE TABLE IF NOT EXISTS missingMessages (
	runId TEXT NOT NULL,
	clusterId INTEGER NOT NULL,
	pubsubTopic TEXT NOT NULL,
	messageHash TEXT NOT NULL,
	msgTimestamp BIGINT NOT NULL,
	storenode TEXT NOT NULL,
	msgStatus TEXT NOT NULL,
	storedAt BIGINT NOT NULL, 
	PRIMARY KEY (messageHash, storenode)
);

CREATE INDEX IF NOT EXISTS idxMsg1 ON missingMessages(storedAt DESC);
CREATE INDEX IF NOT EXISTS idxMsg2 ON missingMessages(runId);

CREATE TABLE IF NOT EXISTS storeNodeUnavailable (
	runId TEXT NOT NULL,
	storenode TEXT NOT NULL,
	requestTime BIGINT NOT NULL, 
	PRIMARY KEY (runId, storenode)
);

CREATE INDEX IF NOT EXISTS idxStr1 ON storeNodeUnavailable(requestTime);
