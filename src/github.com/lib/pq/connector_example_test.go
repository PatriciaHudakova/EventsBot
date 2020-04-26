// +build go1.10

package pq_test

import (
	"database/sql"
	"fmt"
)

func ExampleNewConnector() {
	name := ""
	connector, err := NewConnector(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := sql.OpenDB(connector)
	defer db.Close()

	// Use the DB
	txn, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	txn.Rollback()
}
