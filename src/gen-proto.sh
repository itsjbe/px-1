#!/bin/bash

protoc -I proto proto/* --go_out=plugins=grpc:./phoenix/message