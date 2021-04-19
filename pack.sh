#! /bin/bash
set -x
NAME=go-cpuminer

CGO_ENABLED=0 go build
zip $(NAME)_linux64.zip $NAME

GOOS=darwin GOARCH=amd64 go build
zip $(NAME)_mac.zip $NAME

GOOS=windows GOARCH=amd64 go build
zip $(NAME)_windows64.zip $NAME.exe

GOOS=linux GOARCH=mipsle go build
zip $(NAME)_mipsle.zip $NAME

GOOS=linux GOARCH=arm go build
zip $(NAME)_arm.zip $NAME

GOOS=linux GOARCH=mips go build
zip $(NAME)_mips.zip $NAME

GOOS=windows GOARCH=386 go build
zip $(NAME)_windows32.zip $NAME.exe

GOOS=linux GOARCH=arm64 go build
zip $(NAME)_arm64.zip $NAME

GOOS=linux GOARCH=mips64 go build
zip $(NAME)_mips64.zip $NAME

GOOS=linux GOARCH=mips64le go build
zip $(NAME)_mips64le.zip $NAME
