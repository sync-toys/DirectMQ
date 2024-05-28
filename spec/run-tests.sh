#!/bin/bash

echo Building golang SDK test agent
cd ./agents/go && go build -buildvcs=false -gcflags="all=-N -l" -o bin/go-agent && cd ../../

echo Running specification tests
# cd ./tests && ginkgo -v --focus "" && cd ../
cd ./tests && ginkgo && cd ../

echo Done
