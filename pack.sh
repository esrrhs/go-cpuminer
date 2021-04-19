#! /bin/bash
set -x
NAME="go-cpuminer"
rm *.zip -f

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

for os in $os_all; do
  for arch in $arch_all; do
    GOOS=$os GOARCH=$arch go build
    zip ${NAME}_${os}_${arch}".zip" $NAME
  done
done
