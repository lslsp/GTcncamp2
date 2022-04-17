FROM golang:1.18 AS build
ENV GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . /go/src/httpserver
WORKDIR /go/src/httpserver

RUN go build -o /bin/httpserver

FROM alpine
COPY --from=build /bin/httpserver /bin/httpserver
EXPOSE 80
ENTRYPOINT ["bin/httpserver"]