package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormConn() (*gorm.DB, error) {
	dsn := "root:daqidan123@tcp(127.0.0.1:3306)/imooc?charset=utf8"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetResultRow Retrieve the return value -- fetch one row.
func GetResultRowGorm(db *gorm.DB) map[string]string {
	var columns []string

	rows, err := db.Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println(err)
	}

	for _, columnType := range columnTypes {
		columns = append(columns, columnType.Name())
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	if rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
		}
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v)
			}
		}
	}

	return record
}

// GetResultRows Retrieve the return value -- all
func GetResultRowsGorm(db *gorm.DB) map[int]map[string]string {
	var columns []string

	rows, err := db.Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println(err)
	}

	for _, columnType := range columnTypes {
		columns = append(columns, columnType.Name())
	}

	vals := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))
	for k := range vals {
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
