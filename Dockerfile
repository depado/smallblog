# Build Step
FROM golang:1.25.5-alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 as builder

# Dependencies
RUN apk update && apk add --no-cache make git

# Source
WORKDIR $GOPATH/src/github.com/Depado/smallblog
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp

# Final Step
FROM gcr.io/distroless/static@sha256:4b2a093ef4649bccd586625090a3c668b254cfe180dee54f4c94f3e9bd7e381e
COPY --from=builder /tmp/smallblog /go/bin/smallblog
COPY templates ./templates
COPY assets ./assets
ENTRYPOINT ["/go/bin/smallblog"]
CMD ["serve", "--server.host=0.0.0.0", "--server.port=8000"]
