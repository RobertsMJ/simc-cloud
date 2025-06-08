# Use Ubuntu as the base image
FROM ubuntu:22.04

# Avoid interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install basic dependencies
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    gnupg \
    lsb-release \
    software-properties-common

# Copy scripts directory and run install-docker.sh to install Docker
COPY scripts/ /scripts/
RUN bash /scripts/install-docker.sh

# Install Go (latest stable)
RUN curl -fsSL https://go.dev/dl/ | grep linux-amd64.tar.gz | head -n 1 | \
    awk -F'"' '{print $4}' | xargs -I{} curl -fsSL -o /tmp/go.tar.gz https://go.dev{} && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"


# Install Node.js (latest LTS)
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && \
    apt-get install -y nodejs

# Install Terraform (latest)
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
    tee /etc/apt/sources.list.d/hashicorp.list && \
    apt-get update && apt-get install -y terraform

# Install AWS SAM CLI dependencies
RUN apt-get update && apt-get install -y \
    unzip \
    python3 \
    python3-pip

# Install AWS SAM CLI (latest stable release)
RUN pip3 install aws-sam-cli

# Clean up
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Set default shell
SHELL ["/bin/bash", "-c"]
