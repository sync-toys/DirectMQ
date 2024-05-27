FROM ghcr.io/xtruder/nix-devcontainer:latest

WORKDIR /workspace

COPY --chown=1000:1000 ./flake.* ./
RUN nix develop

COPY --chown=1000:1000 ./ /workspace/
