package tablespace

import (
	"github.com/elastic/beats/libbeat/common"
	s "github.com/elastic/beats/metricbeat/schema"
	c "github.com/elastic/beats/metricbeat/schema/mapstrstr"
)

var (
	schema = s.Schema{
		"tablespace_name":          c.Str("TABLESPACE_NAME"),
		"block_size":               c.Int("BLOCK_SIZE"),
		"status":                   c.Str("STATUS"),
		"contents":                 c.Str("CONTENTS"),
		"logging":                  c.Str("LOGGING"),
		"force_logging":            c.Str("FORCE_LOGGING"),
		"extent_management":        c.Str("EXTENT_MANAGEMENT"),
		"segment_space_management": c.Str("SEGMENT_SPACE_MANAGEMENT"),
		"bigfile":                  c.Str("BIGFILE"),
		"encrypted":                c.Str("ENCRYPTED"),
		"size_bytes":               c.Int("SIZE_BYTES"),
		"free_bytes":               c.Int("FREE_BYTES"),
		"max_size_bytes":           c.Int("MAX_SIZE_BYTES"),
		"max_free_bytes":           c.Int("MAX_FREE_BYTES"),
	}
)

func eventMapping(input map[string]interface{}) common.MapStr {
	data, _ := schema.Apply(input)
	return data
}
