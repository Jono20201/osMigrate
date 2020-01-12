package main

import (
	migrationPkg "./migrations"
	"./schemaHistory"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const VERSION = 0.1

func main() {
	fmt.Println(fmt.Sprintf("osMigrate v%g", VERSION))

	username := flag.String("username", "", "username for connecting to the database")
	password := flag.String("password", "", "password for connecting to the database")
	hostname := flag.String("hostname", "127.0.0.1", "database hostname")
	port := flag.Uint("port", 3306, "database port")
	schema := flag.String("schema", "", "schema/database name to target")
	schemaTable := flag.String("schemaTable", "flyway_schema_history", "table name for storing migration history")

	flag.Parse()

	if *schema == "" {
		log.Fatal("Schema must be provided.")
	}

	connectionParams := ConnectionParams{
		username: *username,
		password: *password,
		hostname: *hostname,
		port:     *port,
		schema:   *schema,
	}

	db := connect(connectionParams)
	defer db.Close()

	schemaHistory.CreateIfNotExists(*schemaTable, db)
	history := schemaHistory.RetrieveHistory(*schemaTable, db)
	currentVersion := schemaHistory.CurrentSchemaVersion(history)
	anyFailures, failedMigration := schemaHistory.AnyFailures(history)
	fmt.Println(fmt.Sprintf("Current schema verion: %v", currentVersion))

	if anyFailures {
		log.Fatalf("Unable to run migrations due to a previous failure with \"%s\". Please resolve manually.", failedMigration.Script)
	}

	versionedMigrations := migrationPkg.LoadVersionedMigrations("./fixtures")
	repeatableMigrations := migrationPkg.LoadRepeatableMigrations("./fixtures")

	for _, migration := range versionedMigrations {
		if _, found := (*history)[migration.Filename]; found {
			continue
		}

		history, err := migrationPkg.Apply(db, &migration)

		schemaHistory.Put(db, *schemaTable, history)

		if err != nil {
			log.Fatalf("Error applying migration: %v", err)
		}
	}

	for _, migration := range repeatableMigrations {
		if existingMigration, found := (*history)[migration.Filename]; found && migration.Checksum == existingMigration.Checksum {
			continue
		}

		history, err := migrationPkg.Apply(db, &migration)

		schemaHistory.Put(db, *schemaTable, history)

		if err != nil {
			log.Fatalf("Error applying migration: %v", err)
		}
	}
}
