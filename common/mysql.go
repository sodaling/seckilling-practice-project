package common

import (
	"database/sql"
	"fmt"
	"seckilling-practice-project/configs"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var defaultDb *sql.DB
var defaultDBlock sync.RWMutex
var _mysqlDsnTemp = "%s:%s@tcp(%s:3306)/miaosha?charset=utf8"

func NewMysqlConn() (*sql.DB, error) {
	config := configs.Cfg
	dsn := fmt.Sprintf(_mysqlDsnTemp, config.Database.Username, config.Database.Password, config.Server.Host)
	db, err := sql.Open("mysql", dsn)
	return db, err
}

func DefaultDb() (*sql.DB, error) {
	defaultDBlock.RLock()
	if defaultDb != nil {
		defaultDBlock.RUnlock()
		return defaultDb, nil
	}
	defaultDBlock.RUnlock()
	defaultDBlock.Lock()
	var err error
	if defaultDb == nil {
		defaultDb, err = NewMysqlConn()
	}
	defaultDBlock.Unlock()
	if err != nil {
		return &sql.DB{}, err
	}
	return defaultDb, nil
}

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
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	columns, _ := rows.Columns()
	vals := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))

	for i := range scans {
		scans[i] = &vals[i]
	}
	i := 0
	result := make(map[int]map[string]string)

	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			row[columns[k]] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
