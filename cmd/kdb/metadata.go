package main

import (
  "fmt"
  "io"
  "reflect"
  "strings"
)

type Field struct {
  Name      string
  CleanName string
  Type      reflect.Type
}

type Struct struct {
  Name      string
  CleanName string
  Fields    []Field
}

type Code []string

func (c *Code) Append(args ...string) {
  *c = append(*c, args...)
}

func (c *Code) Appendf(format string, args ...interface{}) {
  *c = append(*c, fmt.Sprintf(format, args...))
}

type Metadata struct {
  Args    []string
  Package string
  Structs []Struct

  // filled by Create()
  ImportCode  Code
  InsertStmts Code
  SelectStmts Code
  StructCode  Code
}

func (md *Metadata) Create() *Metadata {
  if *types == "null" {
    md.ImportCode.Append("import \"database/sql\"\n")
  }

  for _, strct := range md.Structs {
    var args []string
    var sqlArgs []string

    // Output the table struct
    md.StructCode.Appendf("type %s struct {", strct.CleanName)
    for _, field := range strct.Fields {
      tag := ""
      if *sqlstruct {
        tag = "`sql:\"" + field.Name + "\"`"
      }

      typ := field.Type.String()
      if !strings.HasPrefix(typ, "[]") {
        switch *types {
        case "null":
          typ = "sql.Null" + capitalize(typ)
        case "pointer":
          typ = "*" + typ
        }
      }

      md.StructCode.Appendf("%s %s%s;", field.CleanName, typ, tag)

      sqlArgs = append(sqlArgs, field.CleanName)
      args = append(args, "&t."+field.CleanName)
    }
    md.StructCode.Appendf("};")

    // Args function
    md.StructCode.Appendf("func (t *%s) Args() []interface{} {return []interface{}{%s}};", strct.CleanName, strings.Join(args, ","))

    // Insert statements
    insertValues := strings.Repeat("?, ", len(sqlArgs))
    md.InsertStmts.Appendf(
      `"INSERT INTO %s (%s) VALUES (%s)"`,
      strct.CleanName,
      strings.Join(sqlArgs, ","),
      insertValues[:len(insertValues)-2])

    // Select statements
    md.SelectStmts.Appendf(
      `"SELECT %s FROM %s"`,
      strings.Join(sqlArgs, ","),
      strct.CleanName)
  }

  return md
}

func (md *Metadata) Output(w io.Writer) {
  fmt.Fprint(w, "// GENERATED BY dbtogo (github.com/kdar/dbtogo); DO NOT EDIT\n")
  fmt.Fprintf(w, "// ---args: %s\n", strings.Join(md.Args, " "))
  fmt.Fprintf(w, "package %s\n", md.Package)

  fmt.Fprint(w, strings.Join(md.ImportCode, ""))

  fmt.Fprintf(w, "type Arger interface {Args() []interface{}};")

  fmt.Fprint(w, "var InsertStmts = map[string]string{\n")
  for n, i := range md.InsertStmts {
    fmt.Fprintf(w, "\"%s\": %s,\n", md.Structs[n].CleanName, i)
  }
  fmt.Fprint(w, "}\n")

  fmt.Fprint(w, "var SelectStmts = map[string]string{\n")
  for n, i := range md.SelectStmts {
    fmt.Fprintf(w, "\"%s\": %s,\n", md.Structs[n].CleanName, i)
  }
  fmt.Fprint(w, "}\n")

  fmt.Fprint(w, strings.Join(md.StructCode, ""))
}
