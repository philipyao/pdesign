#!/bin/bash
export GOPATH=/home/philip/workspace/pdesign; export GOMAXPROCS=1; go test --bench=. base/log
