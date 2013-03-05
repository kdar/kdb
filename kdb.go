package kdb

import (
  "database/sql"
  "fmt"
  "github.com/kdar/kmap"
  // "github.com/kisielk/sqlstruct"
  // "reflect"
  "strings"
)

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

  res, err := db.Exec(fmt.Sprintf("INSERT INTO Messages %s VALUES (%s)", fieldsSql, strings.Join(variables, ",")), values...)
  return res, err
}
