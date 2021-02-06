#!/bin/bash

#CUR=$(pwd)
SOURCE_DIR="proto"
TARGET_DIR="swagger/json"
# 可以使用 -I 参数指定 GOOGLE_API 的目录
# https://github.com/googleapis/googleapis
# GOOGLE_API="googleapis-common-protos-1_3_1"
# 使用 yml 配置接口，可以不用引入 GOOGLE_API
# allow_delete_body=true  可以让 delete 请求支持从请求体获取数据
OUT_OPTS="allow_delete_body=true,logtostderr=true,grpc_api_configuration=$SOURCE_DIR/service.yml:$TARGET_DIR"

echo "protoc -I $SOURCE_DIR --swagger_out=$OUT_OPTS $SOURCE_DIR/service.proto"
protoc -I "$SOURCE_DIR" --swagger_out="$OUT_OPTS" "$SOURCE_DIR/service.proto"

VERSION=$(cat version)
sed -i "s/version not set/$VERSION/g" swagger/json/service.swagger.json
