package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func main() {
	a, b := "delete from abc", "ab"
	sq := squirrel.StatementBuilder.
		Select(" min("+a+") as _from", "max("+b+") as _to").
		// Distinct().
		From("abc").
		OrderBy("a DESC")

	q, _, _ := sq.ToSql()

	fmt.Println(q)
}
