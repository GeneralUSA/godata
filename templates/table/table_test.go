package table

import (
	"fmt"
	"testing"
)

type TestRow struct {
	Col1 string
	Col2 string
}

func TestTable(t *testing.T) {
	table := Table{}
	table.Columns = map[string]string{"Col1": "Column 1", "Col2": "Column 2"}
	table.Rows = []interface{}{TestRow{"One", "Two"}}
	fmt.Println(table.Render())
}
