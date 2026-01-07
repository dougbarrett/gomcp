// Package metadata handles storing and retrieving scaffold metadata.
// This allows tracking what domains were scaffolded and with what parameters,
// enabling future sync/upgrade capabilities.
package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dbb1dev/go-mcp/internal/types"
)

const (
	// MetadataDir is the directory where metadata files are stored.
	MetadataDir = ".mcp"
	// MetadataFile is the name of the metadata file.
	MetadataFile = "scaffold-metadata.json"
	// CurrentVersion is the current metadata schema version.
	CurrentVersion = "1.0.0"
)

// ProjectMetadata contains all scaffold metadata for a project.
type ProjectMetadata struct {
	Version string                    `json:"version"`
	Domains map[string]DomainMetadata `json:"domains"`
	Wizards map[string]WizardMetadata `json:"wizards,omitempty"`
}

// DomainMetadata contains metadata for a single scaffolded domain.
type DomainMetadata struct {
	ScaffoldedAt      time.Time                 `json:"scaffolded_at"`
	UpdatedAt         time.Time                 `json:"updated_at,omitempty"`
	ScaffolderVersion string                    `json:"scaffolder_version"`
	Input             types.ScaffoldDomainInput `json:"input"`
}

// WizardMetadata contains metadata for a single scaffolded wizard.
type WizardMetadata struct {
	ScaffoldedAt      time.Time                  `json:"scaffolded_at"`
	UpdatedAt         time.Time                  `json:"updated_at,omitempty"`
	ScaffolderVersion string                     `json:"scaffolder_version"`
	Domain            string                     `json:"domain"`
	Input             types.ScaffoldWizardInput  `json:"input"`
}

// Store handles reading and writing scaffold metadata.
type Store struct {
	projectDir string
	mu         sync.RWMutex
}

// NewStore creates a new metadata store for the given project directory.
func NewStore(projectDir string) *Store {
	return &Store{
		projectDir: projectDir,
	}
}

// metadataPath returns the full path to the metadata file.
func (s *Store) metadataPath() string {
	return filepath.Join(s.projectDir, MetadataDir, MetadataFile)
}

// Load reads the project metadata from disk.
// Returns empty metadata if the file doesn't exist.
func (s *Store) Load() (*ProjectMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.metadataPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ProjectMetadata{
				Version: CurrentVersion,
				Domains: make(map[string]DomainMetadata),
				Wizards: make(map[string]WizardMetadata),
			}, nil
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta ProjectMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Ensure maps are initialized
	if meta.Domains == nil {
		meta.Domains = make(map[string]DomainMetadata)
	}
	if meta.Wizards == nil {
		meta.Wizards = make(map[string]WizardMetadata)
	}

	return &meta, nil
}

// Save writes the project metadata to disk.
func (s *Store) Save(meta *ProjectMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Join(s.projectDir, MetadataDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Marshal with indentation for readability
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	path := s.metadataPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// SaveDomain saves or updates metadata for a single domain.
func (s *Store) SaveDomain(domainName string, input types.ScaffoldDomainInput, scaffolderVersion string) error {
	meta, err := s.Load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	existing, exists := meta.Domains[domainName]

	domainMeta := DomainMetadata{
		ScaffolderVersion: scaffolderVersion,
		Input:             input,
	}

	if exists {
		// Preserve original scaffold time, update the updated time
		domainMeta.ScaffoldedAt = existing.ScaffoldedAt
		domainMeta.UpdatedAt = now
	} else {
		domainMeta.ScaffoldedAt = now
	}

	meta.Domains[domainName] = domainMeta
	return s.Save(meta)
}

// SaveWizard saves or updates metadata for a single wizard.
func (s *Store) SaveWizard(wizardName, domain string, input types.ScaffoldWizardInput, scaffolderVersion string) error {
	meta, err := s.Load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	// Use a composite key of domain:wizardName for uniqueness
	key := domain + ":" + wizardName
	existing, exists := meta.Wizards[key]

	wizardMeta := WizardMetadata{
		ScaffolderVersion: scaffolderVersion,
		Domain:            domain,
		Input:             input,
	}

	if exists {
		// Preserve original scaffold time, update the updated time
		wizardMeta.ScaffoldedAt = existing.ScaffoldedAt
		wizardMeta.UpdatedAt = now
	} else {
		wizardMeta.ScaffoldedAt = now
	}

	meta.Wizards[key] = wizardMeta
	return s.Save(meta)
}

// GetDomain retrieves metadata for a specific domain.
func (s *Store) GetDomain(domainName string) (*DomainMetadata, bool, error) {
	meta, err := s.Load()
	if err != nil {
		return nil, false, err
	}

	domain, exists := meta.Domains[domainName]
	if !exists {
		return nil, false, nil
	}

	return &domain, true, nil
}

// ListDomains returns a list of all scaffolded domain names.
func (s *Store) ListDomains() ([]string, error) {
	meta, err := s.Load()
	if err != nil {
		return nil, err
	}

	domains := make([]string, 0, len(meta.Domains))
	for name := range meta.Domains {
		domains = append(domains, name)
	}
	return domains, nil
}

// RemoveDomain removes metadata for a domain.
func (s *Store) RemoveDomain(domainName string) error {
	meta, err := s.Load()
	if err != nil {
		return err
	}

	delete(meta.Domains, domainName)
	return s.Save(meta)
}

// Exists checks if metadata exists for a domain.
func (s *Store) Exists(domainName string) (bool, error) {
	_, exists, err := s.GetDomain(domainName)
	return exists, err
}
