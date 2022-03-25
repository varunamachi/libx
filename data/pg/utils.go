package pg

import "strings"

// InsertQuery - creates a basic insert/create query from table name
// and fields. This query should only be used with named exec i.e
// with name parameters. The field names are also used to name the parameters
func InsertQuery(table string, fields ...string) string {
	builder := strings.Builder{}
	builder.WriteString("INSERT INTO ")
	builder.WriteString(table)
	builder.WriteString("(")
	for index, field := range fields {
		builder.WriteString(field)
		if index != len(fields)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")

	}
	builder.WriteString(") VALUES (")
	for index, field := range fields {
		builder.WriteString(":")
		builder.WriteString(field)
		if index != len(fields)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString(")")
	return builder.String()
}
