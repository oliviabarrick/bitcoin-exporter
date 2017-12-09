FROM golang:1.9.2

WORKDIR /go/src/usps
COPY . .

RUN go-wrapper download

RUN go get -u github.com/golang/dep/cmd/dep
