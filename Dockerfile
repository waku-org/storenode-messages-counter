FROM golang:1.20 as builder

RUN mkdir -p /go/src/github.com/waku-org/storenode-messages

WORKDIR /go/src/github.com/waku-org/storenode-messages

ADD . .

RUN make

# Copy the binary to the second image
FROM debian:12.5-slim

LABEL maintainer="richard@status.im"
LABEL source="https://github.com/waku-org/storenode-messages"
LABEL description="Storenode message count verifier"

COPY --from=builder /go/src/github.com/waku-org/storenode-messages/build/storeverif /usr/local/bin/storeverif

ENTRYPOINT ["/usr/local/bin/storeverif"]
CMD ["-help"]