package module

import (
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal/github"
)

// Version represents a single version of a module.
type Version struct {
	Version string `json:"version"` // The version number of the provider. Correlates to a tag in the module repository
}

// Metadata represents all the metadata for a module. This includes the list of
// versions available for the module.
type Metadata struct {
	Versions []Version `json:"versions"`
}

type Identifier struct {
	Namespace    string // The module namespace
	Name         string // The module name
	TargetSystem string // The module target system
}

// Module represents a single module.
type Module struct {
	Identifier
	Metadata
	Repository github.Repository // A GitHub client for the module
	Log        *slog.Logger      // A logger for the module
}

func (m *Module) UpdateMetadata() error {
	tags, err := m.Repository.ListTags()
	if err != nil {
		// TODO make this a custom error that the caller can handle
		m.Log.Error("Unable to fetch semver tags, skipping", slog.Any("err", err))
		return nil
	}

	tags = tags.FilterSemver()

	// Merge current versions with new versions
	for _, t := range tags {
		found := false
		for _, v := range m.Metadata.Versions {
			if v.Version == t {
				found = true
				break
			}
		}
		if !found {
			m.Metadata.Versions = append(m.Metadata.Versions, Version{Version: t})
		}
	}

	slices.SortFunc(m.Metadata.Versions, func(a, b Version) int {
		return github.SemverTagSort(a.Version, b.Version)
	})

	return nil
}
