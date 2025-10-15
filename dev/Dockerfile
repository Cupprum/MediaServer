FROM fedora:42

RUN dnf update -y && \
    dnf install -y \
        kubectl \
        ssh \
        jq

RUN useradd x42
USER x42

# Create kube config directory
RUN mkdir --parents /home/x42/.kube

WORKDIR /workspace

ENTRYPOINT ["/workspace/dev/bootstrap.sh"]
