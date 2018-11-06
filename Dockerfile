# Build container
FROM golang:alpine as builder 

RUN apk update && apk add git

RUN mkdir -p /go/src/marathon_exporter 
RUN mkdir -p /go/bin
RUN mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH 

COPY . $GOPATH/src/marathon_exporter
WORKDIR $GOPATH/src/marathon_exporter

RUN go build -o marathon_exporter .

# Execution container
FROM golang:alpine

COPY --from=builder /go/src/marathon_exporter/marathon_exporter /marathon_exporter

ENTRYPOINT ["/marathon_exporter"]

EXPOSE 9088
