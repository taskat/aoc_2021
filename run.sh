#!/bin/bash
day=$1
part=$2
cd $day
go build -o main.exe main.go
echo $part | ./main.exe
rm -f main.exe