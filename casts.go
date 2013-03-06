// Contains casts to help in converting from pointers
// to non-pointers.
package kdb

func String(v *string) string {
  if v != nil {
    return *v
  }
  return ""
}

func Int64(v *int64) int64 {
  if v != nil {
    return *v
  }
  return int64(0)
}

func Uint64(v *uint64) uint64 {
  if v != nil {
    return *v
  }
  return uint64(0)
}

func Float64(v *float64) float64 {
  if v != nil {
    return *v
  }
  return float64(0)
}

func Byteslice(v *[]byte) []byte {
  if v != nil {
    return *v
  }
  return []byte{}
}

func Uint8slice(v *[]uint8) []uint8 {
  if v != nil {
    return *v
  }
  return []uint8{}
}

func Stringp(v string) *string {
  return &v
}

func Int64p(v int64) *int64 {
  return &v
}

func Uint64p(v uint64) *uint64 {
  return &v
}

func Float64p(v float64) *float64 {
  return &v
}

func Byteslicep(v []byte) *[]byte {
  return &v
}

func Uint8slicep(v []uint8) *[]uint8 {
  return &v
}
