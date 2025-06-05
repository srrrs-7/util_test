FROM docker:28.2.2-dind

# Install required tools
ARG KUBECTL_VERSION="1.29.1"
ARG KIND_VERSION="0.20.0"
ARG HELM_VERSION="3.14.0"

# Install tools with security best practices
RUN apk add --no-cache --update \
    curl \
    ca-certificates \
    && \
    # Install kind
    curl -Lo /usr/local/bin/kind "https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-amd64" && \
    chmod +x /usr/local/bin/kind && \
    # Install kubectl
    curl -Lo /usr/local/bin/kubectl "https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
    chmod +x /usr/local/bin/kubectl && \
    # Install helm
    curl -fsSL https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz | tar -xzO linux-amd64/helm > /usr/local/bin/helm && \
    chmod +x /usr/local/bin/helm && \
    # Verify binaries
    kubectl version --client && \
    kind --version && \
    helm version && \
    # Clean up
    rm -rf /var/cache/apk/*