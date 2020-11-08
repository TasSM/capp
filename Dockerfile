FROM golang:alpine as builder

WORKDIR /build

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o /bin/main cmd/capp/main.go

FROM alpine
# Requires $REDIS_HOST at runtime (host:port)

RUN apk --no-cache add ca-certificates
WORKDIR /dist
COPY --from=builder /bin/main .

EXPOSE 8080

CMD ["/dist/main"]