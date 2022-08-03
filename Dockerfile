FROM golang:1.17.5-alpine3.15 as build
WORKDIR /src/prometheus-kafka-adapter

COPY go.mod go.sum *.go ./

ADD . /src/prometheus-kafka-adapter

RUN go build -o /prometheus-kafka-adapter

FROM alpine:3.15

COPY schemas/metric.avsc /schemas/metric.avsc
COPY --from=build /prometheus-kafka-adapter /

CMD /prometheus-kafka-adapter
