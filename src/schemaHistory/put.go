package schemaHistory

import (
	"database/sql"
	"fmt"
	"log"
)

func Put(db *sql.DB, tableName string, history SchemaHistory) {
	stmt, err := db.Prepare(fmt.Sprintf(`INSERT INTO %[1]s
(installed_rank, version, description, type, script, checksum, installed_by, installed_on, execution_time, success)
SELECT IFNULL(MAX(installed_rank), 0) + 1, ?, ?, ?, ?, ?, CURRENT_USER(), ?, ?, ? FROM %[1]s`, tableName))

	if err != nil {
		log.Fatalf("Error creating prepared statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		history.Version,
		history.Description,
		history.Type,
		history.Script,
		history.Checksum,
		history.InstalledOn,
		history.ExecutionTime,
		history.Success,
	)

	if err != nil {
		log.Fatalf("Error running insert: %v", err)
	}
}
