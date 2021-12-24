#! /bin/sh

CGO_ENABLED=0 GOPROXY=https://goproxy.cn GOOS=linux GOARCH=amd64 go build \
-ldflags="-s -w" -o futong-yw-monitor-center monitor-center/main.go