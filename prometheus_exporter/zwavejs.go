package prometheus_exporter

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"regexp"
	"sync"
	"time"
)

type ZWaveLastActive struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

type FloatSensorReading struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

var zWaveNodeLastActive = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "zwave_node_last_active"}, []string{"node"})
var zWaveSensorTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "zwave_sensor_temperature"}, []string{"sensor"})
var zWaveSensorHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "zwave_sensor_humidity"}, []string{"sensor"})
var zWaveSensorIlluminance = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "zwave_sensor_illuminance"}, []string{"sensor"})

var zWaveLastUpdate = promauto.NewGauge(prometheus.GaugeOpts{Name: "zwave_last_update"})

func processZWaveMessage(logger *slog.Logger, topic string, payload string) error {
	logger.Info("received message", "topic", topic, "payload", payload)

	if found := regexp.MustCompile("^zwave/([^/]+)/lastActive$").FindStringSubmatch(topic); len(found) > 0 {
		lastActive := ZWaveLastActive{}
		err := json.Unmarshal([]byte(payload), &lastActive)
		if err != nil {
			return err
		}

		zWaveNodeLastActive.With(prometheus.Labels{"node": found[1]}).Set(float64(lastActive.Value) / 1000)
	} else if found := regexp.MustCompile("^zwave/([^/]+)/sensor_multilevel/endpoint_0/([^/]+)$").FindStringSubmatch(topic); len(found) > 0 {
		sensorReading := FloatSensorReading{}
		err := json.Unmarshal([]byte(payload), &sensorReading)
		if err != nil {
			return err
		}

		var metric *prometheus.GaugeVec
		switch found[2] {
		case "Air_temperature":
			metric = zWaveSensorTemperature
		case "Humidity":
			metric = zWaveSensorHumidity
		case "Illuminance":
			metric = zWaveSensorIlluminance
		default:
			logger.Info(
				"unknown sensor_multilevel sensor reading",
				"sensor", found[1],
				"metric", found[2],
				"value", sensorReading.Value,
			)
			return nil
		}

		metric.With(prometheus.Labels{"sensor": found[1]}).Set(sensorReading.Value)
	}

	zWaveLastUpdate.SetToCurrentTime()

	return nil
}

func CollectZWaveMetrics(logger *slog.Logger, config MosquittoConfig, wg *sync.WaitGroup, quit chan bool) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientId)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetCleanSession(false)

	choke := make(chan MQTT.Message)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- msg
	})

	logger.Info("connecting to mqtt")
	client := MQTT.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		logger.Error("couldn't connect to mqtt", "error", token.Error())
		return
	}

	logger.Info("subscribing to mqtt topics")
	token = client.Subscribe("zwave/#", byte(2), nil)
	if token.WaitTimeout(5*time.Second) && token.Error() != nil {
		logger.Error("couldn't subscribe to mqtt topics", "error", token.Error())
		return
	}

	logger.Info("waiting for zwave updates on mqtt")
	for {
		select {
		case message := <-choke:
			topic := message.Topic()
			payload := string(message.Payload())

			err := processZWaveMessage(logger, topic, payload)
			if err != nil {
				logger.Error(
					"couldn't process mqtt message",
					"topic",
					message.Topic(),
					"payload",
					string(message.Payload()),
					"error",
					err,
				)
			}
		case <-quit:
			logger.Info("disconnecting from mqtt gracefully")
			client.Disconnect(250)
			wg.Done()
			return
		}
	}
}
