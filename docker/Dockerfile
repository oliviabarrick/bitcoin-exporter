FROM golang:1.9.2
WORKDIR /go/src/github.com/justinbarrick/bitcoin-exporter
COPY . .
RUN go get -d -t ./...
RUN go build cmd/bitcoin-exporter.go

FROM debian:stretch
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=0 /go/src/github.com/justinbarrick/bitcoin-exporter/bitcoin-exporter /app/bitcoin-exporter
WORKDIR /app/
ENTRYPOINT [ "/app/bitcoin-exporter" ]
