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

func Test_ZWave_LastActive(t *testing.T) {
	err := processZWaveMessage(
		sLogForTesting,
		"zwave/bathroom_sensor/lastActive",
		"{\"time\":1711922310802,\"value\":1711922310552}",
	)
	require.NoError(t, err)

	newValue := getGaugeVecValue(t, zWaveNodeLastActive, []string{"bathroom_sensor"})

	assert.Equal(t, float64(1711922310552)/1000, newValue)
}
