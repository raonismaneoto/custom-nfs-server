FROM golang:1.18 AS builder

WORKDIR /go/src/nfs

COPY ./nfs-server ./

WORKDIR /go/src/nfs

RUN go install -v

FROM ubuntu:latest

ARG PORT_ARG=5000

COPY --from=builder /go/bin /usr/local/bin

RUN echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc

WORKDIR /usr/local/bin

EXPOSE $PORT_ARG

CMD ("nfs-server 1>logs 2> output.err")