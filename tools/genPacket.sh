#!/usr/bin/env bash
python parse.py ./
go fmt packet_name.go
mv packet_name.go ../src/proto/
