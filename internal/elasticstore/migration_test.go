package elasticstore

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type memoryMigrationState struct {
	finished map[string]bool
	marked   []string
	checkErr error
	markErr  error
}

func (s *memoryMigrationState) Finished(_ context.Context, name string) (bool, error) {
	return s.finished[name], s.checkErr
}

func (s *memoryMigrationState) MarkFinished(_ context.Context, name string) error {
	if s.markErr != nil {
		return s.markErr
	}
	s.finished[name] = true
	s.marked = append(s.marked, name)
	return nil
}

func TestRunMigrationsRunsOnlyUnfinishedInOrder(t *testing.T) {
	state := &memoryMigrationState{finished: map[string]bool{"002": true}}
	var ran []string
	migrations := []Migration{
		{Name: "001", Run: func(context.Context, *opensearchapi.Client) error { ran = append(ran, "001"); return nil }},
		{Name: "002", Run: func(context.Context, *opensearchapi.Client) error { ran = append(ran, "002"); return nil }},
		{Name: "003", Run: func(context.Context, *opensearchapi.Client) error { ran = append(ran, "003"); return nil }},
	}

	err := runMigrations(context.Background(), nil, state, migrations)
	if err != nil {
		t.Fatalf("runMigrations() error = %v", err)
	}
	if want := []string{"001", "003"}; !reflect.DeepEqual(ran, want) {
		t.Fatalf("ran = %v, want %v", ran, want)
	}
	if want := []string{"001", "003"}; !reflect.DeepEqual(state.marked, want) {
		t.Fatalf("marked = %v, want %v", state.marked, want)
	}
}

func TestRunMigrationsStopsOnFailureWithoutMarking(t *testing.T) {
	state := &memoryMigrationState{finished: make(map[string]bool)}
	migrationErr := errors.New("failed change")
	secondRan := false
	migrations := []Migration{
		{Name: "001", Run: func(context.Context, *opensearchapi.Client) error { return migrationErr }},
		{Name: "002", Run: func(context.Context, *opensearchapi.Client) error { secondRan = true; return nil }},
	}

	err := runMigrations(context.Background(), nil, state, migrations)
	if !errors.Is(err, migrationErr) {
		t.Fatalf("runMigrations() error = %v, want %v", err, migrationErr)
	}
	if secondRan {
		t.Fatal("second migration ran after the first failed")
	}
	if len(state.marked) != 0 {
		t.Fatalf("marked = %v, want none", state.marked)
	}
}

func TestRunMigrationsRejectsInvalidListBeforeRunning(t *testing.T) {
	tests := []struct {
		name       string
		migrations []Migration
		wantError  string
	}{
		{name: "empty name", migrations: []Migration{{Run: func(context.Context, *opensearchapi.Client) error { return nil }}}, wantError: "name is empty"},
		{name: "nil function", migrations: []Migration{{Name: "001"}}, wantError: "has no Run function"},
		{name: "duplicate name", migrations: []Migration{
			{Name: "001", Run: func(context.Context, *opensearchapi.Client) error { return nil }},
			{Name: "001", Run: func(context.Context, *opensearchapi.Client) error { return nil }},
		}, wantError: "duplicate migration name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &memoryMigrationState{finished: make(map[string]bool)}
			err := runMigrations(context.Background(), nil, state, tt.migrations)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("runMigrations() error = %v, want error containing %q", err, tt.wantError)
			}
		})
	}
}
