#!/bin/bash

rm -rf /sdk/go/protocol
rm -rf /sdk/cpp/lib/protocol
# rm -rf /sdk/ts/protocol

buf lint --path ./directmq
buf generate --path ./directmq

mkdir -p /sdk/go/protocol
mkdir -p /sdk/cpp/lib/protocol
# mkdir -p /sdk/ts/protocol
