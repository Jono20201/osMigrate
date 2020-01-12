package migrations

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type MigrationType string

const (
	VersionedMigration  MigrationType = "V"
	RepeatableMigration MigrationType = "R"
)

type Migration struct {
	Path     string
	Filename string
	Checksum uint32
	Content  string
	Type     MigrationType
}

func LoadRepeatableMigrations(folderPath string) []Migration {
	return LoadMigrations(folderPath, RepeatableMigration)
}

func LoadVersionedMigrations(folderPath string) []Migration {
	return LoadMigrations(folderPath, VersionedMigration)
}

func LoadMigrations(folderPath string, migrationType MigrationType) []Migration {
	var migrations []Migration

	files, err := filepath.Glob(fmt.Sprintf("%s/**/%s*.sql", folderPath, migrationType))

	if err != nil {
		log.Fatalf("Unable to search file system: %v", err)
	}

	for _, path := range files {
		file, err := os.Open(path)

		if err != nil {
			log.Fatalf("Error opening %s: %v", path, err)
		}

		content, err := ioutil.ReadAll(file)

		if err != nil {
			log.Fatalf("Error reading %s: %v", path, err)
		}

		checksum := crc32.ChecksumIEEE(content)

		migrations = append(migrations, Migration{
			Path:     path,
			Filename: filepath.Base(path),
			Checksum: checksum,
			Content:  string(content),
			Type:     migrationType,
		})

		_ = file.Close()
	}

	return migrations
}
