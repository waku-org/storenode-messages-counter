CREATE TABLE IF NOT EXISTS syncTopicStatus (
	clusterId INTEGER NOT NULL,
	pubsubTopic VARCHAR NOT NULL,
	lastSyncTimestamp INTEGER NOT NULL,
	PRIMARY KEY (clusterId, pubsubTopic)
) WITHOUT ROWID;


CREATE TABLE IF NOT EXISTS missingMessages (
	runId VARCHAR NOT NULL,
	clusterId INTEGER NOT NULL,
	pubsubTopic VARCHAR NOT NULL,
	messageHash VARCHAR NOT NULL,
	msgTimestamp INTEGER NOT NULL,
	storenode VARCHAR NOT NULL,
	msgStatus VARCHAR NOT NULL,
	storedAt INTEGER NOT NULL, 
	PRIMARY KEY (messageHash, storenode)
) WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS idxMsg1 ON missingMessages(storedAt DESC);
CREATE INDEX IF NOT EXISTS idxMsg2 ON missingMessages(runId);

CREATE TABLE IF NOT EXISTS storeNodeUnavailable (
	runId VARCHAR NOT NULL,
	storenode VARCHAR NOT NULL,
	requestTime INTEGER NOT NULL, 
	PRIMARY KEY (runId, storenode)
) WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS idxStr1 ON storeNodeUnavailable(requestTime);
