#! /bin/bash
#set -x
NAME="go-cpuminer"
rm *.zip -f

for line in $(go tool dist list); do
  os=$(echo "$line" | awk -F"/" '{print $1}')
  arch=$(echo "$line" | awk -F"/" '{print $2}')
  echo "os="$os" arch="$arch" start build"
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build
  if [ $? -ne 0 ]; then
    echo "os="$os" arch="$arch" build fail"
    exit 1
  fi
  zip ${NAME}_"${os}"_"${arch}"".zip" $NAME
  if [ $? -ne 0 ]; then
    echo "os="$os" arch="$arch" zip fail"
    exit 1
  fi
  echo "os="$os" arch="$arch" done build"
done

echo "all done"

