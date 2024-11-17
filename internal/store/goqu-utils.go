package store

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
)

// PGArrayExpr takes a slice of items and returns a literal expression
// with a Postgres array that can be used by goqu.
// It does this by using goqu natively and then manipulating the string
// to surround it in an array literal.
//
// To use this in a struct, you can do something like:
//
//	type MyModel struct {
//	  Tags   []string        `json:"tags" db:"-"`
//	  DbTags goqu.Expression `json:"-" db:"tags"`
//	}
//	body := `{"tags":["x", "y"]}`
//	m := MyModel{}
//	_ = json.Unmarshal([]byte(body), &m)
//	m.DbTags = goqux.PGArrayExpr(m.Tags)
//	sql, _, _ := goqu.Insert("modeltable").Rows(m).ToSQL()
func PGArrayExpr[T any](arr []T) goqu.Expression {
	if len(arr) == 0 {
		return goqu.L("'{}'")
	}
	lit := goqu.V(arr)
	selectSql, _, err := goqu.From(lit).ToSQL()
	if err != nil {
		panic(err)
	}
	valuesSql := strings.TrimPrefix(selectSql, "SELECT * FROM ")
	if valuesSql == selectSql {
		panic("expected go to output an (invalid) 'SELECT * FROM (x, y, z)' for the slice")
	}
	if valuesSql[0] != '(' || valuesSql[len(valuesSql)-1] != ')' {
		panic("expected goqu to output '(x, y, z)' but is missing parens")
	}
	arraySql := fmt.Sprintf("ARRAY[%s]", valuesSql[1:len(valuesSql)-1])
	return goqu.L(arraySql)
}
