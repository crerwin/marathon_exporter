FROM golang:alpine as builder 

RUN apk update && apk add git

RUN mkdir /build 

ADD . /build/
WORKDIR /build

RUN go get -d -v

RUN go build -o marathon_exporter

FROM scratch

COPY --from=builder /build/marathon_exporter /marathon_exporter/marathon_exporter
WORKDIR /marathon_exporter

ENTRYPOINT ["./marathon_exporter"]

EXPOSE 9088
