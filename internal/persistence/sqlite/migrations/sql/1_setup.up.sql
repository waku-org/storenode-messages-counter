CREATE TABLE IF NOT EXISTS syncTopicStatus (
	clusterId INTEGER NOT NULL,
	pubsubTopic VARCHAR NOT NULL,
	lastSyncTimestamp INTEGER NOT NULL,
	PRIMARY KEY (clusterId, pubsubTopic)
) WITHOUT ROWID;


CREATE TABLE IF NOT EXISTS missingMessages (
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
