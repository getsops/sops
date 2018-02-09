FROM golang:1.8

COPY . /go/src/go.mozilla.org/sops
WORKDIR /go/src/go.mozilla.org/sops

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/sops ./cmd/sops
RUN apt-get update
RUN apt-get install -y vim python-pip emacs
RUN pip install awscli
ENV EDITOR vim