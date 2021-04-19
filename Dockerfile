FROM golang AS build-env

RUN go get -u github.com/esrrhs/go-cpuminer
RUN go get -u github.com/esrrhs/go-cpuminer/...
RUN go install github.com/esrrhs/go-cpuminer

FROM debian
COPY --from=build-env /go/bin/go-cpuminer .
WORKDIR ./
