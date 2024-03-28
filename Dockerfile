ARG IMAGE=scratch
ARG OS=linux
ARG ARCH=amd64

FROM golang:1.21.5-alpine3.17 as builder

WORKDIR /go/src/github.com/henrywhitaker3/adguard-exporter
COPY . .

RUN apk --no-cache add git alpine-sdk

RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags '-s -w' -o binary ./

FROM $IMAGE

LABEL name="pihole-exporter"

WORKDIR /root/
COPY --from=builder /go/src/github.com/henrywhitaker3/adguard-exporter/binary pihole-exporter

CMD ["./pihole-exporter"]
