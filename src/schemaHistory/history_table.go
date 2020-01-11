package schemaHistory

import (
	"database/sql"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"log"
)

func DoesTableExist(tableName string, db *sql.DB) bool {
	var table string

	row := db.QueryRow(fmt.Sprintf("SHOW TABLES LIKE '%s';", tableName))

	err := row.Scan(&table)

	if err != nil {
		return false
	}

	return table == tableName
}

func CreateIfNotExists(tableName string, db *sql.DB) {
	if DoesTableExist(tableName, db) {
		fmt.Printf("%s already exists.", tableName)
		return
	}

	fmt.Printf("%s does not exist, creating so we can keep track of migrations.", tableName)

	createTableSql := heredoc.Docf(`
		CREATE TABLE IF NOT EXISTS %[1]s
		(
			installed_rank INT                                 NOT NULL
				PRIMARY KEY,
			version        VARCHAR(50)                         NULL,
			description    VARCHAR(200)                        NOT NULL,
			type           VARCHAR(20)                         NOT NULL,
			script         VARCHAR(1000)                       NOT NULL,
			checksum       INT                                 NULL,
			installed_by   VARCHAR(100)                        NOT NULL,
			installed_on   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			execution_time INT                                 NOT NULL,
			success        TINYINT(1)                          NOT NULL
		);
		
		CREATE INDEX %[1]s_s_idx
			ON %[1]s (success);
	`, tableName)

	_, err := db.Exec(createTableSql)

	if err != nil {
		log.Fatalf("Error creating schema history table: %v", err)
	}
}

func RetrieveHistory(tableName string, db *sql.DB) {

}
