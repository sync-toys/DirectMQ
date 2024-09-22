#!/bin/bash

echo Building golang SDK test agent
cd ./agents/go && go build -buildvcs=false -gcflags="all=-N -l" -o bin/go-agent && cd ../../

echo Building cpp SDK test agent
cd ./agents/cpp && cmake -B build -S . -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=ON && cmake --build build && cd ../../
