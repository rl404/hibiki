# Golang base image
FROM golang:1.27.3 AS go_builder
LABEL stage="hibiki_builder"
WORKDIR /hibiki
COPY . .
WORKDIR /hibiki/cmd/hibiki
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o hibiki

# New stage from scratch
FROM alpine:3.24.3
LABEL stage="hibiki_binary"
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 hibiki
COPY --from=go_builder /hibiki/cmd/hibiki/hibiki /usr/local/bin/hibiki
WORKDIR /home/hibiki
USER hibiki
ENTRYPOINT ["/usr/local/bin/hibiki"]
