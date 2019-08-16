#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

protoc chat/chatpb/chat.proto --go_out=plugins=grpc:.
