FROM debian:bookworm-slim

RUN apt-get update && apt-get install --no-install-recommends -y \
      awscli \
      azure-cli \
      curl \
      gnupg \
      vim \
    && rm -rf /var/lib/apt/lists/*

ENV EDITOR vim

# Glob pattern to match the binary for the current architecture
COPY sops* /usr/local/bin/sops

ENTRYPOINT ["sops"]
