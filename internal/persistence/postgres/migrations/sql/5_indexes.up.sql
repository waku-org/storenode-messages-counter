
CREATE INDEX IF NOT EXISTS idxMsg3 ON missingMessages(storedAt DESC, fleet);
CREATE INDEX IF NOT EXISTS idxMsg4 ON missingMessages(fleet, clusterId, messageHash);
CREATE INDEX IF NOT EXISTS idxFleetStore ON storenodeunavailable(fleet);

ALTER TABLE synctopicstatus DROP CONSTRAINT synctopicstatus_pkey;
ALTER TABLE synctopicstatus ADD CONSTRAINT synctopicstatus_pkey PRIMARY KEY (clusterId, pubsubTopic, fleet);

ALTER TABLE storenodeunavailable DROP CONSTRAINT storenodeunavailable_pkey;
ALTER TABLE storenodeunavailable ADD CONSTRAINT storenodeunavailable_pkey PRIMARY KEY (runId, storenode, fleet);

ALTER TABLE missingMessages DROP CONSTRAINT missingmessages_pkey;
ALTER TABLE missingMessages ADD CONSTRAINT missingmessages_pkey PRIMARY KEY (messageHash, storenode, fleet);
