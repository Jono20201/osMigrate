package schemaHistory

import (
	"database/sql"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"log"
	"time"
)

type SchemaHistory struct {
	InstalledRank uint32
	Version       string
	Description   string
	Type          string
	Script        string
	Checksum      uint32
	InstalledBy   string
	InstalledOn   time.Time
	ExecutionTime int64
	Success       bool
}

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
		fmt.Println(fmt.Sprintf("%s already exists.", tableName))
		return
	}

	fmt.Println(fmt.Sprintf("%s does not exist, creating so we can keep track of migrations.", tableName))

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

func RetrieveHistory(tableName string, db *sql.DB) *map[string]SchemaHistory {
	historyRecords := make(map[string]SchemaHistory)

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))

	if err != nil {
		log.Fatalf("Error retrieving history: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var history SchemaHistory

		err = rows.Scan(
			&history.InstalledRank,
			&history.Version,
			&history.Description,
			&history.Type,
			&history.Script,
			&history.Checksum,
			&history.InstalledBy,
			&history.InstalledOn,
			&history.ExecutionTime,
			&history.Success,
		)

		if err != nil {
			log.Fatalf("Error mapping to history: %v", err)
		}

		historyRecords[history.Script] = history
	}

	return &historyRecords
}
