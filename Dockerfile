# Build Step
FROM golang:1.27.3-alpine@sha256:f2c872bcee949b7be9c490efd02b48a8a2fbf947cdf533f1c58525ec9d8170fd as builder

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
FROM gcr.io/distroless/static@sha256:9e0130bec991f471279b70f1630ada11e46dad439ef2949544f96001225329fb
COPY --from=builder /tmp/smallblog /go/bin/smallblog
COPY templates ./templates
COPY assets ./assets
ENTRYPOINT ["/go/bin/smallblog"]
CMD ["serve", "--server.host=0.0.0.0", "--server.port=8000"]
