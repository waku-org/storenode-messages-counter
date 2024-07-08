# Storenode Message Counter

used to compare the number of Status messages across Status nodes to understand the potential discrepancies and odd behaviour of messages being inserted in the past in a given pubsub topic.

# Querying the data

The following tables are available in https://superset.bi.status.im/superset/sqllab/ `telemetry` database:

### `missingMessages`
This table contains the list of missing messages, and in which storenode the messages was identified to be missing. The `msgStatus` column indicates the possible scenarios in which a message can be considered missing:
- `does_not_exist`. The message is not stored in the storenode's db
- `unknown`. It was not possible to retrieve information about this message. The store node was likely not available during the time the query to verify the message status was done.

### `storeNodeUnavailable`
This table records timestamps on which a specific storenode was not available. This information should be used along with data from https://kibana.infra.status.im/ and https://grafana.infra.status.im/ to determine the reason why the storenode was not available.

# Browsing the logs

You can browse the logs in kibana at https://kibana.infra.status.im/goto/0be1cef0-1782-11ef-8a2b-91a1a792bd9c

Assuming the link is not available use the following parameters:
- `logsource`: `node-01.he-eu-hel1.tmetry.misc`
- `program`: `docker/telemetry-counter`

# Development

You need to setup a postgres db as such:
1) Create an user with a password
2) Create a db
3) Execute `make build`


Then you can run the program with
```
./build/storemsgcounter --storenode=some_multiaddress --storenode=some_multiaddress  --pubsub-topic=some_pubsubtopic --cluster-id=16 --db-url=postgres://user:password@127.0.0.1:5432/telemetry
```

A dockerfile is also available for ease of setup
```
docker build -t wakuorg/storenode-messages:latest .

docker run wakuorg/storenode-messages:latest --help
```