FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /prometheus-exporter

CMD ["/prometheus-exporter"]

FROM debian:latest

COPY --from=build /prometheus-exporter /prometheus-exporter
COPY /dashboards /dashboards

ENTRYPOINT ["/prometheus-exporter"]
