FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/dao-server
COPY . $GOPATH/src/dao-server
RUN go build .

EXPOSE 9000
ENTRYPOINT ["./dao-server"]
