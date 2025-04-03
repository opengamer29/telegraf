//go:build !custom || processors || processors.lookup

package all

import _ "github.com/influxdata/telegraf/plugins/processors/cumulative_sum" // register plugin
