package kdb

import (
  "database/sql"
  "fmt"
  "github.com/kdar/kmap"
  // "github.com/kisielk/sqlstruct"
  // "reflect"
  "reflect"
  "strings"
)

type Arger interface {
  Args() []interface{}
}

func GetMaps(rows *sql.Rows) ([]kmap.Map, error) {
  var maps []kmap.Map

  cols, err := rows.Columns()
  if err != nil {
    return nil, err
  }

  values := make([]interface{}, len(cols))
  scanArgs := make([]interface{}, len(cols))

  for x := 0; x < len(cols); x++ {
    scanArgs[x] = &values[x]
  }

  for rows.Next() {
    m := kmap.Make()
    err := rows.Scan(scanArgs...)
    if err == nil {
      for n, c := range cols {
        m[strings.ToLower(c)] = values[n]
      }
      maps = append(maps, m)
    }
  }

  return maps, nil
}

func QueryMap(db *sql.DB, query string, args ...interface{}) (kmap.Map, error) {
  rows, err := db.Query(query, args...)
  if err != nil {
    return nil, err
  }

  ret, err := GetMaps(rows)
  if len(ret) > 0 {
    return ret[0], err
  }

  return nil, err
}

func QueryMaps(db *sql.DB, query string, args ...interface{}) ([]kmap.Map, error) {
  rows, err := db.Query(query, args...)
  if err != nil {
    return nil, err
  }

  ret, err := GetMaps(rows)
  if len(ret) > 0 {
    return ret, err
  }

  return nil, err
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

// func QueryStruct(db *sql.DB, query string, strct interface{}, args ...interface{}) error {
//   rows, err := db.Query(query, args...)
//   if err != nil {
//     return err
//   }

//   if rows.Next() {
//     err = sqlstruct.Scan(strct, rows)
//     if err != nil {
//       return err
//     }
//   }

//   return nil
// }

// // use it like:
// // var accounts []*Account
// // var account Account
// // QueryStructs(db, "select * from Accounts", &accounts, &account)
// func QueryStructs(db *sql.DB, query string, strcts interface{}, strct interface{}, args ...interface{}) error {
//   rows, err := db.Query(query, args...)
//   if err != nil {
//     return err
//   }

//   vof := reflect.ValueOf(strcts)
//   //into := reflect.New(vof.Type().Elem().Elem())

//   for rows.Next() {
//     err := sqlstruct.Scan(strct, rows)
//     if err != nil {
//       fmt.Println(err)
//       return err
//     }

//     vof.Elem().Set(reflect.Append(vof.Elem(), reflect.ValueOf(strct)))
//   }

//   return nil
// }

func Fields(names []string) string {
  for i, _ := range names {
    names[i] = strings.Replace(names[i], `"`, `\"`, -1)
  }
  return `("` + strings.Join(names, `", "`) + `")`
}

func InsertMap(db *sql.DB, table string, m kmap.Map) (sql.Result, error) {
  var fields []string
  var values []interface{}
  var variables []string

  for key, value := range m {
    fields = append(fields, key)
    values = append(values, value)
    variables = append(variables, "?")
  }

  fieldsSql := Fields(fields)

  res, err := db.Exec(fmt.Sprintf("INSERT INTO %s %s VALUES (%s)", table, fieldsSql, strings.Join(variables, ",")), values...)
  return res, err
}
