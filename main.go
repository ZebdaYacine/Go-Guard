package main

import (
	"go-gaurd/api/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server.InitFibreServer()
}
