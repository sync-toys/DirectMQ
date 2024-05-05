#!/bin/bash

rm -rf /sdk/go/protocol
# rm -rf /sdk/cpp/protocol
# rm -rf /sdk/ts/protocol

buf lint src
buf generate src

mkdir -p /sdk/go/protocol
# mkdir -p /sdk/cpp/protocol
# mkdir -p /sdk/ts/protocol
