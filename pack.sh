#! /bin/bash
#set -x
NAME="go-cpuminer"

export GO111MODULE=on

#go tool dist list
build_list=$(go tool dist list)

go mod tidy

rm pack -rf
rm pack.zip -f
mkdir pack

BUILD_TIME="${BUILD_TIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"
GIT_COMMIT="${GIT_COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}"
GIT_BRANCH="${GIT_BRANCH:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")}"

LDFLAGS="-s -w -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.GitBranch=${GIT_BRANCH}'"

for line in $build_list; do
  os=$(echo "$line" | awk -F"/" '{print $1}')
  arch=$(echo "$line" | awk -F"/" '{print $2}')
  echo "os="$os" arch="$arch" start build"
  if [ "$os" == "android" ]; then
    continue
  fi
  if [ "$os" == "ios" ]; then
    continue
  fi
  if [ "$arch" == "wasm" ]; then
    continue
  fi
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="${LDFLAGS}"
  if [ $? -ne 0 ]; then
    echo "os="$os" arch="$arch" build fail"
    exit 1
  fi
  if [ "$os" = "windows" ]; then
    zip ${NAME}_"${os}"_"${arch}"".zip" $NAME".exe"
    if [ $? -ne 0 ]; then
      echo "os="$os" arch="$arch" zip fail"
      exit 1
    fi
    mv ${NAME}_"${os}"_"${arch}"".zip" pack/
    rm $NAME".exe" -f
  else
    zip ${NAME}_"${os}"_"${arch}"".zip" $NAME
    if [ $? -ne 0 ]; then
      echo "os="$os" arch="$arch" zip fail"
      exit 1
    fi
    mv ${NAME}_"${os}"_"${arch}"".zip" pack/
    rm $NAME -f
  fi
  echo "os="$os" arch="$arch" done build"
done

zip pack.zip pack/ -r

echo "all done"
