# Based on docker-library's golang 1.6 alpine and wheezy docker files.
# https://github.com/docker-library/golang/blob/master/1.6/alpine/Dockerfile
# https://github.com/docker-library/golang/blob/master/1.6/wheezy/Dockerfile
FROM buildpack-deps:wheezy-scm

ENV GOLANG_VERSION tip
ENV GOLANG_SRC_REPO_URL https://go.googlesource.com/go

ENV GOLANG_BOOTSTRAP_VERSION 1.6.2
ENV GOLANG_BOOTSTRAP_URL https://golang.org/dl/go$GOLANG_BOOTSTRAP_VERSION.linux-amd64.tar.gz
ENV GOLANG_BOOTSTRAP_SHA256 e40c36ae71756198478624ed1bb4ce17597b3c19d243f3f0899bb5740d56212a
ENV GOLANG_BOOTSTRAP_PATH /usr/local/bootstrap

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		git \
	&& rm -rf /var/lib/apt/lists/*

# Setup the Bootstrap
RUN mkdir -p "$GOLANG_BOOTSTRAP_PATH" \
	&& curl -fsSL "$GOLANG_BOOTSTRAP_URL" -o golang.tar.gz \
	&& echo "$GOLANG_BOOTSTRAP_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C "$GOLANG_BOOTSTRAP_PATH" -xzf golang.tar.gz \
	&& rm golang.tar.gz

# Get and build Go tip
RUN export GOROOT_BOOTSTRAP=$GOLANG_BOOTSTRAP_PATH/go \
	&& git clone "$GOLANG_SRC_REPO_URL" /usr/local/go \
	&& cd /usr/local/go/src \
	&& ./make.bash \
	&& rm -rf "$GOLANG_BOOTSTRAP_PATH" /usr/local/go/pkg/bootstrap 

# Build Go workspace and environment
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" \
	&& chmod -R 777 "$GOPATH"

WORKDIR $GOPATH
