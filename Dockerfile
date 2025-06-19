FROM golang:1.24-alpine3.22 AS builder

WORKDIR /apps/ascii-art-web

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -ldflags "-s -w" ./cmd/server

FROM alpine:3.22


LABEL org.opencontainers.image.title="ascii-art-web-docker"
LABEL org.opencontainers.image.version="1.0"
LABEL org.opencontainers.image.description="A web server that converts user input into ASCII art"
LABEL org.opencontainers.image.authors="ayedige, amirmaga, ademeu"

WORKDIR /

COPY --from=builder /apps/ascii-art-web/server /server

COPY --from=builder /apps/ascii-art-web/assets /assets

EXPOSE 8080

ENTRYPOINT [ "./server" ]
