#! /bin/bash
set -x

CGO_ENABLED=0 go build
zip $NAME_linux64.zip $NAME

GOOS=darwin GOARCH=amd64 go build
zip $NAME_mac.zip $NAME

GOOS=windows GOARCH=amd64 go build
zip $NAME_windows64.zip $NAME.exe

GOOS=linux GOARCH=mipsle go build
zip $NAME_mipsle.zip $NAME

GOOS=linux GOARCH=arm go build
zip $NAME_arm.zip $NAME

GOOS=linux GOARCH=mips go build
zip $NAME_mips.zip $NAME

GOOS=windows GOARCH=386 go build
zip $NAME_windows32.zip $NAME.exe

GOOS=linux GOARCH=arm64 go build
zip $NAME_arm64.zip $NAME

GOOS=linux GOARCH=mips64 go build
zip $NAME_mips64.zip $NAME

GOOS=linux GOARCH=mips64le go build
zip $NAME_mips64le.zip $NAME

