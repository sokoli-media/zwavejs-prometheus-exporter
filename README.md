# zwavejs-prometheus-exporter

Simple project to export ZWave metrics to Prometheus.

Docker image: `sokolimedia/zwavejs-prometheus-exporter:latest`

Project exports http api on `:9000` with metrics at `/metrics` url.

Required environmental variables:
 * `MOSQUITTO_BROKER`: host of your Mosquitto instance
 * `MOSQUITTO_CLIENT_ID`: customer's client id (set to something unique, e.g. name of this project)
 * `MOSQUITTO_USERNAME`: Mosquitto username
 * `MOSQUITTO_PASSWORD`: Mosquitto password

Prerequisites:
 * zwavejs installed
 * mosquitto installed
 * mqtt enabled
