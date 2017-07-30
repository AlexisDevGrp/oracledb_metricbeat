package status

import (
	"fmt"
	"runtime/debug"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/odbaeu/oracledb_metricbeat/module/oracledb"
	"github.com/pkg/errors"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("oracledb", "status", New); err != nil {
		panic(err)
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	oraDB   oracledb.OraDB
	counter int
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	config := struct{}{}

	logp.Warn("EXPERIMENTAL: The oracledb status metricset is experimental")

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		counter:       1,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	// Recover panic is for deugging only!!!
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic!", err)
			fmt.Println(string(debug.Stack()))
		}
	}()

	// Open Dataase connection
	oraDB, err := oracledb.NewDB(m.HostData().URI)
	if err != nil {
		failEvent := []common.MapStr{{"status": "OFFLINE"}}
		return failEvent, errors.Wrap(err, "oracledb-status open db connection failed")
	}
	defer oraDB.Close()

	// Oracle Database query of this MetricSet
	// Any database version

	// Oracle Database query of this MetricSet
	var qry string
	if oracledb.VersionMatch(oraDB, "0", "11.2") == 1 {
		// 11g
		qry = `SELECT inst_id, instance_name, host_name, version, startup_time, status,
					  parallel, thread#, archiver, instance_role
				FROM  gv$instance`
	} else if oracledb.VersionMatch(oraDB, "12.1", "12.1") == 1 {
		// 12cR1
		qry = `SELECT inst_id, instance_name, host_name, version, startup_time, status,
					  parallel, thread#, archiver, instance_role, con_id, edition
				FROM  gv$instance`
	} else if oracledb.VersionMatch(oraDB, "12.2", "9999") == 1 {
		// 12cR2
		qry = `SELECT inst_id, instance_name, host_name, version, startup_time, status,
					  parallel, thread#, archiver, instance_role, con_id, edition, database_type
				FROM  gv$instance`
	}

	// Load data
	var data []map[string]interface{}
	// data, err = oracledb.ProcessMetric(m.oraDB.Db, qry)
	data, err = oracledb.ProcessMetric(oraDB, qry)
	if err != nil {
		failEvent := []common.MapStr{{"status": "OFFLINE"}}
		m.oraDB.Db = nil
		return failEvent, errors.Wrap(err, "oracledb-status fetch failed")
	}

	// Map data with multiple rows
	event := []common.MapStr{}
	for _, dd := range data {
		event = append(event, eventMapping(dd))
	}
	// fmt.Println("event:", event)

	return event, nil
}
