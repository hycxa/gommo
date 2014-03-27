#!/usr/bin/env bash
python parse.py ../src/proto/

fileNameTab=(
name.go
decode.go
)
for ((i=0; i< ${#fileNameTab[*]}; i++))
do
	go fmt "${fileNameTab[$i]}"
	mv "${fileNameTab[$i]}" ../src/proto/
done

