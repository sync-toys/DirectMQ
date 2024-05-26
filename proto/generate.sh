#!/bin/bash

rm -rf /sdk/go/protocol
rm -rf /sdk/cpp/protocol
# rm -rf /sdk/ts/protocol

buf lint --path ./directmq
buf generate --path ./directmq

mkdir -p /sdk/go/protocol
mkdir -p /sdk/cpp/protocol
# mkdir -p /sdk/ts/protocol
