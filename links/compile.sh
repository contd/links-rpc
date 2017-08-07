#!/usr/bin/env bash

protoc links.proto --go_out=plugins=grpc:.
