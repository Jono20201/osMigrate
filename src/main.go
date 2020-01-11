package main

import (
	"./schemaHistory"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const VERSION = 1.0

func logError(err error) {
	log.Fatalf("Fatal error: %v", err)
}

func main() {
	fmt.Println(fmt.Sprintf("osMigrate v%g", VERSION))

	username := flag.String("username", "", "username for connecting to the database")
	password := flag.String("password", "", "password for connecting to the database")
	hostname := flag.String("hostname", "127.0.0.1", "database hostname")
	schema := flag.String("schema", "", "schema/database name to target")
	schemaTable := flag.String("schemaTable", "flyway_schema_history", "table name for storing migration history")

	flag.Parse()

	connectionParams := ConnectionParams{
		username: *username,
		password: *password,
		hostname: *hostname,
		schema:   *schema,
	}

	db := connect(connectionParams)
	defer db.Close()

	loadMigrations("./fixtures")
	schemaHistory.CreateIfNotExists(*schemaTable, db)
	schemaHistory.RetrieveHistory(*schemaTable, db)
}
