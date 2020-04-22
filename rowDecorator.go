package qorexGoHelpers

import "database/sql"

// https://stackoverflow.com/questions/50046151/loop-through-database-sql-sql-rows-multiple-times
// Struct that holds a row collection and an OnScan callback func
type RowDecorator struct {
	*sql.Rows
	OnScan func([]interface{}, error)
}

// The Scan method overwrites the default Scan behaviour and applies
// the OnScan callback passing on the params if present
func (rows *RowDecorator) Scan(dest ...interface{}) error {
	err := rows.Rows.Scan(dest...)
	if rows.OnScan != nil {
		rows.OnScan(dest, err)
	}
	return err
}

// Wrap adds a callback func to the RowDecorator
func Wrap(rows *sql.Rows, onScan func([]interface{}, error)) *RowDecorator {
	return &RowDecorator{Rows: rows, OnScan: onScan}
}