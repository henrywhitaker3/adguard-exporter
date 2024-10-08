FROM alpine:3.20.2 AS certs

RUN apk add ca-certificates

FROM golang:1.22 AS builder

WORKDIR /build

COPY . /build

RUN go mod download
RUN CGO_ENABLED=0 go build -a -o adguard-exporter main.go

FROM alpine:3.20.2

ARG SREP_VERSION
ENV SREP_VERSION ${SREP_VERSION}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/adguard-exporter /adguard-exporter

ENTRYPOINT [ "/adguard-exporter" ]
