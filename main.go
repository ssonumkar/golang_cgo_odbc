package main

import (
	"github.com/ssonumkar/test-odbc-snowflake/server/prepare"
)

func main() {
	conn_string := "<connection_string>"
	query := "query"
	prepare.Prepare(conn_string, query)
}
