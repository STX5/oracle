FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/oracle
COPY . $GOPATH/src/oracle
RUN go build cmd/worker/worker.go

ENTRYPOINT [ "./worker" ]