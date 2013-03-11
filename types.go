package kdb

import (
  "database/sql"
  "database/sql/driver"
  "time"
)

func Nstring(v string) sql.NullString {
  return sql.NullString{v, true}
}

func Nint64(v int64) sql.NullInt64 {
  return sql.NullInt64{v, true}
}

func Nfloat64(v float64) sql.NullFloat64 {
  return sql.NullFloat64{v, true}
}

func Nbool(v bool) sql.NullBool {
  return sql.NullBool{v, true}
}

type NullTime struct {
  Time  time.Time
  Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
  nt.Time, nt.Valid = value.(time.Time)
  return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
  if !nt.Valid {
    return nil, nil
  }
  return nt.Time, nil
}
