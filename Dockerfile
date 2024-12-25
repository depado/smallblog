# Build Step
FROM golang:1.23-alpine@sha256:6c5c9590f169f77c8046e45c611d3b28fe477789acd8d3762d23d4744de69812 AS builder

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
FROM gcr.io/distroless/static@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58
COPY --from=builder /tmp/smallblog /go/bin/smallblog
COPY templates ./templates
COPY assets ./assets
ENTRYPOINT ["/go/bin/smallblog"]
CMD ["serve", "--server.host=0.0.0.0", "--server.port=8000"]
