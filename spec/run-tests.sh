#!/bin/bash

task agents:install agents:build

echo Running specification tests
# cd ./tests && ginkgo -v --focus "" && cd ../
cd ./tests && ginkgo && cd ../

echo Done
