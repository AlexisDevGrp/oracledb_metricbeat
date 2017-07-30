package fra

import (
	"github.com/elastic/beats/libbeat/common"
	s "github.com/elastic/beats/metricbeat/schema"
	c "github.com/elastic/beats/metricbeat/schema/mapstrstr"
)

var (
	schema = s.Schema{
		"file_type":                 c.Str("FILE_TYPE"),
		"percent_space_used":        c.Float("PERCENT_SPACE_USED"),
		"percent_space_reclaimable": c.Float("PERCENT_SPACE_RECLAIMABLE"),
		"number_of_files":           c.Int("NUMBER_OF_FILES"),
		// "TEST_INT":      c.Int("TESTINT"),
		// "TEST_DATE":     c.Time(time.RFC3339, "TESTDATE"),
		// "test_string": c.Str("testString"),
		// "test_bool":   c.Bool("testBool"),
		// "test_float":  c.Float("testFloat"),
		// "test_obj": s.Object{
		// 	"test_obj_string": c.Str("testObjString"),
		// },
	}
)

func eventMapping(input map[string]interface{}) common.MapStr {
	data, _ := schema.Apply(input)
	return data
}
