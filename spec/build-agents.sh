#!/bin/bash

echo Building golang SDK test agent
cd ./agents/go && rm -rf bin && go build -buildvcs=false -gcflags="all=-N -l" -o bin/go-agent && cd ../../

echo Building cpp SDK test agent
cd ./agents/cpp && rm -rf build && cmake -B build -S . -DCMAKE_BUILD_TYPE=Debug && cmake --build build && cd ../../
