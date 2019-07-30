FROM golang:1.12

COPY . /go/src/go.mozilla.org/sops
WORKDIR /go/src/go.mozilla.org/sops

RUN CGO_ENABLED=1 make install
RUN apt-get update
RUN apt-get install -y vim python-pip emacs
RUN pip install awscli
ENV EDITOR vim
