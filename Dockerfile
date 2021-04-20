FROM golang AS build-env

RUN GO111MODULE=off go get -u github.com/esrrhs/go-cpuminer
RUN GO111MODULE=off go get -u github.com/esrrhs/go-cpuminer/...
RUN GO111MODULE=off go install github.com/esrrhs/go-cpuminer

FROM debian
COPY --from=build-env /go/bin/go-cpuminer .
WORKDIR ./
