#!/bin/bash

git clone https://github.com/sync-toys/DirectMQ.git /workspace
chown -R 1000:1000 /workspace

# This is fix to prevent docker pull errors
# rm ~/.docker/config.json
# cp /workspace/.devcontainer/docker-config.json ~/.docker/config.json
