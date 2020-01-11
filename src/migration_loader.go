package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Migration struct {
	path     string
	checksum uint32
	content  string
}

func loadMigrations(folderPath string) []Migration {
	var migrations []Migration

	files, err := filepath.Glob(fmt.Sprintf("%s/**/V*.sql", folderPath))

	if err != nil {
		log.Fatalf("Unable to search file system: %v", err)
	}

	for _, path := range files {
		file, err := os.Open(path)

		if err != nil {
			logError(err)
		}

		content, err := ioutil.ReadAll(file)

		if err != nil {
			logError(err)
		}

		checksum := crc32.ChecksumIEEE(content)

		migrations = append(migrations, Migration{
			path:     path,
			checksum: checksum,
			content:  string(content),
		})

		_ = file.Close()
	}

	return migrations
}
