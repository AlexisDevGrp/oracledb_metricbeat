package status

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	s "github.com/elastic/beats/metricbeat/schema"
	c "github.com/elastic/beats/metricbeat/schema/mapstrstr"
)

var (
	schema = s.Schema{
		"inst_id":       c.Int("INST_ID"),
		"host_name":     c.Str("HOST_NAME"),
		"version":       c.Str("VERSION"),
		"instance_name": c.Str("INSTANCE_NAME"),
		"startup_time":  c.Time(time.RFC3339, "STARTUP_TIME"),
		"status":        c.Str("STATUS"),
		"parallel":      c.Str("PARALLEL"),
		"thread_no":     c.Str("THREAD_NO"),
		"archiver":      c.Str("ARCHIVER"),
		"instance_role": c.Str("INSTANCE_ROLE"),
		"con_id":        c.Int("CON_ID"),
		"edition":       c.Str("EDITION"),
		"database_type": c.Str("DATABASE_TYPE"),
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
