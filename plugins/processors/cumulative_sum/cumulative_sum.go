//go:generate ../../../tools/readme_config_includer/generator
package cumulative_sum

import (
	_ "embed"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

//go:embed sample.conf
var sampleConfig string

type CumulativeSum struct {
	Log          telegraf.Logger
	Fields       []string `toml:"fields"`
	DropOriginal bool     `toml:"drop_original"`

	fields map[string]bool

	// TODO: do we need to clear old caches?
	// TODO: does processor's Apply call concurrently? if so we need to use atomic in cache
	cache map[uint64]aggregate
}

type aggregate struct {
	name   string
	tags   map[string]string
	fields map[string]float64
}

func (*CumulativeSum) SampleConfig() string {
	return sampleConfig
}

func (c *CumulativeSum) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, original := range in {
		id := original.HashID()
		if _, ok := c.cache[id]; !ok {
			a := aggregate{
				name:   original.Name(),
				tags:   original.Tags(),
				fields: make(map[string]float64),
			}
			for _, field := range original.FieldList() {
				if c.fields != nil {
					if _, ok := c.fields[field.Key]; !ok {
						continue
					}
				}
				if fv, ok := convert(field.Value); ok {
					a.fields[field.Key] = fv
					original.AddField(field.Key+"_sum", fv)
					if c.DropOriginal {
						original.RemoveField(field.Key)
					}
				}
			}
		} else {
			for _, field := range original.FieldList() {
				if c.fields != nil {
					if _, ok := c.fields[field.Key]; !ok {
						continue
					}
				}
				if fv, ok := convert(field.Value); ok {
					if _, ok := c.cache[id].fields[field.Key]; !ok {
						// hit an uncached field of a cached metric
						c.cache[id].fields[field.Key] = fv
					} else {
						c.cache[id].fields[field.Key] = c.cache[id].fields[field.Key] + fv
					}
					original.AddField(field.Key+"_sum", c.cache[id].fields[field.Key])
					if c.DropOriginal {
						original.RemoveField(field.Key)
					}
				}
			}
		}
	}
	return in
}

func convert(in interface{}) (float64, bool) {
	switch v := in.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func (c *CumulativeSum) Init() error {
	if c.Fields != nil {
		c.fields = make(map[string]bool, len(c.Fields))
		for _, field := range c.Fields {
			c.fields[field] = true
		}
	}
	return nil
}

func init() {
	processors.Add("cumulative_sum", func() telegraf.Processor {
		return &CumulativeSum{
			DropOriginal: true,
		}
	})
}
