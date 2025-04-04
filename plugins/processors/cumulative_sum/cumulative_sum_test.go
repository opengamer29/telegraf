package cumulative_sum

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"testing"
	"time"

	"github.com/influxdata/telegraf/testutil"
)

// TestCumulativeSum perform sum of two metrics
func TestCumulativeSum(t *testing.T) {
	expected := []telegraf.Metric{
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_sum": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_sum": float64(4)},
			time.Unix(0, 0),
		),
	}

	plugin := NewCumulativeSum()
	plugin.Init()

	actual := plugin.Apply(
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": 1},
			time.Unix(0, 0),
		), metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": 3},
			time.Unix(0, 0),
		))
	testutil.RequireMetricsEqual(t, expected, actual)
}

// TestCumulativeSum perform sum of two metrics and left original field
func TestCumulativeSumDropOriginalFalse(t *testing.T) {
	expected := []telegraf.Metric{
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": int64(1), "value_sum": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": int64(3), "value_sum": float64(4)},
			time.Unix(0, 0),
		),
	}

	plugin := NewCumulativeSum()
	plugin.DropOriginalField = false
	plugin.Init()

	actual := plugin.Apply(
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": 1},
			time.Unix(0, 0),
		), metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value": 3},
			time.Unix(0, 0),
		))
	testutil.RequireMetricsEqual(t, expected, actual)
}

// TestCumulativeSum perform sum of two metrics and don't touch string field
func TestCumulativeSumStringField(t *testing.T) {
	expected := []telegraf.Metric{
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(4)},
			time.Unix(0, 0),
		),
	}

	plugin := NewCumulativeSum()
	plugin.Init()

	actual := plugin.Apply(
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(1)},
			time.Unix(0, 0),
		), metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(3)},
			time.Unix(0, 0),
		))
	testutil.RequireMetricsEqual(t, expected, actual)
}

// TestCumulativeSum don't perform sum of two metrics with filtered out fields
func TestCumulativeFieldFilteredOut(t *testing.T) {
	expected := []telegraf.Metric{
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(3)},
			time.Unix(0, 0),
		),
	}

	plugin := NewCumulativeSum()
	plugin.Fields = []string{"another_name"}
	plugin.Init()

	// same as expected
	actual := plugin.Apply(
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value": float64(3)},
			time.Unix(0, 0),
		))
	testutil.RequireMetricsEqual(t, expected, actual)
}

// TestCumulativeSum perform sum of two metrics when field name match config
func TestCumulativeFieldMatch(t *testing.T) {
	expected := []telegraf.Metric{
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(4)},
			time.Unix(0, 0),
		),
	}

	plugin := NewCumulativeSum()
	plugin.Fields = []string{"value"}
	plugin.Init()

	actual := plugin.Apply(
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(1)},
			time.Unix(0, 0),
		),
		metric.New(
			"m1",
			map[string]string{"metric_tag": "from_metric"},
			map[string]interface{}{"value_name": "name", "value_sum": float64(4)},
			time.Unix(0, 0),
		))
	testutil.RequireMetricsEqual(t, expected, actual)
}
