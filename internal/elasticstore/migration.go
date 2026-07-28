package elasticstore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const migrationsIndexName = ".openlibrary-migrations"

// Migration is a named, idempotent change to the OpenSearch store.
//
// A migration is considered finished only after Run returns nil and its name
// has been recorded in OpenSearch. Names must be unique within a migration
// list and must not be changed after a migration has run.
type Migration struct {
	Name string
	Run  func(context.Context, *opensearchapi.Client) error
}

// RunMigrations runs unfinished migrations in the order they are provided.
// It stops at the first error. Migrations that have already finished are
// skipped.
//
// Migration functions should be idempotent: a process can stop after a
// migration succeeds but before its completion is recorded, causing that
// migration to run again on the next attempt.
func RunMigrations(ctx context.Context, client *opensearchapi.Client, migrations []Migration) error {
	if client == nil {
		return errors.New("elasticstore migrations: client is nil")
	}

	return runMigrations(ctx, client, opensearchMigrationState{client: client}, migrations)
}

type migrationState interface {
	Finished(context.Context, string) (bool, error)
	MarkFinished(context.Context, string) error
}

func runMigrations(ctx context.Context, client *opensearchapi.Client, state migrationState, migrations []Migration) error {
	seen := make(map[string]struct{}, len(migrations))
	for _, migration := range migrations {
		if migration.Name == "" {
			return errors.New("elasticstore migrations: migration name is empty")
		}
		if migration.Run == nil {
			return fmt.Errorf("elasticstore migrations: migration %q has no Run function", migration.Name)
		}
		if _, ok := seen[migration.Name]; ok {
			return fmt.Errorf("elasticstore migrations: duplicate migration name %q", migration.Name)
		}
		seen[migration.Name] = struct{}{}
	}

	for _, migration := range migrations {
		finished, err := state.Finished(ctx, migration.Name)
		if err != nil {
			return fmt.Errorf("elasticstore migrations: check %q: %w", migration.Name, err)
		}
		if finished {
			continue
		}

		if err := migration.Run(ctx, client); err != nil {
			return fmt.Errorf("elasticstore migrations: run %q: %w", migration.Name, err)
		}
		if err := state.MarkFinished(ctx, migration.Name); err != nil {
			return fmt.Errorf("elasticstore migrations: mark %q finished: %w", migration.Name, err)
		}
	}

	return nil
}

type opensearchMigrationState struct {
	client *opensearchapi.Client
}

func (s opensearchMigrationState) Finished(ctx context.Context, name string) (bool, error) {
	response, err := s.client.Document.Exists(ctx, opensearchapi.DocumentExistsReq{
		Index:      migrationsIndexName,
		DocumentID: migrationDocumentID(name),
	})
	if err == nil {
		return true, nil
	}
	if response != nil && response.StatusCode == 404 {
		return false, nil
	}

	return false, err
}

func (s opensearchMigrationState) MarkFinished(ctx context.Context, name string) error {
	body, err := json.Marshal(struct {
		Name       string    `json:"name"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
		Name:       name,
		FinishedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Index(ctx, opensearchapi.IndexReq{
		Index:      migrationsIndexName,
		DocumentID: migrationDocumentID(name),
		Body:       strings.NewReader(string(body)),
	})
	return err
}

func migrationDocumentID(name string) string {
	sum := sha256.Sum256([]byte(name))
	return hex.EncodeToString(sum[:])
}
