# Golang base image
FROM golang:1.18 as go_builder
LABEL stage="hibiki_builder"
WORKDIR /hibiki
COPY . .
WORKDIR /hibiki/cmd/hibiki
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o hibiki

# New stage from scratch
FROM alpine:3.16
LABEL stage="hibiki_binary"
RUN apk add --no-cache ca-certificates && apk add --no-cache tzdata
COPY --from=go_builder /hibiki/cmd/hibiki/hibiki /cmd/hibiki/hibiki
WORKDIR /cmd/hibiki
