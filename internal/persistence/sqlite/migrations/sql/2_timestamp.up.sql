CREATE TABLE IF NOT EXISTS missingMessages_new (
	runId TEXT NOT NULL,
	clusterId INTEGER NOT NULL,
	pubsubTopic TEXT NOT NULL,
	messageHash TEXT NOT NULL,
	storenode TEXT NOT NULL,
	msgStatus TEXT NOT NULL,
	storedAt BIGINT NOT NULL, 
	PRIMARY KEY (messageHash, storenode)
);

INSERT INTO missingMessages_new (runId, clusterId, pubsubTopic, messageHash, storenode, msgStatus, storedAt)
SELECT runId, clusterId, pubsubTopic, messageHash, storenode, msgStatus, storedAt
FROM missingMessages;

DROP TABLE missingMessages;

ALTER TABLE missingMessages_new
RENAME TO missingMessages;

