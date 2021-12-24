#!/bin/bash
cd  /home/go/futong-yw-monitor-center/monitor-agent

unixTime=`date +%s`
#echo $unixTime
agentName="ft-agent"
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build \
-ldflags="-s -w -X futong-yw-monitor-center/monitor-base/bg.Version=${unixTime}" \
-o $agentName \
main.go

cp $agentName /home/go/futong-yw-monitor-center/monitor-center/packages/
scp  $agentName 172.16.71.31:/etc/ftcloud/yw/futong-yw-monitor-center/packages/

# update db
json="{\"id\":1,\"agentVersion\":${unixTime},\"downloadpath\":\"http://172.16.71.20:8001/packages/ft-agent\"}"
curl -H 'Content-Type: application/json' -X PUT \
    -d ${json} \
http://172.16.71.20:8001/api/agent/updateVer

