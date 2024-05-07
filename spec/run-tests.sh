#!/bin/bash

echo Building golang SDK test agent
cd ./agents/go && go build -gcflags="all=-N -l" -o bin/go-agent && cd ../../

echo Running specification tests
cd ./tests && ginkgo -v && cd ../

echo Done
