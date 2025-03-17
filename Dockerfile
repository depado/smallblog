# Build Step
FROM golang:1.24.1-alpine@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 as builder

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
FROM gcr.io/distroless/static@sha256:95ea148e8e9edd11cc7f639dc11825f38af86a14e5c7361753c741ceadef2167
COPY --from=builder /tmp/smallblog /go/bin/smallblog
COPY templates ./templates
COPY assets ./assets
ENTRYPOINT ["/go/bin/smallblog"]
CMD ["serve", "--server.host=0.0.0.0", "--server.port=8000"]
