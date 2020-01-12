package schemaHistory

import (
	"github.com/hashicorp/go-version"
	"sort"
)

func CurrentSchemaVersion(history *map[string]SchemaHistory) []*version.Version {
	var versionsRaw []string

	for _, value := range *history {
		if value.Version != "" {
			versionsRaw = append(versionsRaw, value.Version)
		}
	}

	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}

	sort.Sort(version.Collection(versions))

	return versions
}

func AnyFailures(history *map[string]SchemaHistory) (bool, SchemaHistory) {
	for _, value := range *history {
		if !value.Success {
			return true, value
		}
	}

	return false, SchemaHistory{}
}
