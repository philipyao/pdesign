#!/bin/bash
export GOPATH=/home/philip/workspace/pdesign; export GOMAXPROCS=1; go test --bench=. -benchtime=3s base/log
