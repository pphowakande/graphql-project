package main

import (
	in "api/src/init"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	in.Initialize()
}
