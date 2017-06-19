package oracledb

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	_ "github.com/mattn/go-oci8" // Oracle OCI driver
	"github.com/pkg/errors"
)

// OraDB ...
type OraDB struct {
	Db      *sql.DB
	Version string
}

func init() {
	// Register the ModuleFactory function for the "oracledb" module.
	if err := mb.Registry.AddModule("oracledb", NewModule); err != nil {
		panic(err)
	}
}

// NewModule creates a new mb.Module instance
func NewModule(base mb.BaseModule) (mb.Module, error) {
	// Validate that at least one host has been specified.
	config := struct {
		Hosts []string `config:"hosts"    validate:"nonzero,required"`
	}{}
	if err := base.UnpackConfig(&config); err != nil {
		return nil, err
	}

	return &base, nil
}

// NewDB returns a new oracle database handle. The dsn value (data source name)
// must be valid, otherwise an error will be returned.
//
//   DSN Format: username/password@host:port/service_name
func NewDB(ociURL string) (OraDB, error) {
	var (
		err error
		odb OraDB
	)

	// NLS_LANG is set to American format. At least NLS_NUMERIC_CHARACTERS has to be ".,".
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	os.Setenv("NLS_DATE_FORMAT", "YYYY-MM-DD\"T\"HH24:MI:SS")

	// Open DB connection
	odb.Db, err = sql.Open("oci8", ociURL)
	if err != nil {
		return odb, errors.Wrap(err, "sql open failed")
	}

	// Get DB version
	var version []map[string]interface{}
	version, err = ProcessMetric(odb.Db, "SELECT version FROM v$instance")
	if err != nil {
		return odb, errors.Wrap(err, "get Database Version failed.")
	}
	odb.Version = version[0]["VERSION"].(string)

	return odb, nil
}

// ProcessMetric pases database rows and returns a map
func ProcessMetric(db *sql.DB, sql string) ([]map[string]interface{}, error) {
	var resultMap []map[string]interface{}
	// Run query
	queryResult, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer queryResult.Close()

	// pointer on data of a database column
	rowData := map[string]*interface{}{}
	// holds data of a database column
	rowVars := []interface{}{}
	// holds all column names (always upper case, oracle related)
	colNames, err := queryResult.Columns()
	// bring poniters in place
	for _, col := range colNames {
		rowData[col] = new(interface{})
		rowVars = append(rowVars, rowData[col])
	}

	// loop through database result set
	for queryResult.Next() {
		if err := queryResult.Scan(rowVars...); err != nil {
			return nil, err
		}

		// Define header of MapStr for beat message
		row := map[string]interface{}{
		// "@timestamp": common.Time(time.Now()).String(),
		}
		// parse row
		rowMap, err := parseRow(rowData, row)
		if err != nil {
			logp.Err("Error during parseRow function:", err)
		}
		resultMap = append(resultMap, rowMap)
	}

	return resultMap, nil
}

func parseRow(rowData map[string]*interface{}, row map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range rowData {
		if v == nil {
			continue
		}

		// Schema conversion seams to need only strings... so no conversion to correct datatypes is done.
		switch val := (*v).(type) {
		case string:
			// // Oracle's NUMBER seems not to be recognized correctly.
			// // Therefore we convert int and float back to it's real type...
			// if vv, err := strconv.ParseInt(val, 10, 64); err == nil {
			// 	fmt.Println("string-int", vv, k, *v)
			// 	row[k] = vv
			// } else if vv, err := strconv.ParseFloat(val, 64); err == nil {
			// 	fmt.Println("string-float", k, *v)
			// 	row[k] = vv
			// } else {
			// 	fmt.Println("string", k, *v)
			row[k] = val
			// }
		case int64, int32, int:
			// row[k] = val
			row[k] = val.(string)
			// fmt.Println("int", k, *v)
		case []byte:
			row[k] = string(val)
			// fmt.Println("byte", k, *v)
		case time.Time:
			// row[k] = common.Time(val)
			row[k] = common.Time(val).String()
			// fmt.Println("Time", k, *v)
		default:
			logp.Err("Failed to convert column to golang datatype. Key:", k, "- Value:", val)
		}
	}
	return row, nil
}

// VersionMatch checks if oracle version matches to monitor version
//
// oVer must be between lowVer and highVer. Then 1 is returned, else 0.
func VersionMatch(oVer string, lowVer string, highVer string) int {
	var oVersion []int
	var lowVersion []int
	var highVersion []int

	// Avoid null pointer
	if oVer == "" || lowVer == "" || highVer == "" {
		return 0
	}

	// Split version string at "."
	for _, i := range strings.Split(oVer, ".") {
		j, _ := strconv.Atoi(i)
		oVersion = append(oVersion, j)
	}
	for _, i := range strings.Split(lowVer, ".") {
		j, _ := strconv.Atoi(i)
		lowVersion = append(lowVersion, j)
	}
	for _, i := range strings.Split(highVer, ".") {
		j, _ := strconv.Atoi(i)
		highVersion = append(highVersion, j)
	}

	// is lowVersion same or lower than Oracle Version?
	lVerMatch := 0
	if lowVersion[0] < oVersion[0] {
		lVerMatch = 1
	} else if lowVersion[0] == oVersion[0] {
		if lowVersion[1] <= oVersion[1] {
			lVerMatch = 1
		}
	}

	// is highVersion same or higher than Oracle Version?
	hVerMatch := 0
	if highVersion[0] > oVersion[0] {
		hVerMatch = 1
	} else if highVersion[0] == oVersion[0] {
		if highVersion[1] >= oVersion[1] {
			hVerMatch = 1
		}
	}

	// If both match, than this metric is the right one. return 1
	if lVerMatch == 1 && hVerMatch == 1 {
		return 1
	}
	return 0
}
