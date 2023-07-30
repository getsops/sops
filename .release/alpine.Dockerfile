FROM alpine:3.18

RUN apk --no-cache add \
      ca-certificates \
      vim \
  && update-ca-certificates

ENV EDITOR vim

# Glob pattern to match the binary for the current architecture
COPY sops* /usr/local/bin/sops

ENTRYPOINT ["sops"]
