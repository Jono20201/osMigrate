package migrations

import (
	"../schemaHistory"
	"database/sql"
	"fmt"
	"regexp"
	"time"
)

func Apply(db *sql.DB, migration *Migration) (schemaHistory.SchemaHistory, error) {
	history := CreateMigrationHistory(*migration)
	startTime := time.Now()
	_, err := db.Exec(migration.Content)
	elapsed := time.Since(startTime)

	history.ExecutionTime = elapsed.Milliseconds()
	history.Success = err == nil

	if err != nil {
		return history, err
	}

	fmt.Println(fmt.Sprintf("Successfully applied migration \"%s\" in %d ms", migration.Path, elapsed.Milliseconds()))
	return history, nil
}

func CreateMigrationHistory(migration Migration) schemaHistory.SchemaHistory {
	var expression *regexp.Regexp

	if migration.Type == RepeatableMigration {
		expression = regexp.MustCompile(`R__(?P<description>[A-Za-z0-9_]*).*`)
	} else {
		expression = regexp.MustCompile(`V(?P<version>[0-9.]*)__(?P<description>[A-Za-z0-9_]*).*`)
	}

	match := expression.FindStringSubmatch(migration.Filename)
	result := make(map[string]string)
	for i, name := range expression.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	history := schemaHistory.SchemaHistory{
		InstalledRank: 0,
		Version:       result["version"],
		Description:   result["description"],
		Type:          "SQL",
		Script:        migration.Filename,
		Checksum:      migration.Checksum,
		InstalledOn:   time.Now(),
	}

	return history
}
