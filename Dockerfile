FROM alpine:latest

ENV LISTEN_ADDRESS ":8080"

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chown 777 "$GOPATH"

COPY . "$GOPATH/src/github.com/WuHan0608/webhook-listener/"

RUN apk add -U musl-dev go ca-certificates && \
    cd "$GOPATH/src/github.com/WuHan0608/webhook-listener/" && \
    go install -ldflags="-s -w" && \
    apk del -r go && \
    rm -rf /var/cache/apk/*

WORKDIR $GOPATH

CMD ["webhook-listener"]
