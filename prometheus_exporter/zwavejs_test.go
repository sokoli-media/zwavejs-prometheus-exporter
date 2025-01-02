package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

var sLogForTesting = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func getGaugeVecValue(t *testing.T, metric *prometheus.GaugeVec, labels []string) float64 {
	var m = &dto.Metric{}
	if err := metric.WithLabelValues(labels...).Write(m); err != nil {
		t.Fatalf("couldnt get metric with labels: %s", err)
	}
	return m.Gauge.GetValue()
}

func Test_ZWave_Aeotec_MultiSensor7_LastActive(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/lastActive",
		"{\"time\":1711922310802,\"value\":1711922310552}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveNodeLastActive, []string{"bathroom_sensor"})

	assert.Equal(t, float64(1711922310552)/1000, newValue)
}

func Test_ZWave_Aeotec_MultiSensor7_Temperature(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/sensor_multilevel/endpoint_0/Air_temperature",
		"{\"time\":1735855076246,\"value\":25.5}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveSensorTemperature, []string{"bathroom_sensor"})

	assert.Equal(t, 25.5, newValue)
}

func Test_ZWave_Aeotec_MultiSensor7_Humidity(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/sensor_multilevel/endpoint_0/Humidity",
		"{\"time\":1735855076298,\"value\":30}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveSensorHumidity, []string{"bathroom_sensor"})

	assert.Equal(t, 30.0, newValue)
}

func Test_ZWave_Aeotec_MultiSensor7_Illuminance(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/sensor_multilevel/endpoint_0/Illuminance",
		"{\"time\":1735855077071,\"value\":34}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveSensorIlluminance, []string{"bathroom_sensor"})

	assert.Equal(t, 34.0, newValue)
}

func Test_ZWave_Aeotec_MultiSensor7_UnknownSensorReading(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/sensor_multilevel/endpoint_0/ThisMetricDoesntExist",
		"{\"time\":1735855077071,\"value\":34}",
	)
	require.NoError(t, err)
}
