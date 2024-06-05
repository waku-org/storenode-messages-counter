FROM golang:1.21 as builder

RUN mkdir -p /go/src/github.com/waku-org/storenode-messages

WORKDIR /go/src/github.com/waku-org/storenode-messages

ADD . .

RUN make

# Copy the binary to the second image
FROM debian:12.5-slim

LABEL maintainer="richard@status.im"
LABEL source="https://github.com/waku-org/storenode-messages"
LABEL description="Storenode message count verifier"
LABEL commit="unknown"

COPY --from=builder /go/src/github.com/waku-org/storenode-messages/build/storemsgcounter /usr/local/bin/storemsgcounter

ENTRYPOINT ["/usr/local/bin/storemsgcounter"]
CMD ["-help"]
