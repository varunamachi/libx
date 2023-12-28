package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func main() {
	sq := squirrel.StatementBuilder.
		Select("a", "b").
		Distinct().
		From("abc").
		OrderBy("a DESC")

	q, _, _ := sq.ToSql()

	fmt.Println(q)
}
