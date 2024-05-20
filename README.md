# Storenode Message Counter

used to compare the number of Status messages across Status nodes to understand the potential discrepancies and odd behaviour of messages being inserted in the past in a given pubsub topic.


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