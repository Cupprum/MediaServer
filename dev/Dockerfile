FROM fedora:latest

RUN dnf install -y \
    curl \
    unzip \
    jq \
    && dnf clean all

# Download opencode binary directly
RUN curl -fsSL https://opencode.ai/install | bash
ENV PATH="${PATH}:/root/.opencode/bin"

WORKDIR /home/clanker/workspace
