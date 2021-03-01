FROM golang:1.15 as build

WORKDIR /go/src/json-exporter
COPY . .
RUN go get -v
RUN go test -v
RUN CGO_ENABLED=0 go build -v

FROM alpine:latest

COPY --from=build /go/src/json-exporter/json-exporter /json-exporter
RUN chmod +x /json-exporter

EXPOSE 9116
ENTRYPOINT [ "/json-exporter" ]
