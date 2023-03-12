# Build Step
FROM golang:1.20.2-alpine AS builder

# Dependencies
RUN apk update && apk add --no-cache upx make git

# Source
WORKDIR $GOPATH/src/github.com/Depado/smallblog
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp
RUN upx --best --lzma /tmp/smallblog


# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/smallblog /go/bin/smallblog
VOLUME [ "/data" ]
WORKDIR /data
ENTRYPOINT ["/go/bin/smallblog"]
