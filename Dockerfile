FROM ghcr.io/xtruder/nix-devcontainer:latest

WORKDIR /workspace

RUN apt-get update && apt-get install -y gdb
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Copy the flake.nix and flake.lock files and install the dependencies
COPY --chown=1000:1000 ./flake.* ./
RUN nix develop

# Clean /workspace directory in preparation to run on-create.sh script
RUN rm -rf /workspace/*
