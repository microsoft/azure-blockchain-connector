FROM golang:1.10 AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build connector.go

FROM ubuntu:16.04 AS ca-store
RUN apt-get update -y
RUN apt-get install -y ca-certificates
RUN update-ca-certificates -v

FROM scratch
COPY --from=ca-store /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/connector /bin/connector
ENTRYPOINT ["connector"]