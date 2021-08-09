package db

import (
	"github.com/hashicorp/go-memdb"
)

type DB struct {
	Conn *memdb.MemDB
}

// InitDB initializes the database connect
func InitDB(schema *memdb.DBSchema) (DB, error) {
	if schema == nil {
		panic("cannot initialize database: missing schema")
	}

	conn, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return DB{conn}, nil
}

// Query queries the data store
func (d *DB) Query(table string, idx string, args ...interface{}) ([]interface{}, error) {
	if d.Conn == nil {
		panic("database is not initialized")
	}

	txn := d.Conn.Txn(false)

	it, err := txn.Get(table, idx, args...)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		results = append(results, obj)
	}

	return results, nil
}

// Upsert inserts or replaces existing data in the data store
func (d *DB) Upsert(table string, record interface{}) error {
	if d.Conn == nil {
		panic("database is not initialized")
	}

	txn := d.Conn.Txn(true)
	defer txn.Abort()

	err := txn.Insert(table, record)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

// Delete deletes rows in the data store
func (d *DB) Delete(table string, idx string, args ...interface{}) (int, error) {
	if d.Conn == nil {
		panic("database is not initialized")
	}

	txn := d.Conn.Txn(true)
	defer txn.Abort()

	count, err := txn.DeleteAll(table, idx, args...)
	if err != nil {
		return 0, err
	}

	txn.Commit()

	return count, nil
}