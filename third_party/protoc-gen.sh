#!/bin/bash

protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 star-wars-service.proto
protoc --proto_path=pkg/swapi/proto --proto_path=third_party --go_out=plugins=grpc:pkg/swapi swapi-response.proto