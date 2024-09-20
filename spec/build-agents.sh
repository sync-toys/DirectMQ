#!/bin/bash

echo Building golang SDK test agent
cd ./agents/go && go build -buildvcs=false -gcflags="all=-N -l" -o bin/go-agent && cd ../../
