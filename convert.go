package kdb

import (
  "database/sql"
  "database/sql/driver"
  "errors"
  "fmt"
  "reflect"
  "strconv"
  "time"
)

// ConvertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func ConvertAssign(dest, src interface{}) error {
  // Common cases, without reflect.  Fall through.
  switch s := src.(type) {
  case string:
    switch d := dest.(type) {
    case *string:
      *d = s
      return nil
    case *[]byte:
      *d = []byte(s)
      return nil
    }
  case []byte:
    switch d := dest.(type) {
    case *string:
      *d = string(s)
      return nil
    case *interface{}:
      bcopy := make([]byte, len(s))
      copy(bcopy, s)
      *d = bcopy
      return nil
    case *[]byte:
      *d = s
      return nil
    }
  case nil:
    switch d := dest.(type) {
    case *[]byte:
      *d = nil
      return nil
    }
  }

  var sv reflect.Value

  switch d := dest.(type) {
  case *string:
    sv = reflect.ValueOf(src)
    switch sv.Kind() {
    case reflect.Bool,
      reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
      reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
      reflect.Float32, reflect.Float64:
      *d = fmt.Sprintf("%v", src)
      return nil
    }
  case *bool:
    bv, err := driver.Bool.ConvertValue(src)
    if err == nil {
      *d = bv.(bool)
    }
    return err
  case *interface{}:
    *d = src
    return nil
  }

  if scanner, ok := dest.(sql.Scanner); ok {
    return scanner.Scan(src)
  }

  dpv := reflect.ValueOf(dest)
  if dpv.Kind() != reflect.Ptr {
    return errors.New("destination not a pointer")
  }

  if !sv.IsValid() {
    sv = reflect.ValueOf(src)
  }

  dv := reflect.Indirect(dpv)
  if dv.Kind() == sv.Kind() {
    dv.Set(sv)
    return nil
  }

  switch dv.Kind() {
  case reflect.Ptr:
    if src == nil {
      dv.Set(reflect.Zero(dv.Type()))
      return nil
    } else {
      dv.Set(reflect.New(dv.Type().Elem()))
      return ConvertAssign(dv.Interface(), src)
    }
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    s := asString(src)
    i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
    if err != nil {
      return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
    }
    dv.SetInt(i64)
    return nil
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    s := asString(src)
    u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
    if err != nil {
      return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
    }
    dv.SetUint(u64)
    return nil
  case reflect.Float32, reflect.Float64:
    s := asString(src)
    f64, err := strconv.ParseFloat(s, dv.Type().Bits())
    if err != nil {
      return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
    }
    dv.SetFloat(f64)
    return nil
  }

  return fmt.Errorf("unsupported driver -> Scan pair: %T -> %T", src, dest)
}

func asString(src interface{}) string {
  switch v := src.(type) {
  case string:
    return v
  case []byte:
    return string(v)
  }
  return fmt.Sprintf("%v", src)
}

func ScanMapIntoStruct(obj interface{}, objMap map[string][]byte) error {
  dataStruct := reflect.Indirect(reflect.ValueOf(obj))
  if dataStruct.Kind() != reflect.Struct {
    return errors.New("expected a pointer to a struct")
  }

  for key, data := range objMap {
    structField := dataStruct.FieldByName(key)
    if !structField.CanSet() {
      continue
    }

    var v interface{}

    switch structField.Type().Kind() {
    case reflect.Slice:
      v = data
    case reflect.String:
      v = string(data)
    case reflect.Bool:
      v = string(data) == "1"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
      x, err := strconv.Atoi(string(data))
      if err != nil {
        return errors.New("arg " + key + " as int: " + err.Error())
      }
      v = x
    case reflect.Int64:
      x, err := strconv.ParseInt(string(data), 10, 64)
      if err != nil {
        return errors.New("arg " + key + " as int: " + err.Error())
      }
      v = x
    case reflect.Float32, reflect.Float64:
      x, err := strconv.ParseFloat(string(data), 64)
      if err != nil {
        return errors.New("arg " + key + " as float64: " + err.Error())
      }
      v = x
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
      x, err := strconv.ParseUint(string(data), 10, 64)
      if err != nil {
        return errors.New("arg " + key + " as int: " + err.Error())
      }
      v = x
    //Now only support Time type
    case reflect.Struct:
      x, _ := time.Parse("2006-01-02 15:04:05.000 -0700", string(data))
      v = x
    default:
      return errors.New("unsupported type in Scan: " + reflect.TypeOf(v).String())
    }

    structField.Set(reflect.ValueOf(v))
  }

  return nil
}

func ScanStructIntoMap(obj interface{}) (map[string]interface{}, error) {
  dataStruct := reflect.Indirect(reflect.ValueOf(obj))
  if dataStruct.Kind() != reflect.Struct {
    return nil, errors.New("expected a pointer to a struct")
  }

  dataStructType := dataStruct.Type()

  mapped := make(map[string]interface{})

  for i := 0; i < dataStructType.NumField(); i++ {
    field := dataStructType.Field(i)
    fieldName := field.Name

    value := dataStruct.FieldByName(fieldName).Interface()

    mapped[fieldName] = value
  }

  return mapped, nil
}
