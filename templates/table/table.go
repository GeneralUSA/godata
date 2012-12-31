package table

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"
)

var t *template.Template

func init() {
	fs := template.FuncMap{"field": reflectField}
	t = template.Must(template.New("").Funcs(fs).Parse(temp))
}

type Table struct {
	Columns map[string]string
	Rows    []interface{}
}

func (table Table) Render() string {
	buffer := new(bytes.Buffer)
	fmt.Println(table.Rows)
	err := t.Execute(buffer, table)
	fmt.Println(err)
	return buffer.String()
}

func reflectField(item interface{}, field string) (interface{}, error) {
	valueOf := reflect.ValueOf(item)
	if valueOf.Kind() != reflect.Struct {
		return "", errors.New(fmt.Sprintf("Cannot take field %v of type %v", field, valueOf.Kind()))
	}
	fieldValue := valueOf.FieldByName(field)
	return fieldValue.Interface(), nil
}

var temp = `
<table class='table'>
  <thead>
    <tr>
      {{ range .Columns }}
        <th>{{ . }}</th>
      {{ end }}
    </tr>
  </thead>
  <tbody>
    {{ range $row := .Rows }}<tr>
      {{ range $name, $label := $.Columns }}<td>{{ field $row $name }}</td>
      {{ end }}
    </tr>{{ end }}
  </tbody>
</table>
`
