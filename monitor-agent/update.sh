#!/bin/bash
#set -e
 pushd /usr/local/bin/ && \
  curl -O http://172.16.71.20:8001/packages/ft-agent && \
  chmod 777 ft-agent && \
   curl -O http://172.16.71.20:8001/packages/config.yaml && \
   chmod 755 config.yaml && \
 nohup ./ft-agent -w=false -i=false -o=true &


 pushd /usr/local/bin/ && \
  curl -O http://172.16.8.145:8001/packages/ft-agent && \
  chmod 777 ft-agent && \
   curl -O http://172.16.8.145:8001/packages/config.yaml && \
   chmod 755 config.yaml && \
 nohup ./ft-agent -w=false -i=true -o=false &