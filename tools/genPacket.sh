#!/usr/bin/env bash
python parse.py ../src/proto/
go fmt packet_name.go
go fmt packet_decode.go
mv packet_name.go ../src/proto/
mv packet_decode.go ../src/proto/
