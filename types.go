package helper

import (
  "database/sql/driver"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var s NullString
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  if s.Valid {
//     // use s.String
//  } else {
//     // NULL value
//  }
//
type NullString struct {
  String string
  Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
  if value == nil {
    ns.String, ns.Valid = "", false
    return nil
  }
  ns.Valid = true
  return ConvertAssign(&ns.String, value)
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
  if !ns.Valid {
    return nil, nil
  }
  return ns.String, nil
}

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt64 struct {
  Int64 int64
  Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
  if value == nil {
    n.Int64, n.Valid = 0, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Int64, value)
}

// Value implements the driver Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Int64, nil
}

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat64 struct {
  Float64 float64
  Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
  if value == nil {
    n.Float64, n.Valid = 0, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Float64, value)
}

// Value implements the driver Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Float64, nil
}

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBool struct {
  Bool  bool
  Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
  if value == nil {
    n.Bool, n.Valid = false, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Bool, value)
}

// Value implements the driver Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Bool, nil
}

type NullUint64 struct {
  Uint64 uint64
  Valid  bool
}

func (n *NullUint64) Scan(value interface{}) error {
  if value == nil {
    n.Uint64, n.Valid = 0, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Uint64, value)
}

func (n NullUint64) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Uint64, nil
}

type NullByteslice struct {
  Byteslice []byte
  Valid     bool
}

func (n *NullByteslice) Scan(value interface{}) error {
  if value == nil {
    n.Byteslice, n.Valid = []byte{}, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Byteslice, value)
}

func (n NullByteslice) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Byteslice, nil
}

type NullUint8slice struct {
  Uint8slice []byte
  Valid      bool
}

func (n *NullUint8slice) Scan(value interface{}) error {
  if value == nil {
    n.Uint8slice, n.Valid = []byte{}, false
    return nil
  }
  n.Valid = true
  return ConvertAssign(&n.Uint8slice, value)
}

func (n NullUint8slice) Value() (driver.Value, error) {
  if !n.Valid {
    return nil, nil
  }
  return n.Uint8slice, nil
}
