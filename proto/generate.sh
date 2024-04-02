#!/bin/bash

rm -rf gen/cpp/*
rm -rf gen/go/*
rm -rf gen/ts/*

rm -rf /projects/cpp/protocol
rm -rf /projects/go/protocol
rm -rf /projects/ts/protocol

buf lint src
buf generate src

mkdir -p /projects/cpp/protocol/v1
mkdir -p /projects/go/protocol/v1
mkdir -p /projects/ts/protocol/v1

cp -r gen/cpp/directmq/v1/* /projects/cpp/protocol/v1/
cp -r gen/go/protocol/* /projects/go/protocol/v1/
cp -r gen/ts/directmq/v1/* /projects/ts/protocol/v1/
