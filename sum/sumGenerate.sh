#!/usr/bin/env bash

protoc sumpb/sum.proto --go_out=plugins=grpc:.
