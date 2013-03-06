package helper

import (
  "database/sql"
  "reflect"
)

type Arger interface {
  Args() []interface{}
}

// Querys the database for one row, and sets the data in arger.
// Usage:
//  var account Accounts // implements Arger
//  found, err := Query(db, `select * from Accounts where username = ?`, &account, "kevin")
func Query(db *sql.DB, query string, arger Arger, args ...interface{}) (found bool, err error) {
  rows, err := db.Query(query, args...)
  if err != nil {
    return false, err
  }

  if rows.Next() {
    found = true
    err := rows.Scan(arger.Args()...)
    if err != nil {
      return false, err
    }
  }

  return found, nil
}

// Queries database for many rows, and sets the strcts interface as
// the returned rows.
// Usage:
//  var strcts []Accounts // Each Accounts implements Arger
//  err := helper.QueryAll(db, `select * from Accounts`, &strcts, reflect.TypeOf(Accounts{}))
func QueryAll(db *sql.DB, query string, strcts interface{}, typ reflect.Type, args ...interface{}) error {
  rows, err := db.Query(query, args...)
  if err != nil {
    return err
  }

  vof := reflect.ValueOf(strcts)

  for rows.Next() {
    strct := reflect.New(typ)
    err := rows.Scan(strct.Interface().(Arger).Args()...)
    if err != nil {
      return err
    }

    vof.Elem().Set(reflect.Append(vof.Elem(), reflect.Indirect(strct)))
  }

  return nil
}
