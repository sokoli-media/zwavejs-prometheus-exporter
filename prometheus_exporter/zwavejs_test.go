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

func Test_ZWave_Aeotec_Meter_TotalConsumption(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/sowa_power_outlet/meter/endpoint_0/value/65537",
		"{\"time\":1735906852206,\"value\":1}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveMeterTotalConsumption, []string{"sowa_power_outlet"})

	assert.Equal(t, 1.0, newValue)
}

func Test_ZWave_Aeotec_Meter_Power(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/sowa_power_outlet/meter/endpoint_0/value/66049",
		"{\"time\":1735906853203,\"value\":3.395}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveMeterPower, []string{"sowa_power_outlet"})

	assert.Equal(t, 3.395, newValue)
}

func Test_ZWave_Aeotec_Meter_Voltage(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/sowa_power_outlet/meter/endpoint_0/value/66561",
		"{\"time\":1735906854203,\"value\":240.71}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveMeterVoltage, []string{"sowa_power_outlet"})

	assert.Equal(t, 240.71, newValue)
}

func Test_ZWave_Aeotec_Meter_Current(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/sowa_power_outlet/meter/endpoint_0/value/66817",
		"{\"time\":1735906855204,\"value\":0.023}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveMeterCurrent, []string{"sowa_power_outlet"})

	assert.Equal(t, 0.023, newValue)
}
