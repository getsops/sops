FROM golang:1.7

COPY . /go/src/go.mozilla.org/sops
WORKDIR /go/src/go.mozilla.org/sops
RUN go get -d -v
RUN go test ./...
RUN go install -v ./...
