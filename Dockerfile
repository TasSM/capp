FROM golang:1.14 as builder

WORKDIR /build

COPY src/ .
RUN go mod download

COPY application/ .

RUN GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o main .

FROM alpine:3
# Requires $REDIS_HOST at runtime (host:port)

RUN apk --no-cache add ca-certificates
WORKDIR /dist
RUN mkdir /dist/web
COPY application/web/ ./web/
COPY --from=builder /build/main .

EXPOSE 8080

CMD ["/dist/main"]