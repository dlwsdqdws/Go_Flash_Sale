package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// NewMysqlConn Create a MySQL connection.
func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:daqidan123@tcp(127.0.0.1:3306)/imooc?charset=utf8")
	return
}

// GetResultRow Retrieve the return value -- fetch one row.
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

// GetResultRows Retrieve the return value -- all
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	// Return all columns.
	columns, _ := rows.Columns()
	// All values of a row's columns are represented using []byte.
	vals := make([][]byte, len(columns))
	// Fill in data for one row.
	scans := make([]interface{}, len(columns))
	// scans references vals to fill in data into a []byte slice.
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := columns[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
