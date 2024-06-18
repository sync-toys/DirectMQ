#!/bin/bash

rm -rf /sdk/go/protocol
rm -rf /sdk/cpp/lib/dmq-proto-embedded
# rm -rf /sdk/ts/protocol

mkdir -p /sdk/go/protocol
mkdir -p /sdk/cpp/lib/dmq-proto-embedded
# mkdir -p /sdk/ts/protocol

buf lint --path ./directmq
buf generate --path ./directmq
cp CMakeLists.txt /sdk/cpp/lib/dmq-proto-embedded
