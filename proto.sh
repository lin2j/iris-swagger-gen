#!/bin/bash

mkdir -p common/message

#
protoc --proto_path=./proto --go_out=./common/message ./proto/message.proto
#
sed -i "s/,omitempty//g" ./common/message.pb.go

