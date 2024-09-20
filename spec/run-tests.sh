#!/bin/bash

/bin/sh ./build-agents.sh

echo Running specification tests
# cd ./tests && ginkgo -v --focus "" && cd ../
cd ./tests && ginkgo && cd ../

echo Done
